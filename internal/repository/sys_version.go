package repository

import (
	"goboot/internal/model"
	"goboot/internal/repository/repo"
)

type SysVersionRepository struct {
	*repo.Repository
}

func NewSysVersionRepository(r *repo.Repository) *SysVersionRepository {
	return &SysVersionRepository{Repository: r}
}

func (repo *SysVersionRepository) GetLatestReleaseVersion(name string, platform string) *model.SysAppVersion {
	var m *model.SysAppVersion
	r := repo.R.Where(model.SysAppVersion{Name: name, Platform: platform, Status: model.AppNewStatus}).Order("id desc").Find(&m)
	if r.Error != nil || r.RowsAffected < 0 {
		return nil
	}
	return m
}

func (repo *SysVersionRepository) GetLatestVersion(name string, platform string) *model.SysAppVersion {
	var m *model.SysAppVersion
	r := repo.R.Where(model.SysAppVersion{Name: name, Platform: platform, Status: model.AppNewStatus}).Order("id desc").Find(&m)
	if r.Error != nil || r.RowsAffected < 0 {
		return nil
	}
	return m
}
