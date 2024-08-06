package repository

import (
	"goboot/internal/model"
	"goboot/internal/repository/repo"

	"gorm.io/gorm"
)

type AccountHistoryRepository struct {
	*repo.Repository
}

func NewAccountHistoryRepository(r *repo.Repository) *AccountHistoryRepository {
	return &AccountHistoryRepository{Repository: r}
}

func (repo *AccountHistoryRepository) NewModel(userId uint) model.IAccountHistory {
	return model.NewAccountHistoryModel(userId, repo.L)
}

func (repo *AccountHistoryRepository) Record(ah model.IAccountHistory, account *model.Account, tradeCoinIncr float64, mainCoinIncr float64, lockMainCoinIncr float64, secCoinIncr float64, lockSecCoinIncr float64, typeVal int8, tx *gorm.DB) (model.IAccountHistory, error) {
	return ah.Record(ah, account, tradeCoinIncr, mainCoinIncr, lockMainCoinIncr, secCoinIncr, lockSecCoinIncr, typeVal, tx)
}
