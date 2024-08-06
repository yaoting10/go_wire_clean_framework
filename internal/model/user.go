package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string     `gorm:"type:varchar(64);unique;not null;comment:用户名"`
	Password          string     `gorm:"type:varchar(128);not null;comment:用户登录密码"`
	GooAuthCode       string     `gorm:"type:varchar(128);default:'';comment:谷歌验证码"`
	AuthBind          int        `gorm:"type:tinyint;default:0;comment:是否绑定谷歌验证器"`
	UniqueNumber      string     `gorm:"type:varchar(30);unique;comment:编码"`
	NickName          string     `gorm:"type:varchar(64);default:'';comment:用户昵称"`
	Status            int        `gorm:"type:tinyint;default:1;comment:状态"`
	Avatar            string     `gorm:"type:varchar(256);default:'';comment:用户头像"`
	Gender            int        `gorm:"type:tinyint;default:0;comment:性别"`
	Phone             string     `gorm:"type:varchar(11);default:'';comment:用户手机号"`
	Email             string     `gorm:"type:varchar(64);unique;not null;comment:用户邮箱"`
	Banner            string     `gorm:"type:varchar(128);default:'';comment:背景图url"`
	Intro             string     `gorm:"type:varchar(512);default:'';comment:简介"`
	Name              string     `gorm:"type:varchar(64);default:'';comment:姓名"`
	BirthDay          *time.Time `gorm:"type:datetime;default:null;comment:生日"`
	Address           string     `gorm:"type:varchar(128);default:'';comment:所在地"`
	ParentId          uint       `gorm:"default:0;type:bigint;comment:上级id"`
	GrandparentId     uint       `gorm:"default:0;type:bigint;comment:'上上级id'"`
	InviteNum         int        `gorm:"type:int;default:0;comment:邀请人数"`
	ActiveTime        time.Time  `gorm:"type:datetime;default:'1970-01-01 00:00:00';comment:最近活跃时间"`
	TwitterClientId   string     `gorm:"type:varchar(128);default:'';comment:'twitter clientId'"`
	TwitterUserId     string     `gorm:"varchar(128);default:'';comment:twitter用户id"`
	TwitterUserName   string     `gorm:"varchar(128);default:'';comment:twitter用户name"`
	TwitterUserAvatar string     `gorm:"varchar(512);default:'';comment:twitter用户头像"`
	ParentMust        bool       `gorm:"type:tinyint;default:0;comment:'是否必需有父级'"`
	PetCount          int        `gorm:"type:integer;default:0;comment:'宠物数量'"`
}

const (
	UserStatusForbidden  = -2 // 禁用
	UserStatusLocked     = -1 // 锁定
	UserStatusRegistered = 1  // 已注册
	UserStatusVerified   = 2  // 已验证
)

const (
	UserGenderUnknown = iota
	UserGenderMale
	UserGenderFemale

	UserWithdrawPermit = 1
)

func (u *User) IsLocked() bool {
	return u.Status == UserStatusLocked
}

func (u *User) IsForbidden() bool {
	return u.Status == UserStatusForbidden
}
