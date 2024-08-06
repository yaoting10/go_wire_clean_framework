package repository

import (
	"errors"
	"fmt"
	"goboot/internal/model"
	"goboot/internal/repository/repo"

	"gorm.io/gorm"
)

type AccountRepository struct {
	*repo.Repository
	hisrepo *AccountHistoryRepository
}

func NewAccountRepository(r *repo.Repository, hisrepo *AccountHistoryRepository) *AccountRepository {
	return &AccountRepository{Repository: r, hisrepo: hisrepo}
}

func (repo *AccountRepository) RecordHistory(m *model.Account, incr float64, coinType int8, typeVal int8, tx *gorm.DB) (model.IAccountHistory, error) {
	userId := m.UserId
	if userId <= 0 {
		return nil, fmt.Errorf("invalid userId of account")
	}
	h := repo.hisrepo.NewModel(userId)
	h.SetAccountId(m.ID)
	h.SetUserId(m.UserId)
	h.SetTradeCoin(m.TradeCoin)
	switch coinType {
	case model.TradeCoinTokenType:
		h.SetTradeCoinIncr(incr)
	case model.MainCoinTokenType:
		h.SetMainCoinIncr(incr)
	case model.LockMainCoinType:
		h.SetLockMainCoinIncr(incr)
	case model.SecCoinTokenType:
		h.SetSecCoinIncr(incr)
	case model.LockSecCoinType:
		h.SetLockSecCoinIncr(incr)
	}
	h.SetMainCoin(m.MainCoin)
	h.SetSecCoin(m.SecCoin)
	h.SetLockMainCoin(m.LockMainCoin)
	h.SetLockSecCoin(m.LockSecCoin)
	h.SetType(typeVal)
	r := tx.Save(h)
	if r.Error != nil {
		return nil, r.Error
	}
	if r.RowsAffected < 1 {
		return nil, errors.New("insert error")
	}
	return h, nil
}
