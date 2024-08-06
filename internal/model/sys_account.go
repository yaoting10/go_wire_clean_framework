package model

import (
	"gorm.io/gorm"
)

type SysAccount struct {
	gorm.Model
	Address        string `gorm:"default:'';type:varchar(100);comment:'地址'" json:"address" form:"address"`
	PrivateKey     string `gorm:"default:'';type:varchar(255);comment:'私钥'" json:"privateKey" form:"privateKey"`
	PublicKey      string `gorm:"default:'';type:varchar(255);comment:'公钥'" json:"publicKey" form:"publicKey"`
	Mnemonic       string `gorm:"default:'';type:varchar(255);comment:'助记词'" json:"mnemonic" form:"mnemonic"`
	RechargeState  int8   `gorm:"type:smallint(2);default:1;comment:'状态:充值状态'" json:"rechargeState" form:"rechargeState"`
	RechargeWeight int    `gorm:"type:smallint(2);default:0;comment:'状态:充值权重'" json:"rechargeWeight" form:"rechargeWeight"`
	WithdrawState  int8   `gorm:"type:smallint(2);default:1;comment:'状态:提现状态'" json:"withdrawState" form:"withdrawState"`
	WithdrawWeight int    `gorm:"type:smallint(2);default:0;comment:'状态:充值权重'" json:"withdrawWeight" form:"withdrawWeight"`
}

const (
	SysAccountDisable = iota
	SysAccountEnable
)
