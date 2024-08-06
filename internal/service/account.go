package service

import (
	"fmt"
	"goboot/internal/model"
	"goboot/internal/repository"
	"goboot/internal/vo"
	"goboot/internal/vo/response"

	"github.com/gophero/goal/assert"
	"github.com/gophero/goal/errorx"
	"github.com/gophero/goal/gormx"
	"github.com/gophero/goal/mathx"
	"gorm.io/gorm"
)

type AccountService struct {
	*Service
	repo    *repository.AccountRepository
	sssrepo *repository.SysSettingRepository
}

func NewAccountService(service *Service, repo *repository.AccountRepository, sssrepo *repository.SysSettingRepository) *AccountService {
	return &AccountService{
		Service: service,
		repo:    repo,
		sssrepo: sssrepo,
	}
}

func (srv *AccountService) AddAmountOfTypeWithHistory(act *model.Account, incr float64, coinType int8, typ int8, db *gorm.DB) (*model.Account, model.IAccountHistory, error) {
	assert.True(incr >= 0)
	assert.True(act != nil)
	assert.True(act.ID > 0)
	assert.True(act.UserId > 0)
	var err error
	var ah model.IAccountHistory
	err = db.Transaction(func(tx *gorm.DB) error {
		act, err = srv.addAmount(act, incr, coinType, tx)
		if err != nil {
			return err
		}
		ah, err = srv.repo.RecordHistory(act, incr, coinType, typ, tx)
		return err
	})
	return act, ah, err
}

func (srv *AccountService) AddAmountOfTypeWithHistoryAndMark(act *model.Account, incr float64, coinType int8, typ int8, awardParent, awardGrandParent, fromChild, fromGrandChild float64, db *gorm.DB) (*model.Account, model.IAccountHistory, error) {
	assert.True(incr >= 0)
	assert.True(act != nil)
	assert.True(act.ID > 0)
	assert.True(act.UserId > 0)
	var err error
	var ah model.IAccountHistory
	err = db.Transaction(func(tx *gorm.DB) error {
		act, err = srv.addAmountAndAwardParent(act, incr, coinType, awardParent, awardGrandParent, fromChild, fromGrandChild, tx)
		if err != nil {
			return err
		}
		ah, err = srv.repo.RecordHistory(act, incr, coinType, typ, tx)
		return err
	})
	return act, ah, err
}

func (srv *AccountService) addAmountAndAwardParent(act *model.Account, incr float64, coinType int8, awardParent, awardGrandParent, fromChild, fromGrandChild float64, db *gorm.DB) (*model.Account, error) {
	assert.True(incr > 0)
	var r *gorm.DB

	switch coinType {
	case model.TradeCoinTokenType:
		act.TradeCoin = mathx.Add(act.TradeCoin, incr)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("trade_coin", "version").Updates(model.Account{TradeCoin: act.TradeCoin, Version: act.Version + 1})
	case model.SecCoinTokenType:
		act.SecCoin = mathx.Add(act.SecCoin, incr)
		// to award to superisor
		if awardParent > 0 {
			act.GiveFriendSecCoin += awardParent
		}
		if awardGrandParent > 0 {
			act.GiveSecFriendSecCoin += awardGrandParent
		}
		if fromChild > 0 {
			act.FriendGiveSecCoin += fromChild
		}
		if fromGrandChild > 0 {
			act.SecFriendGiveSecCoin += fromGrandChild
		}
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).
			Select("sec_coin", "version", "give_friend_sec_coin", "give_sec_friend_sec_coin", "friend_give_sec_coin", "sec_friend_give_sec_coin").
			Updates(model.Account{
				SecCoin: act.SecCoin, Version: act.Version + 1, GiveFriendSecCoin: act.GiveFriendSecCoin,
				GiveSecFriendSecCoin: act.GiveSecFriendSecCoin, FriendGiveSecCoin: act.FriendGiveSecCoin, SecFriendGiveSecCoin: act.SecFriendGiveSecCoin,
			})
	case model.MainCoinTokenType:
		act.MainCoin = mathx.Add(act.MainCoin, incr)
		// to award to superisor
		if awardParent > 0 {
			act.GiveFriendMainCoin += awardParent
		}
		if awardGrandParent > 0 {
			act.GiveSecFriendMainCoin += awardGrandParent
		}
		if fromChild > 0 {
			act.FriendGiveMainCoin += fromChild
		}
		if fromGrandChild > 0 {
			act.SecFriendGiveMainCoin += fromGrandChild
		}
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).
			Select("main_coin", "version", "give_friend_main_coin", "give_sec_friend_main_coin", "friend_give_main_coin", "sec_friend_give_main_coin").
			Updates(model.Account{
				MainCoin: act.MainCoin, Version: act.Version + 1, GiveFriendMainCoin: act.GiveFriendMainCoin,
				GiveSecFriendMainCoin: act.GiveSecFriendMainCoin, FriendGiveMainCoin: act.FriendGiveMainCoin, SecFriendGiveMainCoin: act.SecFriendGiveMainCoin,
			})
	case model.LockMainCoinType:
		act.LockMainCoin = mathx.Add(act.LockMainCoin, incr)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("lock_main_coin", "version").Updates(model.Account{LockMainCoin: act.LockMainCoin, Version: act.Version + 1})
	case model.LockSecCoinType:
		act.LockSecCoin = mathx.Add(act.LockSecCoin, incr)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("lock_sec_coin", "version").Updates(model.Account{LockSecCoin: act.LockSecCoin, Version: act.Version + 1})
	}

	if r.Error != nil || r.RowsAffected <= 0 {
		return nil, gormx.UpdateError
	}
	return act, db.Error
}

func (srv *AccountService) addAmount(act *model.Account, incr float64, coinType int8, db *gorm.DB) (*model.Account, error) {
	assert.True(incr > 0)
	var r *gorm.DB
	switch coinType {
	case model.TradeCoinTokenType:
		act.TradeCoin = mathx.Add(act.TradeCoin, incr)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("trade_coin", "version").Updates(model.Account{TradeCoin: act.TradeCoin, Version: act.Version + 1})
	case model.SecCoinTokenType:
		act.SecCoin = mathx.Add(act.SecCoin, incr)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("sec_coin", "version").Updates(model.Account{SecCoin: act.SecCoin, Version: act.Version + 1})
	case model.MainCoinTokenType:
		act.MainCoin = mathx.Add(act.MainCoin, incr)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("main_coin", "version").Updates(model.Account{MainCoin: act.MainCoin, Version: act.Version + 1})
	case model.LockMainCoinType:
		act.LockMainCoin = mathx.Add(act.LockMainCoin, incr)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("lock_main_coin", "version").Updates(model.Account{LockMainCoin: act.LockMainCoin, Version: act.Version + 1})
	}

	if r.Error != nil || r.RowsAffected <= 0 {
		return nil, gormx.UpdateError
	}
	return act, db.Error
}

func (srv *AccountService) SubAmountOfTypeWithHistory(act *model.Account, amount float64, coinType int8, typ int8, db *gorm.DB) (*model.Account, model.IAccountHistory, error) {
	assert.True(amount > 0)
	assert.True(act != nil)
	assert.True(act.ID > 0)
	assert.True(act.UserId > 0)

	var err error
	var ah model.IAccountHistory
	err = db.Transaction(func(tx *gorm.DB) error {
		act, err = srv.subAmountOfType(act, amount, coinType, tx)
		if err != nil {
			return err
		}
		ah, err = srv.repo.RecordHistory(act, -amount, coinType, typ, tx)
		return err
	})
	return act, ah, err
}

func (srv *AccountService) subAmountOfType(act *model.Account, amount float64, coinType int8, db *gorm.DB) (*model.Account, error) {
	var r *gorm.DB
	switch coinType {
	case model.TradeCoinTokenType:
		if act.TradeCoin < amount {
			return nil, errorx.NewPreferredErrf("insufficient balance")
		}
		act.SecCoin = mathx.Sub(act.SecCoin, amount)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("sec_coin", "version").Updates(model.Account{SecCoin: act.SecCoin, Version: act.Version + 1})
	case model.MainCoinTokenType:
		if act.MainCoin < amount {
			return nil, errorx.NewPreferredErrf("insufficient balance")
		}
		act.MainCoin = mathx.Sub(act.MainCoin, amount)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("main_coin", "version").Updates(model.Account{MainCoin: act.MainCoin, Version: act.Version + 1})
	case model.LockMainCoinType:
		if act.LockMainCoin < amount {
			return nil, errorx.NewPreferredErrf("insufficient balance")
		}
		act.LockMainCoin = mathx.Sub(act.LockMainCoin, amount)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("lock_main_coin", "version").Updates(model.Account{LockMainCoin: act.LockMainCoin, Version: act.Version + 1})
	case model.LockSecCoinType:
		if act.LockSecCoin < amount {
			return nil, errorx.NewPreferredErrf("insufficient balance")
		}
		act.LockSecCoin = mathx.Sub(act.LockSecCoin, amount)
		r = db.Model(act).Where("user_id = ? and version = ?", act.UserId, act.Version).Select("lock_sec_coin", "version").Updates(model.Account{LockSecCoin: act.LockSecCoin, Version: act.Version + 1})
	default:
		return nil, errorx.NewPreferredErrf("coin type is error")
	}
	if r.Error != nil || r.RowsAffected <= 0 {
		return nil, gormx.UpdateError
	}
	return act, r.Error
}

func (srv *AccountService) GetByUserId(userId uint) (account *model.Account) {
	r := srv.R.R.Where("user_id=?", userId).First(&account)
	if r.Error != nil || r.RowsAffected == 0 {
		return nil
	}
	return account
}

func (srv *AccountService) GetByUserIdTx(tx *gorm.DB, userId uint) (account *model.Account) {
	r := tx.Where("user_id=?", userId).First(&account)
	if r.Error != nil || r.RowsAffected == 0 {
		return nil
	}
	return account
}

func (srv *AccountService) GetHistory(userId uint, typeVal int, page vo.Page) []response.AccountHistoryVo {
	acHistories := make([]response.AccountHistoryVo, 0)
	sql := ""
	switch typeVal {
	case model.TradeCoinTokenType:
		sql = "SELECT (trade_coin_incr) AS amount, trade_coin AS balance, type, created_at FROM %s WHERE user_id=? AND trade_coin_incr <> 0"
	case model.MainCoinTokenType:
		sql = "SELECT (main_coin_incr) AS amount, main_coin AS balance, type, created_at FROM %s WHERE user_id=? AND main_coin_incr <> 0"
	case model.LockMainCoinType:
		sql = "SELECT (lock_main_coin_incr) AS amount, lock_main_coin AS balance, type, created_at FROM %s WHERE user_id=? AND lock_main_coin_incr <> 0"
	case model.SecCoinTokenType:
		sql = "SELECT (sec_coin_incr) AS amount, sec_coin AS balance, type, created_at FROM %s WHERE user_id=? AND sec_coin_incr <> 0"
	case model.LockSecCoinType:
		sql = "SELECT (lock_sec_coin_incr) AS amount, lock_sec_coin AS balance, type, created_at FROM %s WHERE user_id=? AND lock_sec_coin_incr <> 0"
	default:
		return acHistories
	}
	sql = fmt.Sprintf(sql, "account_histories")
	srv.repo.R.Raw(sql+" ORDER BY created_at DESC LIMIT ?,?", userId, (page.PageNum-1)*page.PageSize, page.PageSize).Scan(&acHistories)
	return acHistories
}
