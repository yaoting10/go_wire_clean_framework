package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/errorx"
	"goboot/internal/service"
	"goboot/internal/vo/request"
	"goboot/internal/vo/response"
)

type SysSettingHandler struct {
	*Handler
	sysSettingService *service.SysSettingService
}

func NewSysSettingHandler(handler *Handler, sysSettingService *service.SysSettingService) *SysSettingHandler {
	return &SysSettingHandler{
		Handler:           handler,
		sysSettingService: sysSettingService,
	}
}

func (api *SysSettingHandler) RegisterRouter(r *gin.Engine) {
	g := api.Group("/system/setting", r)
	{
		g.GET("", api.getSystemSetting)
	}
}

type SystemSettingRequest struct {
	request.Request
	SysKey string `form:"sysKey" json:"sysKey"`
}

type SettingRewardReq struct {
	request.Request
	Emails  string `form:"emails" json:"emails"`
	Amounts string `form:"amounts" json:"amounts"`
}

func (api *SysSettingHandler) getSystemSetting(c *gin.Context) {
	req := &SystemSettingRequest{}
	errorx.Throw(c.Bind(req))
	keys := []string{"privacy_url", "terms_url", "use_help", "award_rule", "free_nft_white_paper", "video_earnings_rule", "cute_replenishment_project"}
	settingMap := api.sysSettingService.GetSystemSettingsValByKeys(keys...)
	response.Ctx(c).OkJson(map[string]interface{}{
		"privacy":                  settingMap["privacy_url"],
		"terms":                    settingMap["terms_url"],
		"taskUseHelp":              settingMap["use_help"],
		"awardRule":                settingMap["award_rule"],
		"freeNftWhitePaper":        settingMap["free_nft_white_paper"],
		"videoEarningsRule":        settingMap["video_earnings_rule"],
		"cuteReplenishmentProject": settingMap["cute_replenishment_project"],
	})
}
