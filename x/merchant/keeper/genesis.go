package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"natchain/x/merchant/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) error {
	if err := k.Params.Set(ctx, genState.Params); err != nil { return err }
	for _, m := range genState.Merchants {
		if err := k.Merchants.Set(ctx, m.Address, m); err != nil { return err }
	}
	for _, p := range genState.PaymentRecords {
		if err := k.PaymentRecords.Set(ctx, p.Id, p); err != nil { return err }
	}
	return k.PaymentCount.Set(ctx, genState.NextPaymentId)
}

func (k Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil { return nil, err }
	var merchants []types.Merchant
	k.Merchants.Walk(ctx, nil, func(key string, value types.Merchant) (bool, error) {
		merchants = append(merchants, value)
		return false, nil
	})
	var payments []types.PaymentRecord
	k.PaymentRecords.Walk(ctx, nil, func(key string, value types.PaymentRecord) (bool, error) {
		payments = append(payments, value)
		return false, nil
	})
	count, _ := k.PaymentCount.Get(ctx)
	return &types.GenesisState{Params: params, Merchants: merchants, PaymentRecords: payments, NextPaymentId: count}, nil
}
