package handler

import (
	"context"
	"fmt"
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/mailx"
	"github.com/gophero/goal/redisx"
	"github.com/gophero/goal/stringx"
	"goboot/internal/consts"
	"goboot/internal/model"
	"goboot/internal/service"
	"goboot/internal/vo/request"
	"goboot/internal/vo/response"
	"goboot/pkg/util"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/i18n"
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/assert"
	"github.com/gophero/goal/conv"
)

const (
	verifyCodeLogin = 1

	loginErrorDuration  = 5 * time.Minute
	loginLockedDuration = 10 * time.Minute
	loginErrorKey       = "login_error_%s"
	loginLockedKey      = "login_lock_%v"
)

type LoginHandler struct {
	*Handler
	rc   redisx.Client
	msdr *mailx.Sender

	usrv *service.UserService
	ssrv *service.SysSettingService
	vsrv *service.AppVersionService
	dsrv *service.SysDeviceService
}

func NewLoginHandler(handler *Handler, rc redisx.Client, msdr *mailx.Sender,
	usrv *service.UserService, ssrv *service.SysSettingService, vsrv *service.AppVersionService,
	dsrv *service.SysDeviceService) *LoginHandler {
	return &LoginHandler{Handler: handler, rc: rc, msdr: msdr, usrv: usrv, ssrv: ssrv, vsrv: vsrv, dsrv: dsrv}
}

type InviteCodeReq struct {
	InviteCode string `json:"inviteCode" form:"inviteCode"`
}

type RegisterReq struct {
	Email      string `json:"email" form:"email"`
	VerifyCode string `json:"verifyCode" form:"verifyCode"`
	Password   string `json:"password" form:"password"`
	InviteCodeReq
}

func (api *LoginHandler) RegisterRouter(r *gin.Engine) {
	g := api.Group("/login", r)
	{
		g.POST("/sign_in", api.Login)              // 登录
		g.POST("/sign_up", api.Register)           // 注册
		g.GET("/version_info", api.VersionInfo)    // 版本信息
		g.GET("/email_check", api.CheckMail)       // 邮箱合法性检查
		g.POST("/verify_code", api.SendVerifyCode) // 发送邮箱验证码
		g.POST("/reset_pwd", api.ResetPassword)    //重置密码
	}
}

func (api *LoginHandler) CheckMail(c *gin.Context) {
	var loginReq request.SendVerifyRequest
	request.Bind(c, &loginReq)
	assert.Require(loginReq.Email != "", "missing email")
	if !api.usrv.CheckValidEmail(loginReq.Email, true) {
		response.Ctx(c).BadReqI18n("invalid_email")
		return
	}
	if api.usrv.GetByEmail(loginReq.Email) != nil {
		response.Ctx(c).BadReqI18n("email_registered")
		return
	}
	response.Ctx(c).Ok()
}

func (api *LoginHandler) Register(c *gin.Context) {
	req := &RegisterReq{}
	loginReq := request.Bind(c, request.LoginReq{})
	errorx.Throw(c.Bind(req))
	assert.NotBlank(req.Email)
	assert.NotBlank(req.VerifyCode)
	// assert.NotBlank(req.InviteCode)
	assert.NotBlank(req.Password)

	if !api.usrv.CheckValidEmail(req.Email, true) {
		response.Ctx(c).ServerI18n("invalid_email")
		return
	}

	usr := api.usrv.GetByEmail(req.Email)
	if usr != nil {
		response.Ctx(c).ServerI18n("email_registered")
		return
	}
	//校验验证码
	if !api.usrv.CheckVerifyCode(req.Email, req.VerifyCode) {
		response.Ctx(c).ServerI18n("invalid_verify_code")
		return
	}

	usr = &model.User{Email: req.Email, Password: req.Password}
	if api.usrv.Register(usr, req.InviteCode) {
		api.rc.Del(context.Background(), req.Email)
		token := util.ProduceToken(api.rc, c, usr.ID)
		// 绑定设备，出错不影响登录
		appChan := c.GetHeader(consts.AppChannelKey)
		var cn string
		if appChan != "" {
			cn = consts.GetChanName(conv.StrToInt(appChan))
			if cn == "" {
				cn = loginReq.Platform
			}
		} else {
			cn = loginReq.Platform
		}
		loginReq.Channel = cn
		loginReq.SysLang = c.GetHeader(consts.SysLang)
		api.bindDevice(usr.ID, loginReq)
		response.Ctx(c).OkJson(map[string]any{"token": token})
	} else {
		response.Ctx(c).ServerErrf("register failed")
	}
}

func (api *LoginHandler) bindDevice(userId uint, req request.LoginReq) {
	defer func() {
		if err := recover(); err != nil {
			api.L.Error("bind device error: ", err)
		}
	}()
	// if req.DeviceId != "" {
	api.usrv.BindDevice(userId, &request.SysDeviceVo{
		DeviceId:          req.DeviceId, // 设备ID
		Channel:           req.Channel,
		Platform:          req.Platform,    // 平台，mac 或 android
		VersionInfo:       req.VersionInfo, // 版本信息
		UserId:            userId,
		DeviceModel:       req.DeviceModel, // 设备型号
		DeviceBrand:       req.DeviceBrand, // 设备品牌
		DeviceType:        req.DeviceType,  // 设备类型
		DeviceOrientation: req.DeviceOrientation,
		DevicePixelRatio:  req.DevicePixelRatio,
		System:            req.System, // 操作系统及版本
		SysLang:           req.SysLang,
	})
	// } else {
	// 	config.Logger.Warning("No deviceId found, user id: %d, email: %s, device brand: %s", userId, req.Email, req.DeviceBrand)
	// }
}

func (api *LoginHandler) SendVerifyCode(c *gin.Context) {
	// token := c.PostForm("cap_token")
	// key := g.Conf.System.CaptchaKey
	// reCAPTCHA, _ := recaptcha.NewReCAPTCHA(key, recaptcha.V2, 10*time.Second) // for v2 API get your secret from https://www.google.com/recaptcha/admin
	// err := reCAPTCHA.Verify(token)
	// if err != nil {
	//	response.Ctx(c).BadReq()
	//	return
	// }
	var loginReq request.SendVerifyRequest
	request.Bind(c, &loginReq)
	assert.NotBlank(loginReq.Email)
	if !api.usrv.CheckValidEmail(loginReq.Email, false) {
		response.Ctx(c).ServerI18n("invalid_email")
		return
	}
	// 判断是否校验邀请码
	//if loginReq.Type == "register" {
	//	bl, _ := api.usrv.VerifyInviteCode(loginReq.InviteCode)
	//	if !bl {
	//		response.Ctx(c).BadReqI18n("invalid_invite_code")
	//		return
	//	}
	//}
	api.msdr.SendVerifyCode(api.rc, loginReq.Email, "")
	response.Ctx(c).Ok()
}

func (api *LoginHandler) Login(c *gin.Context) {
	req := request.Bind(c, request.LoginReq{})
	if c.Request.UserAgent() == "insomnia/6.4.1" {
		c.String(http.StatusBadRequest, "lower version")
		return
	}
	// em := c.PostForm("email")
	em := req.Email
	assert.Require(strings.TrimSpace(em) != "", "missing email address")
	usr := api.usrv.GetByEmail(em)
	if usr == nil {
		response.Ctx(c).BadReqf(i18n.MustGetMessage("password_error"))
		return
	}
	if usr.IsLocked() || isLoginLock(api.rc, c, em) {
		response.Ctx(c).BadReqf(i18n.MustGetMessage("password_error"))
		return
	}
	// 校验密码或验证码是否正确
	if req.Type != verifyCodeLogin {
		password := req.Password
		assert.Require(strings.TrimSpace(password) != "", "missing password")
		if usr.Password != "" && usr.Password != password {
			count := loginError(api.rc, c, em)
			if count > 10 {
				loginLocked(api.rc, c, em)
				response.Ctx(c).BadReqf(i18n.MustGetMessage("password_error"))
				return
			}
			response.Ctx(c).BadReqf(i18n.MustGetMessage("password_error"))
			return
		}
	} else {
		assert.Require(strings.TrimSpace(req.VerifyCode) != "", "missing verification code")
		if !api.usrv.CheckVerifyCode(req.Email, req.VerifyCode) {
			response.Ctx(c).ServerI18n("invalid_verify_code")
			return
		}
	}

	if usr.IsForbidden() {
		response.Ctx(c).BadReqf("cheat account, suspend service!")
		return
	}
	token := util.ProduceToken(api.rc, c, usr.ID)
	// 绑定设备，出错不影响登录
	appChan := c.GetHeader(consts.AppChannelKey)
	var cn string
	if appChan != "" {
		cn = consts.GetChanName(conv.StrToInt(appChan))
		if cn == "" {
			cn = req.Platform
		}
	} else {
		cn = req.Platform
	}
	req.Channel = cn
	req.SysLang = c.GetHeader(consts.SysLang)
	api.bindDevice(usr.ID, req)
	if req.Type == 1 {
		api.rc.Del(context.Background(), req.Email)
	}

	response.Ctx(c).OkJson(map[string]any{"token": token})
}

func (api *LoginHandler) VersionInfo(c *gin.Context) {
	versionReq := c.GetHeader("Version")
	assert.Require(versionReq != "", i18n.MustGetMessage("param_error"))
	appChan := conv.StrToInt(c.GetHeader(consts.AppChannelKey))
	version := api.vsrv.GetLatestVersion(consts.AppName, consts.GetChanName(appChan))
	//api.L.Infof("appChan %d", appChan)
	//api.L.Infof("version %s", versionReq)
	pass := 1
	// 获取系统参数
	keys := []string{consts.HeartbeatSwitch, consts.HeartbeatInterval, consts.IOSPassKey}
	settingMap := api.ssrv.GetSystemSettingsValByKeys(keys...)
	heartbeat := gin.H{"switch": settingMap[consts.HeartbeatSwitch], "interval": settingMap[consts.HeartbeatInterval]}

	if versionReq > version.VersionInfo {
		pass = conv.StrToInt(settingMap[consts.IOSPassKey])
	}

	if versionReq >= version.VersionInfo {
		response.Ctx(c).OkJson(map[string]any{
			"pass":      pass,
			"version":   nil,
			"heartbeat": heartbeat,
		})
		return
	}
	// 封装数据
	rep := map[string]any{
		"pass": pass,
		"version": &response.VersionResp{
			Version: version.VersionInfo,
			Url:     version.Url,
			Content: version.Content,
			Force:   version.Force,
		},
		"heartbeat": heartbeat,
	}
	response.Ctx(c).OkJson(rep)
}

func (api *LoginHandler) ResetPassword(c *gin.Context) {
	req := request.Bind(c, request.LoginReq{})
	//校验验证码
	if !api.usrv.CheckVerifyCode(req.Email, req.VerifyCode) {
		response.Ctx(c).ServerI18n("invalid_verify_code")
		return
	}
	err := api.usrv.ResetPwd(req.Email, req.Password)
	if err == nil {
		api.rc.Del(context.Background(), req.Email)
		response.Ctx(c).Ok()
		return
	}
	response.Ctx(c).BadReqf(err.Error())
}

// 判断ios是否过审
func (api *LoginHandler) isIosAudit(appChannel int, val string) bool {
	if val == "0" {
		if appChannel == consts.AppChannelIosApp {
			return true
		}
	}
	return false
}

func loginError(client redisx.Client, c context.Context, email string) int {
	key := fmt.Sprintf(loginErrorKey, email)
	countStr := client.Get(c, key)
	if stringx.IsEmpty(countStr.Val()) {
		client.Set(c, key, 1, loginErrorDuration)
		return 1
	}
	count := errorx.Throwv(countStr.Int())
	count += 1
	client.Set(c, key, count, loginErrorDuration)
	return count
}

func loginLocked(client redisx.Client, c context.Context, email string) {
	key := fmt.Sprintf(loginLockedKey, email)
	client.Set(c, key, 1, loginLockedDuration)
}

func isLoginLock(client redisx.Client, c context.Context, email string) bool {
	key := fmt.Sprintf(loginLockedKey, email)
	countStr := client.Get(c, key)
	if stringx.IsEmpty(countStr.Val()) {
		return false
	}
	return true
}
