package request

import (
	"goboot/internal/vo"
)

type UserSubscribesRequest struct {
	Request
	vo.Page
	//SrcUserId uint `form:"srcUserId" json:"srcUserId"`
	//DestUserId uint `form:"destUserId" json:"destUserId"`
	UserId uint `form:"userId" json:"userId"`
}
