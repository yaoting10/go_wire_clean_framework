package model

import (
	"gorm.io/gorm"
)

type SysDevice struct {
	gorm.Model
	DeviceId          string `gorm:"type:VARCHAR(255);not null;comment:'设备id'" json:"deviceId" form:"deviceId"`
	Platform          string `gorm:"type:VARCHAR(32);default:'';comment:'平台 ios,android'" json:"platform" form:"platform"`
	Channel           string `gorm:"type:VARCHAR(32);default:'';comment:'app渠道'" json:"channel" form:"channel"`
	VersionInfo       string `gorm:"type:VARCHAR(50);default:'';comment:'版本号'" json:"VersionInfo" form:"VersionInfo"`
	UserId            uint   `gorm:"type:bigint;default:0;comment:'用户id'" json:"userId" form:"userId"`
	DeviceModel       string `gorm:"type:varchar(512);default:'';comment:'设备型号'" json:"deviceModel" form:"deviceModel"`
	DeviceBrand       string `gorm:"type:varchar(512);default:'';comment:'设备品牌'" json:"deviceBrand" form:"deviceBrand"`
	DeviceType        string `gorm:"type:varchar(512);default:'';comment:'设备类型'" json:"deviceType" form:"deviceType"`
	DeviceOrientation string `gorm:"type:varchar(512);default:'';comment:'设备方向'" json:"deviceOrientation" form:"deviceOrientation"`
	DevicePixelRatio  string `gorm:"type:varchar(512);default:'';comment:'设备像素比'" json:"devicePixelRatio" form:"devicePixelRatio"`
	System            string `gorm:"type:varchar(512);default:'';comment:'操作系统及版本'" json:"system" form:"system"`
	PreUpdate         bool   `gorm:"type:tinyint;default:0;comment:'预更新'" json:"preUpdate" form:"preUpdate"`
	SysLang           string `gorm:"type:varchar(32);default:'';comment:'操作系统语言'" json:"sysLang" form:"sysLang"`
}
