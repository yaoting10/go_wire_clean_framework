package request

type LoginReq struct {
	Email      string `form:"email"`
	Password   string `form:"password"`
	VerifyCode string `json:"verifyCode" form:"verifyCode"`
	Type       int    `form:"type"`
	SysDeviceVo
}

type SendVerifyRequest struct {
	Request
	Email      string ` form:"email" json:"email"`
	Type       string `form:"type" json:"type"`
	VerifyCode string `form:"verifyCode" json:"verifyCode"`
	InviteCode string `json:"inviteCode" form:"inviteCode"`
}

type AppVersionReq struct {
	Request
	DeviceId string `json:"deviceId" form:"deviceId"`
	Platform int    `json:"platform" form:"platform"`
	Version  string `json:"version" form:"version"`
}

type ResetPwdReq struct {
	Email      string `json:"email" form:"email"`
	Password   string `json:"password" form:"password"`
	VerifyCode string `json:"verifyCode" form:"verifyCode"`
}

type SystemSettingRequest struct {
	Request
	SysKey string `form:"sysKey" json:"sysKey"`
}
