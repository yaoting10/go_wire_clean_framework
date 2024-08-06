package request

type VideoDiscusslikeRequest struct {
	Request
	VideoId   uint `form:"videoId" json:"videoId"`
	DiscussId uint `form:"discussId" json:"discussId"`
	Type      int  `form:"type" json:"type"`
}
