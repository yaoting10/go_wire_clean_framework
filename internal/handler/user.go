package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/uuid"
	"goboot/internal/config"
	"goboot/internal/service"
	"goboot/internal/vo/request"
	"goboot/internal/vo/response"
	"goboot/pkg/util"
	"path/filepath"
)

type UserHandler struct {
	*Handler
	conf   config.Conf
	usrv   *service.UserService
	versrv *service.AppVersionService
}

func NewUserHandler(handler *Handler, conf config.Conf, usrv *service.UserService, versrv *service.AppVersionService) *UserHandler {
	return &UserHandler{
		Handler: handler,
		conf:    conf,
		usrv:    usrv,
		versrv:  versrv,
	}
}

func (api *UserHandler) RegisterRouter(r *gin.Engine) {
	var g = api.Group("/user", r)
	{
		// user
		g.GET("", api.findByUserId)
		g.PUT("", api.updUser)
		g.POST("avatar", api.uploadAvatar)
		g.POST("banner", api.uploadBanner)
		g.POST("logout", api.logout)
		g.POST("upd_pwd", api.updPwd)
		g.DELETE("", api.Del) // 注销账户
		g.POST("report", api.report)
		g.GET("/language", api.getLanguage)
	}

	// SLB健康检查
	hg := api.Group("/health", r)
	{
		hg.GET("", api.healthCheck)
	}
}

// uploadAvatar 上传用户头像
func (api *UserHandler) uploadAvatar(c *gin.Context) {
	req := request.BindBase(c)
	file, err := c.FormFile("file")
	f, err := file.Open()
	if err != nil {
		response.Ctx(c).BadReqf(err.Error())
		return
	}
	file.Filename = uuid.UUID32() + filepath.Ext(file.Filename)
	url, _, illegal, err := api.usrv.UploadAvatar(req.LoginUserId, f, file.Filename, c.RemoteIP())
	if err != nil {
		response.Ctx(c).ServerErrf(err.Error())
		return
	}
	if illegal {
		response.Ctx(c).BadReqI18n("illegal_image")
	} else {
		response.Ctx(c).OkJson(gin.H{"url": url})
	}
}

// uploadBanner 上传用户个人背景图
func (api *UserHandler) uploadBanner(c *gin.Context) {
	req := request.BindBase(c)
	file, err := c.FormFile("file")
	f, err := file.Open()
	if err != nil {
		response.Ctx(c).BadReqf(err.Error())
		return
	}
	file.Filename = uuid.UUID32() + filepath.Ext(file.Filename)
	url, illegal, err := api.usrv.UploadBanner(req.LoginUserId, f, file.Filename, c.RemoteIP())
	if err != nil {
		response.Ctx(c).ServerErrf(err.Error())
		return
	}
	if illegal {
		response.Ctx(c).BadReqI18n("illegal_image")
	} else {
		response.Ctx(c).OkJson(gin.H{"url": url})
	}
}

// findByUserId 查询用户信息
func (api *UserHandler) findByUserId(c *gin.Context) {
	var req request.InfoReq
	request.Bind(c, &req)
	if req.UserId == 0 {
		req.UserId = req.LoginUserId
	}
	usr := api.usrv.GetById(req.UserId)
	// 判断是否查询自己的信息
	email, telegram, discord := "", "", ""
	if req.LoginUserId == usr.ID {
		email = usr.Email
		//telegram = usr.TelegramString()
		//discord = usr.DiscordString()
	}
	response.Ctx(c).OkJson(gin.H{
		"avatar":   api.conf.AwsS3().FormatUrl(usr.Avatar),
		"gender":   usr.Gender,
		"number":   usr.UniqueNumber,
		"email":    email,
		"telegram": telegram,
		"discord":  discord,
		"intro":    usr.Intro,
		"nickName": usr.NickName,
	})
}

// updUser 更新用户信息
func (api *UserHandler) updUser(c *gin.Context) {
	var req request.UpdateUserReq
	request.Bind(c, &req)
	api.usrv.UpdateUser(req)
	response.Ctx(c).Ok()
}

// updPwd 修改密码
func (api *UserHandler) updPwd(c *gin.Context) {
	req := request.Bind(c, &request.UpdPwdReq{})
	dbUser := api.usrv.GetById(req.LoginUserId)
	if dbUser.Password != req.OldPassword {
		response.Ctx(c).BadReqI18n("password_error")
		return
	}
	err := api.usrv.UpdPwd(dbUser.ID, req.NewPassword, util.GetToken(c))
	if err == nil {
		response.Ctx(c).Ok()
	} else {
		response.Ctx(c).BadReqf(err.Error())
	}
}

// logout 退出登录
func (api *UserHandler) logout(c *gin.Context) {
	errorx.Throw(api.usrv.SignOut(util.GetToken(c)))
	response.Ctx(c).Ok()
}

// Del 注销账户
func (api *UserHandler) Del(c *gin.Context) {
	req := request.Bind(c, request.DelActReq{})
	dbUser := api.usrv.GetById(req.LoginUserId)
	if dbUser.Password != req.Password {
		response.Ctx(c).ServerI18n("password_error")
		return
	}
	api.usrv.Del(req.LoginUserId)
	errorx.Throw(api.usrv.SignOut(util.GetToken(c)))
	response.Ctx(c).Ok()
}

// report 心跳接口
func (api *UserHandler) report(c *gin.Context) {
	req := request.BindBase(c)
	api.usrv.Heartbeat(req.LoginUserId)
	response.Ctx(c).Ok()
}

func (api *UserHandler) getLanguage(ctx *gin.Context) {
	resp := []response.LanguageResp{
		{
			Name: "English",
			Code: "en",
		},
	}
	response.Ctx(ctx).OkJson(resp)
}

func (api *UserHandler) healthCheck(c *gin.Context) {
	response.Ctx(c).Ok()
}
