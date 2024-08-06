package response

type InviteUserResp struct {
	UserId       uint    `json:"userId"`
	NickName     string  `json:"nickName"`
	Avatar       string  `json:"avatar"`
	Email        string  `json:"email"`
	Telegram     string  `json:"telegram"`
	Discord      string  `json:"discord"`
	VipLevel     int     `json:"vipLevel"`
	ParentReward float64 `json:"parentReward"`
	InviteNum    int     `json:"inviteNum"`
	TeamNum      int     `json:"teamNum"`
	IsOnline     bool    `json:"isOnline"`
}

type InviteCodeResp struct {
	Code   string `json:"code"`
	Status int8   `json:"status"`
}
