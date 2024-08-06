package request

import "goboot/internal/vo"

type TradeCancelReq struct {
	Request
	TradeId uint `json:"tradeId" form:"tradeId"`
}

type TradeDealReq struct {
	Request
	TradeId uint `json:"tradeId" form:"tradeId"`
}

type TradeListReq struct {
	Request
	vo.Page
	Type      int8   `form:"type"`
	Lvl       int8   `form:"lvl"`
	Sort      int8   `form:"sort"`
	Dir       string `form:"dir"`
	Gender    int8   `form:"gender"`
	StarLvl   int8   `form:"starLvl"`
	LvlSort   int8   `form:"lvlSort"`
	EquipType int8   `form:"equipType"`
}

type CreateReq struct {
	Request
	Type     int8    `json:"type" form:"type"`
	DomainId uint    `json:"domainId" form:"domainId"`
	Price    float64 `json:"price" form:"price"`
}
