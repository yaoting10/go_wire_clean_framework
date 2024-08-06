package request

import "goboot/internal/vo"

type InviteUserListReq struct {
	Request
	UserId uint `form:"userId"`
	vo.Page
}
