package model

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	TradeCoin             float64 `gorm:"type:decimal(40,20);default:0;comment:''" json:"tradeCoin" form:"tradeCoin"`
	MainCoin              float64 `gorm:"type:decimal(40,20);default:0;comment:''" json:"mainCoin" form:"mainCoin"`
	LockMainCoin          float64 `gorm:"type:decimal(40,20);default:0;comment:''" json:"lockMainCoin" form:"lockMainCoin"`
	SecCoin               float64 `gorm:"type:decimal(40,20);default:0;comment:''" json:"secCoin" form:"secCoin"`
	LockSecCoin           float64 `gorm:"type:decimal(40,20);default:0;comment:''" json:"lockSecCoin" form:"lockSecCoin"`
	UserId                uint    `gorm:"default:0;unique;comment:''" json:"userId" form:"userId"`
	State                 int8    `gorm:"type:smallint(2);default:0;comment:''" json:"state" form:"state"`
	Version               int64   `gorm:"type:bigint;default:1;comment:'version'"`
	TradePwd              string  `gorm:"type:varchar(128);comment:"`
	Address               string  `gorm:"type:varchar(128);comment:"`
	GiveFriendMainCoin    float64 `gorm:"type:decimal(40,20);default:0;comment:''"`
	GiveSecFriendMainCoin float64 `gorm:"type:decimal(40,20);default:0;comment:''"`
	FriendGiveMainCoin    float64 `gorm:"type:decimal(40,20);default:0;comment:''"`
	SecFriendGiveMainCoin float64 `gorm:"type:decimal(40,20);default:0;comment:''"`
	GiveFriendSecCoin     float64 `gorm:"type:decimal(40,20);default:0;comment:''"`
	GiveSecFriendSecCoin  float64 `gorm:"type:decimal(40,20);default:0;comment:''"`
	FriendGiveSecCoin     float64 `gorm:"type:decimal(40,20);default:0;comment:''"`
	SecFriendGiveSecCoin  float64 `gorm:"type:decimal(40,20);default:0;comment:''"`
	WithdrawQuota         float64 `gorm:"type:decimal(20,8);default:0;comment:''"`
}

const (
	Forbidden  = -1 // 封禁
	Suspicious = -2 // 可疑
)

const (
	AwardParentCoinType = SecCoinTokenType
)

const (
	TradeCoinTokenType = iota + 1
	SecCoinTokenType
	MainCoinTokenType
	LockMainCoinType
	LockSecCoinType
)

const (
	_               = iota
	ActTypeRecharge // 账户充值
	ActTypeWithdraw // 账户提现
	BindTwitter     // 绑定推特
	FollowTwitter   // 关注推特
	RetweetTwitter  // 转发推特
	PostTwitter     // 分享推特
	JoinDiscord     // 加入社区
	Invitefriend    // 邀请好友
	ActTypeRenewal

	ActTypeSwapConsume         // 兑换消耗
	ActTypeSwapReceive         // 兑换获得
	ActTypeAwardFromChild      // 奖励上级
	ActTypeAwardFromGrandchild // 奖励上上级
)

func GetCapitalFlowType() map[int8]string {
	resp := map[int8]string{}
	resp[ActTypeRecharge] = "Recharge"
	resp[ActTypeWithdraw] = "Withdraw"
	resp[BindTwitter] = "Link X"
	resp[FollowTwitter] = "Follow X"
	resp[RetweetTwitter] = "Retweet X"
	resp[PostTwitter] = "Post X"
	resp[JoinDiscord] = "Join Discord"
	resp[Invitefriend] = "Invite friends"
	resp[ActTypeRenewal] = "Renewal period"
	resp[ActTypeSwapConsume] = "Swap consume"
	resp[ActTypeSwapReceive] = "Swap receive"
	resp[ActTypeAwardFromChild] = "Friend award"
	resp[ActTypeAwardFromGrandchild] = "Friend award"
	return resp
}

const MinSwapSec = 200000.0
