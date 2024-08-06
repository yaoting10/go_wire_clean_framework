package model

import "gorm.io/gorm"

type UserSubscribe struct {
	gorm.Model
	SrcUserId  uint `gorm:"type:bigint;not null;comment:用户id"`
	DestUserId uint `gorm:"type:bigint;not null;comment:用户id"`
	Status     int  `gorm:"type:tinyint;default:1;comment:状态"`
}
