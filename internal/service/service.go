package service

import (
	"github.com/gophero/goal/logx"
	"github.com/gophero/goal/uuid"
	"goboot/internal/config"
	"goboot/internal/repository/repo"
)

type Service struct {
	R   *repo.Repository
	L   *logx.Logger
	sid *uuid.Sid
	C   config.Conf
}

func NewService(repo *repo.Repository, logger *logx.Logger, sid *uuid.Sid, conf config.Conf) *Service {
	return &Service{
		R:   repo,
		L:   logger,
		sid: sid,
		C:   conf,
	}
}
