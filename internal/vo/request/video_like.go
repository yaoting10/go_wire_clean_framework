package request

import (
	"goboot/internal/vo"
)

type VideolikeRequest struct {
	Request
	vo.Page
	VideoId uint `form:"videoId" json:"videoId"`
}
