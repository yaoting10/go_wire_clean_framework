package model

import (
	"gorm.io/gorm"
	"time"
)

type UserTeam struct {
	gorm.Model
	UserId        uint      `gorm:"default:0;type:bigint;unique;comment:'用户id'"`
	ParentPath    string    `gorm:"type:text;comment:父级路径"`
	TeamCount     int       `gorm:"default:0;type:bigint;comment:'团队数量'"`
	TeamCountTime time.Time `gorm:"type:datetime;comment:'最后统计团队数量时间'"`
}
