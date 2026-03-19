package types

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type AuthKeeper interface {
	GetAccount(context.Context, sdk.AccAddress) sdk.AccountI
}

type BankKeeper interface {
	SpendableCoins(context.Context, sdk.AccAddress) sdk.Coins
	SendCoins(ctx context.Context, from, to sdk.AccAddress, amt sdk.Coins) error
}
