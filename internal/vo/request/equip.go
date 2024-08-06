package request

import "goboot/internal/vo"

type EquipSynthBaseReq struct {
	Request
	Id uint `form:"id"`
}

type EquipBreakdownReq struct {
	EquipSynthBaseReq `json:"synthBaseReq"`
	Ids               string `form:"ids" json:"ids"`
}

type EquipListReq struct {
	Request
	Type   int8 `form:"type"`
	Status int8 `form:"status"`
	vo.Page
}
