package service

import (
	"github.com/gophero/goal/assert"
	"github.com/gophero/goal/errorx"
	"goboot/internal/model"
	"goboot/internal/repository"
)

type SysSettingService struct {
	sysRepo *repository.SysSettingRepository
	*Service
}

func NewSysSettingService(service *Service, sysRepo *repository.SysSettingRepository) *SysSettingService {
	return &SysSettingService{
		sysRepo: sysRepo,
		Service: service,
	}
}

func (srv *SysSettingService) GetSystemSetting(id int64) *model.SysSetting {
	assert.True(id > 0)
	return srv.sysRepo.GetSystemSetting(id)
}

func (srv *SysSettingService) GetSystemSettingByKey(key string) *model.SysSetting {
	assert.True(key != "")
	return srv.sysRepo.GetSystemSettingByKey(key)
}

func (srv *SysSettingService) GetSystemSettingsValByKeys(keys ...string) map[string]string {
	assert.True(len(keys) > 0)
	return srv.sysRepo.GetSystemSettingsValByKeys(keys...)
}

func (srv *SysSettingService) UpdateByKey(key string, value string) (int64, error) {
	if key == "" || value == "" {
		return 0, errorx.New("missing parameters")
	}
	r := srv.sysRepo.W.Model(&model.SysSetting{}).Where("sys_key", key).Update("sys_value", value)
	return r.RowsAffected, nil
}
