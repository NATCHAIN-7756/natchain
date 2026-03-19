package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"natchain/x/wallet/types"
)

func (k Keeper) InitGenesis(ctx sdk.Context, genState types.GenesisState) error {
	// Set params
	if err := k.Params.Set(ctx, genState.Params); err != nil {
		return err
	}

	// Set wallets
	for _, wallet := range genState.Wallets {
		if err := k.Wallets.Set(ctx, wallet.Address, wallet); err != nil {
			return err
		}
	}

	// Set transfer records
	for _, record := range genState.TransferRecords {
		if err := k.TransferRecords.Set(ctx, record.Id, record); err != nil {
			return err
		}
	}

	// Set next transfer ID
	if genState.NextTransferId > 0 {
		if err := k.NextTransferId.Set(ctx, genState.NextTransferId); err != nil {
			return err
		}
	}

	return nil
}

func (k Keeper) ExportGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return nil, err
	}

	var wallets []types.Wallet
	if err := k.Wallets.Walk(ctx, nil, func(key string, value types.Wallet) (bool, error) {
		wallets = append(wallets, value)
		return false, nil
	}); err != nil {
		return nil, err
	}

	var transferRecords []types.TransferRecord
	if err := k.TransferRecords.Walk(ctx, nil, func(key uint64, value types.TransferRecord) (bool, error) {
		transferRecords = append(transferRecords, value)
		return false, nil
	}); err != nil {
		return nil, err
	}

	nextTransferId, err := k.NextTransferId.Get(ctx)
	if err != nil {
		nextTransferId = 0
	}

	return &types.GenesisState{
		Params:          params,
		Wallets:         wallets,
		TransferRecords: transferRecords,
		NextTransferId:  nextTransferId,
	}, nil
}
