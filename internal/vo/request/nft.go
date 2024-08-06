package request

type Nft struct {
	Request
	Id uint `form:"id" json:"id"`
}

type NftUploadLog struct {
	Request
	NftId uint   `json:"nftId" form:"nftId"`
	Desc  string `json:"desc" form:"desc"`
}

type NftAttrAdd struct {
	Request
	Id  uint `form:"id"`  // id
	Str int  `form:"str"` // 增加的力量
	Dex int  `form:"dex"` // 增加的敏捷
	Int int  `form:"int"` // 增加的智力
}

type NftEquipSave struct {
	Request
	NftId uint   `form:"nftId"` // nft id
	EqIds string `form:"eqIds"` // 多个装备id用英文逗号分隔
}

type NftEquipRequest struct {
	Request
	NftId uint `form:"nftId"`
	EqId  uint `form:"eqId"`
}

type NftMintParents struct {
	Request
	Fid uint `form:"fid"`
	Sid uint `form:"sid"`
}

type UnlockSlot struct {
	Request
	Type  int  `form:"type"`
	NftId uint `form:"nftId"`
}

type NftUnlock struct {
	Request
	Type  int  `form:"type"`
	NftId uint `form:"id"`
}

type NftListRequest struct {
	Request
	Type     int    `form:"type"`
	Gender   int    `form:"gender"`
	TokenIds string `form:"tokenIds"`
}
