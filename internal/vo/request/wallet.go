package request

import (
	"goboot/internal/vo"
)

type WalletInfoReq struct {
	Request
	DeviceId string `json:"deviceId" form:"deviceId"`
	Password string `json:"password" form:"password"`
}

type WalletHistoryReq struct {
	Request
	vo.Page
	DeviceId string `json:"deviceId" form:"deviceId"`
	Coin     string `json:"coin" form:"coin"`
}

type WalletImportReq struct {
	Request
	DeviceId string `json:"deviceId" form:"deviceId"`
	Key      string `json:"key" form:"key"`
	Password string `json:"password" form:"password"`
	Name     string `json:"name" form:"name"`
}

type WalletRestoreReq struct {
	Request
	DeviceId string `form:"deviceId"`
	Key      string `form:"key"`
	Password string `form:"password"`
	Name     string `json:"name" form:"name"`
}

type WalletGenReq struct {
	Request
	DeviceId string `json:"deviceId" form:"deviceId"`
	Mnemonic string `json:"mnemonic" form:"mnemonic"`
	Name     string `json:"name" form:"name"`
	Password string `json:"password" form:"password"`
}

type WalletResetPwdReq struct {
	Request
	DeviceId string `json:"deviceId" form:"deviceId"`
	Key      string `json:"key" form:"key"`
	Password string `json:"password" form:"password"`
}

type WalletResetNameReq struct {
	Request
	DeviceId string `json:"deviceId" form:"deviceId"`
	Name     string `json:"name" form:"name"`
}

type WalletSendReq struct {
	Request
	DeviceId string `json:"deviceId" form:"deviceId"`
	To       string `json:"to" form:"to"`
	Amount   string `json:"amount" form:"amount"`
	Coin     string `json:"coin" form:"coin"`
	Password string `json:"password" form:"password"`
}

type VerifyWalletReq struct {
	Request
	Type     int8   `form:"type"`
	DeviceId string `form:"deviceId"`
	Key      string `form:"key"`
	Password string `form:"password"`
}

type WalletSuggestGasReq struct {
	Request
	Type   int    `form:"type"`
	Coin   string `form:"coin"`
	Amount string `form:"amount"`
}

type WalletMaxReq struct {
	Request
	Coin   string `form:"coin"`
	Amount string `form:"amount"`
	Type   int8   `form:"type"`
}
