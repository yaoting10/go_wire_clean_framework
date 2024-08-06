package request

type BlindBox struct {
	Request
	BoxId uint `form:"boxId" json:"boxId"`
}
