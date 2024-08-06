package repository

import "goboot/internal/repository/repo"

type SysDeviceRepository struct {
	*repo.Repository
}

func NewSysDeviceRepository(r *repo.Repository) *SysDeviceRepository {
	return &SysDeviceRepository{Repository: r}
}
