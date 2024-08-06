package request

import (
	"goboot/internal/vo"
)

type VideodiscusRequest struct {
	Request
	vo.Page
	Id        uint   `form:"id" json:"id"`
	VideoId   uint   `form:"videoId" json:"videoId"`
	DiscussId uint   `form:"discussId" json:"discussId"`
	UserId    uint   `form:"userId" json:"userId"`
	Content   string `form:"content" json:"content"`
}
