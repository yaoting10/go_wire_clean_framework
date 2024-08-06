package request

type IdoRequest struct {
	Request
	NftId uint `form:"nftId"`
}
