package service

import (
	"goboot/internal/model"
	"goboot/internal/repository"
	"goboot/internal/vo/request"
	"gorm.io/gorm"
)

type SysDeviceService struct {
	*Service
	sdrpo *repository.SysDeviceRepository
}

func NewSysDeviceService(s *Service, repo *repository.SysDeviceRepository) *SysDeviceService {
	return &SysDeviceService{
		Service: s,
		sdrpo:   repo,
	}
}

func (srv *SysDeviceService) GetLatestByUserId(userId uint) *model.SysDevice {
	var m *model.SysDevice
	r := srv.sdrpo.R.Where(model.SysDevice{UserId: userId}).Order("updated_at desc").Find(&m)
	if r.Error != nil || r.RowsAffected < 1 {
		return nil
	}
	return m
}

func (srv *SysDeviceService) GetDeviceInfoByUser(userId uint, vo *request.SysDeviceVo) *model.SysDevice {
	var r *gorm.DB
	var m *model.SysDevice
	if vo.DeviceId == "" {
		r = srv.sdrpo.R.Where(model.SysDevice{UserId: userId}).Find(&m) // 设备id为空，按照 user_id 查询
	} else {
		r = srv.sdrpo.R.Where(model.SysDevice{UserId: userId, DeviceId: vo.DeviceId}).Find(&m)
	}
	if r.Error != nil || r.RowsAffected < 1 {
		return nil
	}
	return m
}

func (srv *SysDeviceService) CanPreUpdate(deviceId string) bool {
	di := &model.SysDevice{}
	r := srv.sdrpo.R.Model(di).Where("device_id = ?", deviceId).Find(&di)
	if r.Error != nil || r.RowsAffected < 1 {
		return false
	}
	return di.PreUpdate
}

func (srv *SysDeviceService) Create(m *model.SysDevice) error {
	r := srv.sdrpo.W.Create(&m)
	if r.Error != nil {
		return r.Error
	}
	return nil
}
