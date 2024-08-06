package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"goboot/internal/config"
	"goboot/internal/model"
	"goboot/internal/repository"
	"goboot/internal/vo"
	"goboot/internal/vo/request"
	"goboot/internal/vo/response"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-contrib/i18n"
	"github.com/gophero/goal/assert"
	"github.com/gophero/goal/aws/s3"
	"github.com/gophero/goal/date"
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/imagex"
	"github.com/gophero/goal/mailx"
	"github.com/gophero/goal/random"
	"github.com/gophero/goal/redisx"
	"golang.org/x/image/bmp"
	"gorm.io/gorm"
)

type UserService struct {
	*Service

	urpo *repository.UserRepository
	s3c  *s3.Client
	conf config.Conf
	rc   redisx.Client

	sssrv  *SysSettingService
	sdsrv  *SysDeviceService
	imgsrv *ImageService
}

func NewUserService(service *Service, userRepo *repository.UserRepository, conf config.Conf, s3c *s3.Client, rc redisx.Client,
	sdsrv *SysDeviceService, sssrv *SysSettingService, imgsrv *ImageService,
) *UserService {
	return &UserService{
		Service: service,
		s3c:     s3c,
		conf:    conf,
		rc:      rc,
		urpo:    userRepo,
		sssrv:   sssrv,
		sdsrv:   sdsrv,
		imgsrv:  imgsrv,
	}
}

const (
	userAvatarDir          = "avatars"
	userBannerDir          = "banners"
	userMaxAvatarWidth     = 400
	userMaxBannerWidth     = 800
	UserCompressTypeAvatar = 1
	UserCompressTypeBanner = 2
	UniqueNumPrefix        = "E"

	minNickNameLen = 2  // 昵称最小字符长度, 一个中文3个字符
	maxNickNameLen = 45 // 昵称最大字符长度, 一个中文3个字符
	minNameLen     = 2  // 姓名最小字符长度, 一个中文3个字符
	maxNameLen     = 45 // 姓名最大字符长度, 一个中文3个字符
	maxIntroLen    = 90 // 简介最大字符长度，一个中文3个字符

	UserStatusForbidden  = -2 // 禁用
	UserStatusLocked     = -1 // 锁定
	UserStatusRegistered = 1  // 已注册
	UserStatusVerified   = 2  // 已验证
)

const lockVerifyCodeRedisKey = "lockVerifyCode:%s"

func (srv *UserService) Add(user *model.User) error {
	return srv.urpo.W.Save(&user).Error
}

// 再次确保邀请码未被使用，防止并发问题
func (srv *UserService) makesureDailyInviteCodeNotUsed(tx *gorm.DB, inviteCode string) error {
	if inviteCode != "" {
		useDaily := srv.sssrv.GetSystemSettingByKey("use_daily_invite_code").SysValue == "1"
		if useDaily { // 使用每日邀请码，此时一个邀请码只能用一次
			var mcode model.UserInviteCode
			r := tx.Model(mcode).Where("code = ?", inviteCode).Scan(&mcode)
			if r.Error != nil {
				return r.Error
			}
			if mcode.ID > 0 {
				if mcode.HasUsed {
					return errorx.NewPreferredErrf("invite code has been used")
				}
			} else {
				return errorx.NewPreferredErrf("invite code not found")
			}
		}
	}
	return nil
}

func (srv *UserService) Register(u *model.User, inviteCode string) bool {
	u.Name = u.Email[0:strings.Index(u.Email, "@")] // 默认name为邮箱名
	u.UniqueNumber = srv.GenerateUniqueNumber()     // number
	u.Username = u.UniqueNumber                     // username
	u.Status = UserStatusVerified                   // 直接改为验证通过状态
	u.NickName = u.UniqueNumber
	err := srv.urpo.W.Transaction(func(tx *gorm.DB) error {
		errorx.Throw(srv.makesureDailyInviteCodeNotUsed(tx, inviteCode))
		errorx.Throw(tx.Create(&u).Error)
		// act := &model.Account{UserId: u.ID}
		// errorx.Throw(tx.Create(act).Error)
		// 设置父级数据
		if u.ParentId > 0 {
			tx.Model(&model.User{}).Where("id", u.ParentId).Update("invite_num", gorm.Expr("invite_num+1"))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return true
}

func (srv *UserService) SaveUser(u *model.User) bool {
	if u.Name == "" {
		u.Name = u.Email[0:strings.Index(u.Email, "@")] // 默认name为邮箱名
	}
	u.UniqueNumber = srv.GenerateUniqueNumber() // number
	u.Status = UserStatusVerified               // 直接改为验证通过状态
	if u.Username == "" {
		u.Username = u.UniqueNumber
	}
	if u.NickName == "" {
		u.NickName = u.Name
	}
	err := srv.urpo.W.Transaction(func(tx *gorm.DB) error {
		errorx.Throw(tx.Create(&u).Error)
		// act := &model.Account{UserId: u.ID}
		// errorx.Throw(tx.Create(act).Error)
		return nil
	})
	if err != nil {
		srv.L.Errorf("save user error: %v", err)
		return false
	}
	return true
}

func (srv *UserService) GenerateUniqueNumber() (num string) {
	start := 100000
	end := 9999999
	for i := 0; true; i++ {
		num = strconv.Itoa(random.Between(start, end))
		var dbUser model.User
		err := srv.urpo.R.Where("unique_number = ?", num).Find(&dbUser).Error
		errorx.Throw(err)
		if dbUser.ID == 0 {
			return UniqueNumPrefix + num
		}
		if i != 0 && i%3 == 0 {
			start = start * 10
			end = end*10 + 9
		}
	}
	return
}

func (srv *UserService) Del(userId uint) {
	err := srv.urpo.W.Where("id", userId).Delete(&model.User{}).Error
	errorx.Throw(err)
}

func (srv *UserService) GetById(userId uint) *model.User {
	u := &model.User{}
	err := srv.urpo.R.Where("id", userId).First(&u).Error
	errorx.Throw(err)
	return u
}

func (srv *UserService) GetByIdTx(tx *gorm.DB, userId uint) *model.User {
	u := &model.User{}
	err := tx.Where("id", userId).First(&u).Error
	errorx.Throw(err)
	return u
}

func (srv *UserService) GetByUniqueNum(un string) *model.User {
	u := &model.User{}
	err := srv.urpo.R.Where("unique_number = ?", un).Find(&u).Error
	errorx.Throw(err)
	return u
}

func (srv *UserService) GetByEmail(email string) *model.User {
	assert.NotBlank(email)
	u := &model.User{}
	r := srv.urpo.R.Where("email = ?", email).Find(u)
	if r.Error != nil || r.RowsAffected < 1 {
		return nil
	}
	return u
}

func (srv *UserService) ResetPwd(email string, password string) error {
	assert.NotBlank(email)
	assert.NotBlank(password)
	usr := srv.GetByEmail(email)
	if usr == nil {
		return errors.New(i18n.MustGetMessage("account_not_exist"))
	}
	usr.Password = password
	_, err := srv.urpo.Update(usr, srv.urpo.W)
	return err
}

func (srv *UserService) UpdPwd(userId uint, password string, token string) error {
	err := srv.urpo.W.Model(&model.User{}).Where("id=?", userId).Update("password", password).Error
	if err == nil {
		_, err = srv.rc.Del(context.Background(), token).Result()
		if err != nil {
			srv.L.Warnf("delete token from redis failed: %v", err)
		}
		return nil
	} else {
		return err
	}
}

func (srv *UserService) SignOut(token string) error {
	_, err := srv.rc.Del(context.Background(), token).Result()
	return err
}

func (srv *UserService) UpdateUser(request request.UpdateUserReq) {
	dbuser := srv.GetById(request.LoginUserId)
	if request.NickName != "" && dbuser.NickName != request.NickName {
		assert.Require(len(request.NickName) >= minNickNameLen && len(request.NickName) <= maxNickNameLen, i18n.MustGetMessage("nick_name_over_limit"))
		dbuser.NickName = request.NickName
	}
	if request.Gender != 0 {
		dbuser.Gender = request.Gender
	}
	if request.Intro != "" {
		assert.Require(len(request.Intro) <= maxIntroLen, i18n.MustGetMessage("intro_over_limit"))
		dbuser.Intro = request.Intro
	}
	if request.BirthDay != "" {
		t := date.ParseDate(request.BirthDay)
		dbuser.BirthDay = &t
	}
	if request.Address != "" {
		dbuser.Address = request.Address
	}
	err := srv.urpo.W.Updates(&dbuser).Error
	errorx.Throw(err)
}

func (srv *UserService) Heartbeat(userId uint) {
	usr := srv.GetById(userId)
	// 一个小时以内不更新
	if time.Now().Sub(usr.ActiveTime).Hours() > 1 {
		srv.UpdateActiveTime(userId)
	}
}

func (srv *UserService) UpdateActiveTime(userId uint) {
	err := srv.urpo.W.Model(&model.User{}).Where("id=?", userId).Update("active_time", time.Now()).Error
	errorx.Throw(err)
}

// UploadAvatar 上传头像，返回头像 url、是否违规和 error
func (srv *UserService) UploadAvatar(userId uint, file io.Reader, fileName string, ip string) (string, string, bool, error) {
	// 压缩图片
	r := srv.CompressImage(file, fileName, UserCompressTypeAvatar)
	var relativePath string
	url, illegal, err := srv.imgsrv.UploadImage(userId, r, fileName, ip, userAvatarDir, func(url string) error {
		relativePath = url
		dbUser := &model.User{}
		err := srv.urpo.R.Where("id = ?", userId).First(dbUser).Error
		errorx.Throw(err)
		err = srv.urpo.W.Model(&model.User{}).Where("id = ?", userId).Update("avatar", url).Error
		errorx.Throw(err)
		// 删除旧图片
		if strings.Contains(dbUser.Avatar, userAvatarDir) {
			if err := srv.s3c.DeleteFile(srv.conf.AwsS3().Bucket, dbUser.Avatar); err != nil {
				srv.L.Errorf("delete illegal avatar failed, objectKey: %srv, error: %srv", dbUser.Avatar, err.Error())
			}
		}
		return nil
	}, func() {
		if err := srv.s3c.DeleteFile(srv.conf.AwsS3().Bucket, userAvatarDir+"/"+fileName); err != nil {
			srv.L.Errorf("delete illegal avatar failed, objectKey: %srv, error: %srv", fileName, err.Error())
		}
	})
	return url, relativePath, illegal, err
}

func (srv *UserService) CompressImage(file io.Reader, filename string, typ int) (r io.Reader) {
	bs, err := io.ReadAll(file)
	fr := bytes.NewReader(bs)
	srcr := bytes.NewReader(bs)
	errorx.Throw(err)

	var w uint
	switch typ {
	case UserCompressTypeAvatar:
		w = userMaxAvatarWidth
	case UserCompressTypeBanner:
		w = userMaxBannerWidth
	}

	// var dstr io.Reader
	defer func() {
		if err := recover(); err != nil {
			srv.L.Errorf("compress avatar failed, filename: %s, use original image instead", filename)
			r = srcr // 出错了返回原始 reader，不能直接给 dstr 赋值，不能执行到 return
		}
	}()
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".png":
		img, err := png.Decode(fr)
		errorx.Throw(err)
		img = imagex.Thumbnail(w, w, img)
		var buf bytes.Buffer
		errorx.Throw(png.Encode(&buf, img))
		r = bytes.NewReader(buf.Bytes())
	case ".gif":
		img, err := gif.Decode(fr)
		errorx.Throw(err)
		img = imagex.Thumbnail(w, w, img)
		var buf bytes.Buffer
		errorx.Throw(gif.Encode(&buf, img, nil))
		r = bytes.NewReader(buf.Bytes())
	case ".bmp":
		img, err := bmp.Decode(fr)
		errorx.Throw(err)
		img = imagex.Thumbnail(w, w, img)
		var buf bytes.Buffer
		errorx.Throw(bmp.Encode(&buf, img))
		r = bytes.NewReader(buf.Bytes())
	case ".jpg", ".jpeg": // jpg
		img, err := jpeg.Decode(fr)
		errorx.Throw(err)
		img = imagex.Thumbnail(w, w, img)
		var buf bytes.Buffer
		errorx.Throw(jpeg.Encode(&buf, img, nil))
		r = bytes.NewReader(buf.Bytes())
	default:
		r = srcr // 不支持压缩
	}
	// return dstr
	return r
}

// UploadBanner 上传个人背景图，返回 url、是否违规和 error
func (srv *UserService) UploadBanner(userId uint, file multipart.File, fileName string, ip string) (string, bool, error) {
	// 压缩图片
	r := srv.CompressImage(file, fileName, UserCompressTypeBanner)
	return srv.imgsrv.UploadImage(userId, r, fileName, ip, userBannerDir, func(url string) error {
		return srv.urpo.W.Model(&model.User{}).Where("id = ?", userId).Update("banner", url).Error
	}, func() {
		if err := srv.s3c.DeleteFile(srv.conf.AwsS3().Bucket, fileName); err != nil {
			srv.L.Errorf("delete illegal banner failed, objectKey: %srv, error: %srv", fileName, err.Error())
		}
	})
}

func (srv *UserService) BindDevice(userId uint, vo *request.SysDeviceVo) error {
	srcDeviceInfo := srv.sdsrv.GetDeviceInfoByUser(userId, vo)
	var err error
	if srcDeviceInfo == nil {
		vo.UserId = userId
		_, err = srv.urpo.CreateOne(newDeviceInfo(vo))
	} else {
		err = srv.urpo.UpdateFields(copyFields(srcDeviceInfo, vo), nil)
	}
	return err
}

func newDeviceInfo(vo *request.SysDeviceVo) *model.SysDevice {
	return &model.SysDevice{
		DeviceId:          vo.DeviceId,
		Channel:           vo.Channel,
		Platform:          vo.Platform,
		VersionInfo:       vo.VersionInfo,
		UserId:            vo.UserId,
		DeviceModel:       vo.DeviceModel,
		DeviceBrand:       vo.DeviceBrand,
		DeviceType:        vo.DeviceType,
		DeviceOrientation: vo.DeviceOrientation,
		DevicePixelRatio:  vo.DevicePixelRatio,
		System:            vo.System,
		SysLang:           vo.SysLang,
	}
}

func copyFields(m *model.SysDevice, vo *request.SysDeviceVo) *model.SysDevice {
	m.DeviceId = vo.DeviceId
	m.Channel = vo.Channel
	m.Platform = vo.Platform
	m.VersionInfo = vo.VersionInfo
	m.UserId = vo.UserId
	m.DeviceModel = vo.DeviceModel
	m.DeviceBrand = vo.DeviceBrand
	m.DeviceType = vo.DeviceType
	m.DeviceOrientation = vo.DeviceOrientation
	m.DevicePixelRatio = vo.DevicePixelRatio
	m.System = vo.System
	m.SysLang = vo.SysLang
	return m
}

func (srv *UserService) CheckValidEmail(email string, checkDomain bool) bool {
	if !mailx.ValidMailAddress(email) {
		return false
	}
	if !checkDomain {
		return true
	}
	mc := srv.conf.Mail()
	if mc.SupportDomain == "" { // 不限制
		return true
	}
	ds := strings.Split(mc.SupportDomain, ",")
	domain := strings.Split(email, "@")[1]
	for _, v := range ds {
		if v == domain {
			return true
		}
	}
	return false
}

// CheckVerifyCode 校验验证码
func (srv *UserService) CheckVerifyCode(email string, verifyCode string) bool {
	assert.NotBlank(verifyCode)
	// 加入redis锁,防止重复提交
	key := fmt.Sprintf(lockVerifyCodeRedisKey, email)
	cmd := srv.rc.SetNX(context.Background(), key, verifyCode, time.Second*30)
	if cmd.Err() != nil {
		srv.L.Errorf("Lock VerifyCode error email:%s  error:%v", email, cmd.Err())
		return false
	}
	if cmd.Val() {
		defer func() { // 解锁
			srv.rc.Del(context.Background(), key)
		}()
	} else {
		srv.L.Errorf("Lock VerifyCode false email:%s", email)
		return false
	}
	// 验证验证码
	v, err := srv.rc.Get(context.Background(), email).Result()
	if err != nil {
		return false
	}
	bl := v == verifyCode
	if bl {
		//_, err = srv.rc.Del(context.Background(), email).Result()
		if err != nil {
			srv.urpo.L.Errorf("CheckVerifyCode redis del error:%v", err)
		}
	}
	return bl
}

// FindByUserIds 按用户id集合查询获取关注的用户列表
func (srv *UserService) FindByUserIds(destUserIds []uint) (resp map[uint]model.User) {
	where := map[string]interface{}{}
	where["id"] = destUserIds
	var list []model.User
	err := srv.urpo.R.Where(where).Find(&list).Error
	errorx.Throw(err)
	mp := map[uint]model.User{}
	for _, r := range list {
		mp[r.ID] = r
	}
	return mp
}

func (srv *UserService) FindInternalUsers() []*model.User {
	var us []*model.User
	errorx.Throw(srv.urpo.R.Model(model.User{}).Where("internal = ? and status = ?", true, model.UserStatusVerified).Find(&us).Error)
	return us
}

func (srv *UserService) FindInternalUserCount() int64 {
	var cnt int64
	srv.urpo.R.Raw("select count(*) from users where internal = ? and status = ?", true, model.UserStatusVerified).Find(&cnt)
	return cnt
}

func (srv *UserService) FindVipUsers() []*model.User {
	var us []*model.User
	errorx.Throw(srv.urpo.R.Model(model.User{}).Where("vip = ? and status = ?", true, model.UserStatusVerified).Find(&us).Error)
	return us
}

func (srv *UserService) GetUserParentInfoById(id uint) *response.UserParent {
	sql := `select p.id, p.email, p.name, p.avatar, w.address
from (
     select u.id, u.email, u.name, u.avatar from users u
     where u.id = (select parent_id from users where id = ?)
 ) p left join user_wallets w on w.user_id = p.id and w.state = 0`
	var user response.UserParent
	r := srv.urpo.R.Raw(sql, id).Scan(&user)
	if r.Error != nil {
		panic(r.Error)
	}
	return &user
}

func (srv *UserService) GetCountByParentId(userId uint) (count int64) {
	err := srv.urpo.R.Model(&model.User{}).Where("parent_id", userId).Count(&count).Error
	errorx.Throw(err)
	return
}

func (srv *UserService) GetUserByParentId(userId uint, page vo.Page) []model.User {
	list := []model.User{}
	offset := page.PageSize * (page.PageNum - 1)
	err := srv.urpo.R.Where("parent_id = ? and status = ?", userId, model.UserStatusVerified).Offset(offset).Limit(page.PageSize).Order("vip_level desc,active_time desc").Find(&list).Error
	errorx.Throw(err)
	return list
}

type UserAvatar struct {
	Id     uint
	VipLvl int
	Avatar string
}

func (srv *UserService) getUserAvatars(userIds ...uint) ([]*UserAvatar, error) {
	var rs []*UserAvatar
	r := srv.urpo.R.Model(model.User{}).Select("id", "avatar", "vip_level").Where("id in ?", userIds).Find(&rs)
	if r.Error != nil {
		srv.L.Errorf("query user avatars error: %v", r.Error)
		return []*UserAvatar{}, r.Error
	}
	return rs, nil
}

// SumFriendReward 统计好友贡献奖励
func (srv *UserService) SumFriendReward(parentId uint) (amount float64) {
	err := srv.R.R.Raw("SELECT IFNULL(SUM(com_amt),0) FROM users WHERE parent_id=?", parentId).Find(&amount).Error
	errorx.Throw(err)
	return
}

func (srv *UserService) QueryNeedRefreshVipId(cnt uint) ([]uint, error) {
	assert.True(cnt > 0)
	var (
		ids []uint
		err error
		h   = errorx.NewHandler()
	)
	defer h.Done(func(err error) {
		if err != nil {
			srv.L.Errorf("error: %v", h.Err())
		}
	})
	h.Do(func() error {
		/* 		ids, err = srv.urpo.QueryNeedRefreshVipId(srv.urpo.R, int64(cnt)) */
		return err
	})
	if h.HasErr() {
		return []uint{}, h.PreferredOr(errorx.ServerBusy)
	}
	return ids, nil
}

func (srv *UserService) GetByUsername(username string) *model.User {
	assert.NotBlank(username)
	u := &model.User{}
	r := srv.R.W.Where("username = ?", username).Find(u)
	if r.Error != nil || r.RowsAffected < 1 {
		if r.Error != nil {
			srv.L.Errorf("get user by username error:%v", r.Error)
		}
		return nil
	}
	return u
}

func (srv *UserService) GetByIds(userIds []uint) []*model.User {
	var u []*model.User
	err := srv.urpo.R.Where("id in ?", userIds).Find(&u).Error
	errorx.Throw(err)
	return u
}

func (srv *UserService) GetEmailMapByIds(userIds []uint) map[uint]string {
	var users []model.User
	err := srv.urpo.R.Where("id in ?", userIds).Select("id", "email").Find(&users).Error
	errorx.Throw(err)
	resp := map[uint]string{}
	for _, user := range users {
		resp[user.ID] = user.Email
	}
	return resp
}

func (srv *UserService) GetByEmails(emails []string) []*model.User {
	assert.True(len(emails) > 0)
	var us []*model.User
	r := srv.urpo.R.Where("email in ?", emails).Find(&us)
	if r.Error != nil || r.RowsAffected < 1 {
		return nil
	}
	return us
}
