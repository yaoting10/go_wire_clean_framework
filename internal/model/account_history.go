package model

import (
	"errors"
	"fmt"
	"goboot/internal/model/shard"
	"time"

	"github.com/gophero/goal/logx"
	"gorm.io/gorm"
)

func init() {
	RegisterScheme()
}

func RegisterScheme() {
	shard.SchemeManager.RegisterType(historySchemePrefix+"1", shard.NewMetadata(historyTablePrefix+"1", &AccountHistory1{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"2", shard.NewMetadata(historyTablePrefix+"2", &AccountHistory2{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"3", shard.NewMetadata(historyTablePrefix+"3", &AccountHistory3{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"4", shard.NewMetadata(historyTablePrefix+"4", &AccountHistory4{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"5", shard.NewMetadata(historyTablePrefix+"5", &AccountHistory5{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"6", shard.NewMetadata(historyTablePrefix+"6", &AccountHistory6{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"7", shard.NewMetadata(historyTablePrefix+"7", &AccountHistory7{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"8", shard.NewMetadata(historyTablePrefix+"8", &AccountHistory8{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"9", shard.NewMetadata(historyTablePrefix+"9", &AccountHistory9{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"10", shard.NewMetadata(historyTablePrefix+"10", &AccountHistory10{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"11", shard.NewMetadata(historyTablePrefix+"11", &AccountHistory11{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"12", shard.NewMetadata(historyTablePrefix+"12", &AccountHistory12{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"13", shard.NewMetadata(historyTablePrefix+"13", &AccountHistory13{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"14", shard.NewMetadata(historyTablePrefix+"14", &AccountHistory14{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"15", shard.NewMetadata(historyTablePrefix+"15", &AccountHistory15{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"16", shard.NewMetadata(historyTablePrefix+"16", &AccountHistory16{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"17", shard.NewMetadata(historyTablePrefix+"17", &AccountHistory17{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"18", shard.NewMetadata(historyTablePrefix+"18", &AccountHistory18{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"19", shard.NewMetadata(historyTablePrefix+"19", &AccountHistory19{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"20", shard.NewMetadata(historyTablePrefix+"20", &AccountHistory20{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"21", shard.NewMetadata(historyTablePrefix+"21", &AccountHistory21{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"22", shard.NewMetadata(historyTablePrefix+"22", &AccountHistory22{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"23", shard.NewMetadata(historyTablePrefix+"23", &AccountHistory23{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"24", shard.NewMetadata(historyTablePrefix+"24", &AccountHistory24{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"25", shard.NewMetadata(historyTablePrefix+"25", &AccountHistory25{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"26", shard.NewMetadata(historyTablePrefix+"26", &AccountHistory26{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"27", shard.NewMetadata(historyTablePrefix+"27", &AccountHistory27{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"28", shard.NewMetadata(historyTablePrefix+"28", &AccountHistory28{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"29", shard.NewMetadata(historyTablePrefix+"29", &AccountHistory29{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"30", shard.NewMetadata(historyTablePrefix+"30", &AccountHistory30{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"31", shard.NewMetadata(historyTablePrefix+"31", &AccountHistory31{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"32", shard.NewMetadata(historyTablePrefix+"32", &AccountHistory32{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"33", shard.NewMetadata(historyTablePrefix+"33", &AccountHistory33{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"34", shard.NewMetadata(historyTablePrefix+"34", &AccountHistory34{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"35", shard.NewMetadata(historyTablePrefix+"35", &AccountHistory35{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"36", shard.NewMetadata(historyTablePrefix+"36", &AccountHistory36{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"37", shard.NewMetadata(historyTablePrefix+"37", &AccountHistory37{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"38", shard.NewMetadata(historyTablePrefix+"38", &AccountHistory38{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"39", shard.NewMetadata(historyTablePrefix+"39", &AccountHistory39{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"40", shard.NewMetadata(historyTablePrefix+"40", &AccountHistory40{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"41", shard.NewMetadata(historyTablePrefix+"41", &AccountHistory41{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"42", shard.NewMetadata(historyTablePrefix+"42", &AccountHistory42{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"43", shard.NewMetadata(historyTablePrefix+"43", &AccountHistory43{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"44", shard.NewMetadata(historyTablePrefix+"44", &AccountHistory44{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"45", shard.NewMetadata(historyTablePrefix+"45", &AccountHistory45{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"46", shard.NewMetadata(historyTablePrefix+"46", &AccountHistory46{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"47", shard.NewMetadata(historyTablePrefix+"47", &AccountHistory47{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"48", shard.NewMetadata(historyTablePrefix+"48", &AccountHistory48{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"49", shard.NewMetadata(historyTablePrefix+"49", &AccountHistory49{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"50", shard.NewMetadata(historyTablePrefix+"50", &AccountHistory50{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"51", shard.NewMetadata(historyTablePrefix+"51", &AccountHistory51{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"52", shard.NewMetadata(historyTablePrefix+"52", &AccountHistory52{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"53", shard.NewMetadata(historyTablePrefix+"53", &AccountHistory53{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"54", shard.NewMetadata(historyTablePrefix+"54", &AccountHistory54{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"55", shard.NewMetadata(historyTablePrefix+"55", &AccountHistory55{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"56", shard.NewMetadata(historyTablePrefix+"56", &AccountHistory56{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"57", shard.NewMetadata(historyTablePrefix+"57", &AccountHistory57{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"58", shard.NewMetadata(historyTablePrefix+"58", &AccountHistory58{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"59", shard.NewMetadata(historyTablePrefix+"59", &AccountHistory59{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"60", shard.NewMetadata(historyTablePrefix+"60", &AccountHistory60{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"61", shard.NewMetadata(historyTablePrefix+"61", &AccountHistory61{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"62", shard.NewMetadata(historyTablePrefix+"62", &AccountHistory62{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"63", shard.NewMetadata(historyTablePrefix+"63", &AccountHistory63{}))
	shard.SchemeManager.RegisterType(historySchemePrefix+"64", shard.NewMetadata(historyTablePrefix+"64", &AccountHistory64{}))
}

const (
	AccountHistoryShards = 64
	historySchemePrefix  = "account.history"
	historyTablePrefix   = "account_history"
)

func NewAccountHistoryModel(userId uint, logger *logx.Logger) IAccountHistory {
	n := userId % AccountHistoryShards
	mdl, err := shard.SchemeManager.New(fmt.Sprintf(historySchemePrefix+"%d", n+1))
	if err != nil {
		logger.Errorf("select model error: %v", err)
		return nil
	}
	return mdl.(IAccountHistory)
}

func SelHistoryTable(userId uint) string {
	n := userId % AccountHistoryShards
	return shard.SchemeManager.Table(fmt.Sprintf(historySchemePrefix+"%d", n+1))
}

func RawHistoryModel(userId uint) shard.IModel {
	n := userId % AccountHistoryShards
	return shard.SchemeManager.Model(fmt.Sprintf(historySchemePrefix+"%d", n+1))
}

func RawShardHistoryModel(n int) shard.IModel {
	return shard.SchemeManager.Model(fmt.Sprintf(historySchemePrefix+"%d", n))
}

type IAccountHistory interface {
	shard.IModel
	GetTradeCoin() float64
	SetTradeCoin(v float64)
	GetTradeCoinIncr() float64
	SetTradeCoinIncr(v float64)
	GetMainCoin() float64
	SetMainCoin(v float64)
	GetMainCoinIncr() float64
	SetMainCoinIncr(v float64)
	GetLockMainCoin() float64
	SetLockMainCoin(v float64)
	GetLockSecCoin()
	SetLockSecCoin(v float64)
	GetLockMainCoinIncr() float64
	SetLockMainCoinIncr(v float64)
	GetSecCoin() float64
	SetSecCoin(v float64)
	GetSecCoinIncr() float64
	SetSecCoinIncr(v float64)
	GetLockSecCoinIncr() float64
	SetLockSecCoinIncr(v float64)
	GetType() int8
	SetType(v int8)
	GetUserId() uint
	SetUserId(v uint)
	GetAccountId() uint
	SetAccountId(v uint)

	// 为了传递具体的model类型，必须增加IAccountHistory参数，否则无法获取，而获取的是父类BaseAccountHistory的实例

	Record(h IAccountHistory, account *Account, tradeCoinIncr float64, mainCoinIncr float64, lockMainCoinIncr float64, secCoinIncr float64, lockSecCoinIncr float64, typeVal int8, tx *gorm.DB) (IAccountHistory, error)
}

// 基本模型

type BaseAccountHistory struct {
	shard.BaseModel
	ID               uint
	CreatedAt        time.Time
	UpdatedAt        time.Time
	TradeCoin        float64
	TradeCoinIncr    float64
	MainCoin         float64
	MainCoinIncr     float64
	LockMainCoin     float64
	LockMainCoinIncr float64
	SecCoin          float64
	SecCoinIncr      float64
	LockSecCoin      float64
	LockSecCoinIncr  float64
	Type             int8
	UserId           uint
	AccountId        uint
}

func (h *BaseAccountHistory) Record(i IAccountHistory, account *Account, tradeCoinIncr float64, mainCoinIncr float64, lockMainCoinIncr float64, secCoinIncr float64, lockSecCoinIncr float64, typeVal int8, tx *gorm.DB) (IAccountHistory, error) {
	i.SetAccountId(account.ID)
	i.SetUserId(account.UserId)
	i.SetTradeCoin(account.TradeCoin)
	i.SetTradeCoinIncr(tradeCoinIncr)
	i.SetMainCoin(account.MainCoin)
	i.SetMainCoinIncr(mainCoinIncr)
	i.SetLockMainCoin(account.LockMainCoin)
	i.SetLockMainCoinIncr(lockMainCoinIncr)
	i.SetSecCoin(account.SecCoin)
	i.SetSecCoinIncr(secCoinIncr)
	i.SetLockSecCoin(account.LockSecCoin)
	i.SetLockSecCoinIncr(lockSecCoinIncr)
	i.SetType(typeVal)
	r := tx.Save(i)
	if r.Error != nil {
		return nil, r.Error
	}
	if r.RowsAffected < 1 {
		return nil, errors.New("insert error")
	}
	return i, nil
}

func (h *BaseAccountHistory) GetID() uint { return h.ID }

func (h *BaseAccountHistory) GetCreatedAt() *time.Time      { return &h.CreatedAt }
func (h *BaseAccountHistory) GetUpdatedAt() *time.Time      { return &h.UpdatedAt }
func (h *BaseAccountHistory) GetTradeCoin() float64         { return h.TradeCoin }
func (h *BaseAccountHistory) SetTradeCoin(v float64)        { h.TradeCoin = v }
func (h *BaseAccountHistory) GetTradeCoinIncr() float64     { return h.TradeCoinIncr }
func (h *BaseAccountHistory) SetTradeCoinIncr(v float64)    { h.TradeCoinIncr = v }
func (h *BaseAccountHistory) GetMainCoin() float64          { return h.MainCoin }
func (h *BaseAccountHistory) SetMainCoin(v float64)         { h.MainCoin = v }
func (h *BaseAccountHistory) GetMainCoinIncr() float64      { return h.MainCoinIncr }
func (h *BaseAccountHistory) SetMainCoinIncr(v float64)     { h.MainCoinIncr = v }
func (h *BaseAccountHistory) GetLockMainCoin() float64      { return h.LockMainCoin }
func (h *BaseAccountHistory) SetLockMainCoin(v float64)     { h.LockMainCoin = v }
func (h *BaseAccountHistory) GetLockMainCoinIncr() float64  { return h.LockMainCoinIncr }
func (h *BaseAccountHistory) SetLockMainCoinIncr(v float64) { h.LockMainCoinIncr = v }
func (h *BaseAccountHistory) GetSecCoin() float64           { return h.SecCoin }
func (h *BaseAccountHistory) SetSecCoin(v float64)          { h.SecCoin = v }
func (h *BaseAccountHistory) GetSecCoinIncr() float64       { return h.SecCoinIncr }
func (h *BaseAccountHistory) SetSecCoinIncr(v float64)      { h.SecCoinIncr = v }
func (h *BaseAccountHistory) GetLockSecCoin() float64       { return h.LockSecCoin }
func (h *BaseAccountHistory) SetLockSecCoin(v float64)      { h.LockSecCoin = v }
func (h *BaseAccountHistory) GetLockSecCoinIncr() float64   { return h.LockSecCoinIncr }
func (h *BaseAccountHistory) SetLockSecCoinIncr(v float64)  { h.LockSecCoinIncr = v }
func (h *BaseAccountHistory) GetUserId() uint               { return h.UserId }
func (h *BaseAccountHistory) SetUserId(v uint)              { h.UserId = v }
func (h *BaseAccountHistory) GetAccountId() uint            { return h.AccountId }
func (h *BaseAccountHistory) SetAccountId(v uint)           { h.AccountId = v }
func (h *BaseAccountHistory) GetType() int8                 { return h.Type }
func (h *BaseAccountHistory) SetType(v int8)                { h.Type = v }

// 分表模型

type (
	AccountHistory1  struct{ BaseAccountHistory }
	AccountHistory2  struct{ BaseAccountHistory }
	AccountHistory3  struct{ BaseAccountHistory }
	AccountHistory4  struct{ BaseAccountHistory }
	AccountHistory5  struct{ BaseAccountHistory }
	AccountHistory6  struct{ BaseAccountHistory }
	AccountHistory7  struct{ BaseAccountHistory }
	AccountHistory8  struct{ BaseAccountHistory }
	AccountHistory9  struct{ BaseAccountHistory }
	AccountHistory10 struct{ BaseAccountHistory }
	AccountHistory11 struct{ BaseAccountHistory }
	AccountHistory12 struct{ BaseAccountHistory }
	AccountHistory13 struct{ BaseAccountHistory }
	AccountHistory14 struct{ BaseAccountHistory }
	AccountHistory15 struct{ BaseAccountHistory }
	AccountHistory16 struct{ BaseAccountHistory }
	AccountHistory17 struct{ BaseAccountHistory }
	AccountHistory18 struct{ BaseAccountHistory }
	AccountHistory19 struct{ BaseAccountHistory }
	AccountHistory20 struct{ BaseAccountHistory }
	AccountHistory21 struct{ BaseAccountHistory }
	AccountHistory22 struct{ BaseAccountHistory }
	AccountHistory23 struct{ BaseAccountHistory }
	AccountHistory24 struct{ BaseAccountHistory }
	AccountHistory25 struct{ BaseAccountHistory }
	AccountHistory26 struct{ BaseAccountHistory }
	AccountHistory27 struct{ BaseAccountHistory }
	AccountHistory28 struct{ BaseAccountHistory }
	AccountHistory29 struct{ BaseAccountHistory }
	AccountHistory30 struct{ BaseAccountHistory }
	AccountHistory31 struct{ BaseAccountHistory }
	AccountHistory32 struct{ BaseAccountHistory }
	AccountHistory33 struct{ BaseAccountHistory }
	AccountHistory34 struct{ BaseAccountHistory }
	AccountHistory35 struct{ BaseAccountHistory }
	AccountHistory36 struct{ BaseAccountHistory }
	AccountHistory37 struct{ BaseAccountHistory }
	AccountHistory38 struct{ BaseAccountHistory }
	AccountHistory39 struct{ BaseAccountHistory }
	AccountHistory40 struct{ BaseAccountHistory }
	AccountHistory41 struct{ BaseAccountHistory }
	AccountHistory42 struct{ BaseAccountHistory }
	AccountHistory43 struct{ BaseAccountHistory }
	AccountHistory44 struct{ BaseAccountHistory }
	AccountHistory45 struct{ BaseAccountHistory }
	AccountHistory46 struct{ BaseAccountHistory }
	AccountHistory47 struct{ BaseAccountHistory }
	AccountHistory48 struct{ BaseAccountHistory }
	AccountHistory49 struct{ BaseAccountHistory }
	AccountHistory50 struct{ BaseAccountHistory }
	AccountHistory51 struct{ BaseAccountHistory }
	AccountHistory52 struct{ BaseAccountHistory }
	AccountHistory53 struct{ BaseAccountHistory }
	AccountHistory54 struct{ BaseAccountHistory }
	AccountHistory55 struct{ BaseAccountHistory }
	AccountHistory56 struct{ BaseAccountHistory }
	AccountHistory57 struct{ BaseAccountHistory }
	AccountHistory58 struct{ BaseAccountHistory }
	AccountHistory59 struct{ BaseAccountHistory }
	AccountHistory60 struct{ BaseAccountHistory }
	AccountHistory61 struct{ BaseAccountHistory }
	AccountHistory62 struct{ BaseAccountHistory }
	AccountHistory63 struct{ BaseAccountHistory }
	AccountHistory64 struct{ BaseAccountHistory }
)
