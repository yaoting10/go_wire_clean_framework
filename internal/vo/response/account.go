package response

import "time"

type AccountDetailResp struct {
	Amount   string `json:"amount"`
	Price    string `json:"price"`
	USD      string `json:"USD"`
	Type     int    `json:"type"`
	CoinName string `json:"coinName"`
	Icon     string `json:"icon"`
}

type BaseAccountHistoryVo struct {
	Type   int8    `json:"type"`
	Amount float64 `json:"amount"`
	Time   string  `json:"time"`
	Label  string  `json:"label"`
}

type AccountHistoryVo struct {
	BaseAccountHistoryVo
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"-"`
}

type AccountSwapDetailResp struct {
	CoinName string  `json:"coinName"`
	Price    string  `json:"price"`
	Icon     string  `json:"icon"`
	Amount   float64 `json:"amount"`
}

type AccountSwapResp struct {
	CoinName string  `json:"coinName"`
	Amount   float64 `json:"amount"`
}
