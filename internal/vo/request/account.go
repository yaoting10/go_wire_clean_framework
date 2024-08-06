package request

import "goboot/internal/vo"

type AccountHistoryReq struct {
	Request
	vo.Page
	Coin string `form:"coin" json:"coin"`
}

type TransferReq struct {
	Request
	DeviceId string `json:"deviceId" form:"deviceId"`
	Coin     string `json:"coin" form:"coin"`
	Amount   string `json:"amount" form:"amount"`
	Password string `json:"password" form:"password"`
}

type AccountSwapReq struct {
	Request
	Amount string `json:"amount" form:"amount"`
	Type   string `json:"type" form:"type"`
}

type AccountResetPwd struct {
	Request
	Password   string `form:"password"`
	VerifyCode string `json:"verifyCode" form:"verifyCode"`
}

type AccountNFTOutReq struct {
	Request
	TokenId  uint   `form:"tokenId"`
	DeviceId string `form:"deviceId"`
}

type AccountNFTInReq struct {
	Request
	TokenId  uint   `form:"tokenId"`
	DeviceId string `form:"deviceId"`
}

type AccountNFTSendReq struct {
	Request
	TokenId  uint   `form:"tokenId"`
	To       string `form:"to"`
	DeviceId string `form:"deviceId"`
}

type AccountNFTFeeReq struct {
	Request
	Type    int    `form:"type"`
	TokenId string `form:"tokenId"`
}

type AccountFlowReq struct {
	Request
	Coin   string `form:"coin"`
	Amount string `form:"amount"`
	Dir    int    `form:"dir"`
}

type AccountLimitReq struct {
	Request
	Coin   string `form:"coin"`
	Amount string `form:"amount"`
	UserId uint   `form:"userId"`
	To     string `form:"to"`
}

type AccountFixReq struct {
	Request
	DB string `form:"db"`
}

type AccountPreIdoReq struct {
	Request
	EthAmount string `form:"ethAmount"`
	Addr      string `form:"addr"`
}

type AccountIdoSucReq struct {
	Request
	EthAmount string `form:"ethAmount"`
	TxId      uint   `form:"txId"`
	TxHash    string `form:"txHash"`
}

type AccountIdoCancelReq struct {
	Request
	TxId uint `form:"txId"`
}

type AccountIdoUnlockReq struct {
	Request
	Addr string `form:"addr"`
}
