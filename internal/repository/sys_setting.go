package repository

import (
	"goboot/internal/model"
	"goboot/internal/repository/repo"
)

type SysSettingRepository struct {
	*repo.Repository
}

func NewSysSettingRepository(r *repo.Repository) *SysSettingRepository {
	return &SysSettingRepository{Repository: r}
}

func (repo *SysSettingRepository) GetSystemSetting(id int64) *model.SysSetting {
	ss := new(model.SysSetting)
	r := repo.R.Find(&ss, id)
	if r.Error != nil {
		return nil
	}
	if r.RowsAffected > 0 {
		return ss
	}
	return nil
}

func (repo *SysSettingRepository) GetSystemSettingByKey(key string) *model.SysSetting {
	ss := new(model.SysSetting)
	r := repo.R.Where(&model.SysSetting{SysKey: key}).Find(ss)
	if r.Error != nil {
		return nil
	}
	if r.RowsAffected > 0 {
		return ss
	}
	return nil
}

func (repo *SysSettingRepository) GetSystemSettingsValByKeys(keys ...string) map[string]string {
	ss := make([]model.SysSetting, 0)
	r := repo.R.Where("sys_key in ?", keys).Find(&ss)
	if r.Error != nil || len(ss) < len(keys) {
		return nil
	}
	vals := make(map[string]string, 0)
	for i, v := range ss {
		for _, key := range keys {
			if ss[i].SysKey == key {
				vals[key] = v.SysValue
			}
		}
	}
	return vals
}
