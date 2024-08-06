package request

type WebLoginReq struct {
	Email      string `form:"email"`
	Password   string `form:"password"`
	VerifyCode string `json:"verifyCode" form:"verifyCode"`
	CapToken   string `json:"cap_token" form:"cap_token"`
	Type       int    `form:"type"`
	SysDeviceVo
}

type WebSendVerifyRequest struct {
	Request
	GoogleToken string `form:"googleToken" json:"googleToken"`
	Email       string `form:"email" json:"email"`
	Type        string `form:"type" json:"type"`
	VerifyCode  string `form:"verifyCode" json:"verifyCode"`
	InviteCode  string `json:"inviteCode" form:"inviteCode"`
}

type WebInviteCodeReq struct {
	InviteCode string `json:"inviteCode" form:"inviteCode"`
}

type WebRegisterReq struct {
	Email      string `json:"email" form:"email"`
	VerifyCode string `json:"verifyCode" form:"verifyCode"`
	Password   string `json:"password" form:"password"`
	WebInviteCodeReq
}
