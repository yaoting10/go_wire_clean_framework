package request

type AddComplaintReq struct {
	Request
	VideoId uint   `form:"videoId"`
	Type    uint   `form:"type"`
	Content string `form:"content"`
}
