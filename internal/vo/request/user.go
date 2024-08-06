package request

type UpdateUserReq struct {
	Request
	NickName string `form:"nickName"`
	Gender   int    `form:"gender"`
	Intro    string `form:"intro"`
	BirthDay string `form:"birthDay"`
	Address  string `form:"address"`
	Telegram string `form:"telegram"`
	Discord  string `form:"discord"`
}

type InfoReq struct {
	Request
	UserId uint `form:"userId" json:"userId"`
}

type UpdPwdReq struct {
	Request
	OldPassword string `form:"oldPassword"`
	NewPassword string `form:"newPassword"`
}

type DelActReq struct {
	Request
	Password string `form:"password"`
}
