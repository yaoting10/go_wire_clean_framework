package service

import (
	"github.com/gophero/goal/assert"
	"goboot/internal/model"
	"goboot/internal/repository"
)

type AppVersionService struct {
	*Service
	svrepo *repository.SysVersionRepository
}

func NewAppVersionService(s *Service, repo *repository.SysVersionRepository) *AppVersionService {
	return &AppVersionService{
		Service: s,
		svrepo:  repo,
	}
}

func (srv *AppVersionService) GetLatestReleaseVersion(name string, platform string) *model.SysAppVersion {
	assert.True(name != "" && platform != "")
	return srv.svrepo.GetLatestReleaseVersion(name, platform)
}

func (srv *AppVersionService) GetLatestVersion(name string, platform string) *model.SysAppVersion {
	assert.True(name != "" && platform != "")
	return srv.svrepo.GetLatestVersion(name, platform)
}
