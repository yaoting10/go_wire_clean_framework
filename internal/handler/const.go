package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gophero/goal/errorx"
)

const (
	illegalReq           = "illegal_req"
	networkBusy          = "network_busy"
	balanceInsufficient  = "balance_insufficient"
	balanceInsufficientF = "balance_insufficient_f"
	chainHasPaddingTx    = "chain_has_padding_tx"
	accountIllegal       = "account_illegal"
	walletNotExisted     = "wallet_not_existed"
	invalidInviteCode    = "invalid_invite_code"
	accountNotExist      = "account_not_exist"
	invalidEmail         = "invalid_email"
	emailRegistered      = "email_registered"
	invalidVerifyCode    = "invalid_verify_code"
	passwordError        = "password_error"
	walletExisted        = "wallet_existed"
	noWallet             = "no_wallet"
	mnemonicError        = "mnemonic_error"
	walletPasswordError  = "wallet_password_error"
	privateKeyError      = "private_key_error"
	deviceError          = "device_error"
)

var NotOpenErr = fmt.Errorf("It's comming soon")

func notOpen(ctx *gin.Context) {
	errorx.Throw(NotOpenErr)
}
