package request

type GroupDetail struct {
	Request
	Code string `form:"code"`
}

type CreateGroup struct {
	Request
	TempId uint   `form:"tempId"`
	NftIds string `form:"nftIds"`
}

type GroupCreateDetail struct {
	Request
	TempId uint `form:"tempId"`
}

type JoinGroup struct {
	Request
	Code   string `form:"code"`
	NftIds string `form:"nftIds"`
}

type GroupAward struct {
	Request
	GroupId uint `form:"groupId"`
}

type CloseRequest struct {
	Request
	GroupId uint `form:"groupId"`
}
