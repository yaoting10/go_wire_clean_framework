package model

import (
	"gorm.io/gorm"
)

type SysSetting struct {
	gorm.Model
	SysKey      string `gorm:"type:varchar(32);not null;" form:"sysKey"`
	Name        string `gorm:"type:varchar(32);" form:"name"`
	SysValue    string `gorm:"type:varchar(2048);not null;" form:"sysValue"`
	Description string `gorm:"type:varchar(128)" form:"description"`
}
