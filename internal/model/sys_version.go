package model

import (
	"gorm.io/gorm"
)

const (
	AppPlatformIos     = "ios"
	AppPlatformAndroid = "android"
)

const AppNewStatus = 1

type SysAppVersion struct {
	gorm.Model
	Name        string `gorm:"type:VARCHAR(32);default:'';comment:名称" json:"name" form:"name"`
	Platform    string `gorm:"type:VARCHAR(32);default:'';comment:平台 ios,android" json:"platform" form:"platform"`
	Url         string `gorm:"type:VARCHAR(255);default:'';comment:地址" json:"url" form:"url"`
	QrUrl       string `gorm:"type:VARCHAR(255);default:'';comment:预览图" json:"qrUrl" form:"qrUrl"`
	VersionInfo string `gorm:"type:VARCHAR(50);default:'';comment:版本号" json:"VersionInfo" form:"VersionInfo"`
	Content     string `gorm:"type:VARCHAR(4096);default:'';comment:名称" json:"content" form:"content"`
	Status      int8   `gorm:"type:smallint(2);default:1;comment:状态 ：1:最新状态，2:历史版本" json:"status" form:"status"`
	Force       bool   `gorm:"type:tinyint;;comment:是否强制更新 0 否，1 是" json:"force" form:"force"`
	Remark      string `gorm:"type:VARCHAR(255);default:'';comment:备注" json:"remark" form:"remark"`
}
