package model

import (
	"gorm.io/gorm"
)

// ==============================
// 用户邀请码
// ==============================

const (
	UserInviteCodeUsed   = 0
	UserInviteCodeEnable = 1
)

// UserInviteCodeGen 邀请码生成记录表
type UserInviteCodeGen struct {
	gorm.Model
	UserId  uint `gorm:"type:bigint;not null;comment:'用户id'"`
	Num     int  `gorm:"type:smallint;default:0;comment:'今日数量'"`
	Add     int  `gorm:"type:smallint;default:0;comment:'今日增加数量'"`
	NftNum  int  `gorm:"type:smallint;default:0;comment:'当前NFT对应的邀请码数量'"`
	VipNum  int  `gorm:"type:smallint;default:0;comment:'当前vip对应的邀请码数量'"`
	FreeNum int  `gorm:"type:smallint;default:0;comment:'幸运的免费用户的邀请码数量'"`
	GenEnd  bool `gorm:"type:tinyint;default:0;comment:'今日生成完成'"`
}

// UserInviteCode 用户邀请码
type UserInviteCode struct {
	gorm.Model
	UserId    uint   `gorm:"type:bigint;not null;comment:'用户id'"`
	ValidTime string `gorm:"type:varchar(20);default:'';comment:'时间'"`
	Code      string `gorm:"type:varchar(36);default:''；unique;comment:'邀请码'"`
	HasUsed   bool   `gorm:"type:tinyint;default:0;comment:'已经使用'"`
}
