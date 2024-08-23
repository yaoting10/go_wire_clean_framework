package repository

import (
	"github.com/gophero/goal/gormx"
	"goboot/internal/model"
	"goboot/internal/repository/repo"
	"gorm.io/gorm"
)

type UserRepository struct {
	*repo.Repository
}

func NewUserRepository(r *repo.Repository) *UserRepository {
	return &UserRepository{
		Repository: r,
	}
}

func (repo *UserRepository) Lock(m *model.User, db *gorm.DB) error {
	m.Status = model.UserStatusLocked
	_, err := repo.Update(m, db)
	return err
}

func (repo *UserRepository) Update(m *model.User, db *gorm.DB) (*model.User, error) {
	r := db.Save(m)
	return gormx.UpdateResult(r, m)
}

func (repo *UserRepository) Insert(m *model.User, tx *gorm.DB) (*model.User, error) {
	m.Status = model.UserStatusRegistered
	r := tx.Save(&m)
	return gormx.InsertResult(r, m)
}

func (repo *UserRepository) Verified(m *model.User, tx *gorm.DB) bool {
	r := tx.Where("id = ?", m.ID).Update("status", model.UserStatusVerified)
	if r.Error != nil {
		repo.L.Error("verified error: %v", r.Error)
		return false
	}
	return true
}

func (repo *UserRepository) GetUserByIds(db *gorm.DB, ids []uint) []*model.User {
	var users []*model.User
	if len(ids) == 0 {
		return users
	}
	r := db.Model(model.User{}).Where("id in ?", ids).Find(&users)
	if r.Error != nil {
		panic(r.Error)
	}
	return users
}

func (repo *UserRepository) Count(db *gorm.DB) int64 {
	var cnt int64
	db.Model(model.User{}).Count(&cnt)
	return cnt
}

func (repo *UserRepository) CountByStatus(db *gorm.DB, status ...interface{}) int64 {
	var cnt int64
	db.Model(model.User{}).Where("status in (?)", status...).Count(&cnt)
	return cnt
}

func (repo *UserRepository) GetByIdTx(db *gorm.DB, id uint) *model.User {
	var usr model.User
	db.Model(model.User{}).Where("id=?", id).Find(&usr)
	return &usr
}

// QueryNeedRefreshVip 查询需要降低 vip 等级的数据，按照刷新时间从低到高排序
// 可以刷新vip等级的条件：
// 1、未注销的账户
// 2、非内部账户
// 3、状态为正常的账户
// 4、上次刷新时间 calc_vip_time 与当前时间间隔超过 usrconf.UserConfs.VipCalcHours() 个小时
// 5、Vip等级大于0，等级为0的不可能再降级了
// func (repo *UserRepository) QueryNeedRefreshVip(tx *gorm.DB, cnt int64) ([]*model.User, error) {
// 	sql := `select * from users where deleted_at is null and is_internal = 0 and status = ? and vip_level > 0 and
//  TIMESTAMPDIFF(HOUR, calc_vip_time, now()) >= ? limit ?`
// 	r, err := repo.FindRawTx(tx, []*model.User{}, sql, model.UserStatusVerified, democfg.UserConfs.VipCalcHours(), cnt)
// 	if err != nil {
// 		return []*model.User{}, err
// 	}
// 	return r.([]*model.User), nil
// }
// func (repo *UserRepository) QueryNeedRefreshVipId(tx *gorm.DB, cnt int64) ([]uint, error) {
// 	sql := `select id from users where deleted_at is null and is_internal = 0 and status = ? and vip_level > 0 and
//  TIMESTAMPDIFF(HOUR, calc_vip_time, now()) >= ? limit ?`
// 	r, err := repo.FindRawTx(tx, []uint{}, sql, model.UserStatusVerified, democfg.UserConfs.VipCalcHours(), cnt)
// 	if err != nil {
// 		return []uint{}, err
// 	}
// 	return r.([]uint), nil
// }

// QueryChildLvlGTECnt 查询等级大于 p.VipLevel - 1 的下级数量
// Deprecated: 弃用，没有考虑活跃
//func (repo *UserRepository) QueryChildLvlGTECnt(tx *gorm.DB, p *model.User) (int64, error) {
//	cnt, err := repo.FindCntTx(tx, &model.User{}, cond.New("status = ? and is_internal = 0 and parent_id = ? and vip_level >= ?",
//		model.UserStatusVerified, p.ID, p.VipLevel-1))
//	if err != nil {
//		return 0, err
//	}
//	return cnt, nil
//}

// QueryActiveChild3rdLvl 查询所有下级中等级排名第三的活跃用户的等级，上级的 vipLvl = 排名第三的下级vipLvl + 1，如果不足三个下级，返回 -1
// func (repo *UserRepository) QueryActiveChild3rdLvl(tx *gorm.DB, userId uint) (int, error) {
// 	const limit = 3
// 	sql := `select * from users where deleted_at is null and is_internal = 0 and status = ? and parent_id = ?
//                     and TIMESTAMPDIFF(HOUR, active_time, now()) < ? order by vip_level desc limit ?`
// 	var us []*model.User
// 	_, err := repo.FindRawTx(tx, &us, sql, model.UserStatusVerified, userId, democfg.UserConfs.VipActiveHours(), limit)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if len(us) >= limit {
// 		return us[limit-1].VipLevel, nil
// 	}
// 	return -1, nil
// }

func (repo *UserRepository) UpdateTx(tx *gorm.DB, m *model.User, field string, val any) error {
	r := tx.Model(&model.User{}).Where("id=?", m.ID).Update(field, val)
	if r.Error != nil {
		return tx.Error
	}
	return nil
}
