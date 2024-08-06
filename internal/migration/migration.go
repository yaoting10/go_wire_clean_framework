package migration

import (
	"github.com/gophero/goal/logx"
	"go.uber.org/zap"
	"goboot/internal/model"
	"goboot/internal/repository/repo"
	"gorm.io/gorm"
)

type Migrate struct {
	db  *gorm.DB
	log *logx.Logger
}

func NewMigrate(db repo.WDB, log *logx.Logger) *Migrate {
	return &Migrate{
		db:  db,
		log: log,
	}
}
func (m *Migrate) Run() {
	m.log.Info("Auto migrating normal models...")
	if err := m.db.AutoMigrate(
		// 添加model
		&model.SysAccount{},
		&model.SysDevice{},
		&model.SysSetting{},
		&model.SysAppVersion{},
		&model.User{},
		&model.UserInviteCodeGen{},
		&model.UserInviteCode{},
		&model.UserSubscribe{},
		&model.UserTeam{},
	); err != nil {
		m.log.Error("user migrate error", zap.Error(err))
		return
	}
	m.log.Info("AutoMigrate end")
}
