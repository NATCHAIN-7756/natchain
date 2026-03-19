package types

import (
	"errors"
)

// 错误定义
var (
	ErrMaxWalletsReached   = errors.New("max wallets reached")
	ErrWalletAlreadyExists = errors.New("wallet already exists")
	ErrWalletNotFound      = errors.New("wallet not found")
	ErrNotWalletOwner      = errors.New("not wallet owner")
	ErrPartnerNotFound     = errors.New("partner not found")
)
