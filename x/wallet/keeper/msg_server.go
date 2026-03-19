package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"natchain/x/wallet/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(k Keeper) types.MsgServer {
	return &msgServer{Keeper: k}
}

func (m msgServer) SetSafetyLimit(ctx context.Context, req *types.MsgSetSafetyLimit) (*types.MsgSetSafetyLimitResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Get or create wallet
	wallet, err := m.Wallets.Get(sdkCtx, req.Creator)
	if err != nil {
		// Wallet does not exist, create new one
		wallet = types.Wallet{
			Address: req.Creator,
		}
	}
	
	wallet.SafetyLimit = req.SafetyLimit
	if err := m.Wallets.Set(sdkCtx, req.Creator, wallet); err != nil {
		return nil, err
	}
	
	return &types.MsgSetSafetyLimitResponse{}, nil
}

func (m msgServer) SetColdAddress(ctx context.Context, req *types.MsgSetColdAddress) (*types.MsgSetColdAddressResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Get or create wallet
	wallet, err := m.Wallets.Get(sdkCtx, req.Creator)
	if err != nil {
		// Wallet does not exist, create new one
		wallet = types.Wallet{
			Address: req.Creator,
		}
	}
	
	wallet.ColdAddress = req.ColdAddress
	if err := m.Wallets.Set(sdkCtx, req.Creator, wallet); err != nil {
		return nil, err
	}
	
	return &types.MsgSetColdAddressResponse{}, nil
}

// AddPartner adds a new partner to the wallet
func (m msgServer) AddPartner(ctx context.Context, req *types.MsgAddPartner) (*types.MsgAddPartnerResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	partner := types.Partner{
		Owner:     req.Creator,
		Name:      req.Name,
		Address:   req.Address,
		Memo:      req.Memo,
		CreatedAt: sdkCtx.BlockTime().Unix(),
	}
	
	// Store partner with composite key: owner + address
	key := req.Creator + "/" + req.Address
	if err := m.PartnerStore.Set(sdkCtx, key, partner); err != nil {
		return nil, err
	}
	
	return &types.MsgAddPartnerResponse{}, nil
}

// RemovePartner removes a partner from the wallet
func (m msgServer) RemovePartner(ctx context.Context, req *types.MsgRemovePartner) (*types.MsgRemovePartnerResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	key := req.Creator + "/" + req.Address
	if err := m.PartnerStore.Remove(sdkCtx, key); err != nil {
		return nil, err
	}
	
	return &types.MsgRemovePartnerResponse{}, nil
}

// UpdatePartner updates an existing partner
func (m msgServer) UpdatePartner(ctx context.Context, req *types.MsgUpdatePartner) (*types.MsgUpdatePartnerResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	key := req.Creator + "/" + req.Address
	
	partner, err := m.PartnerStore.Get(sdkCtx, key)
	if err != nil {
		return nil, err
	}
	
	partner.Name = req.Name
	partner.Memo = req.Memo
	
	if err := m.PartnerStore.Set(sdkCtx, key, partner); err != nil {
		return nil, err
	}
	
	return &types.MsgUpdatePartnerResponse{}, nil
}

// SetAutoTransfer 设置自动转出配置
func (m msgServer) SetAutoTransfer(ctx context.Context, req *types.MsgSetAutoTransfer) (*types.MsgSetAutoTransferResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	// Get or create wallet
	wallet, err := m.Wallets.Get(sdkCtx, req.Creator)
	if err != nil {
		wallet = types.Wallet{
			Address: req.Creator,
		}
	}
	
	wallet.AutoTransferThreshold = req.Threshold
	wallet.AutoTransferAmount = req.Amount
	wallet.AutoTransferEnabled = req.Enabled
	
	if err := m.Wallets.Set(sdkCtx, req.Creator, wallet); err != nil {
		return nil, err
	}
	
	return &types.MsgSetAutoTransferResponse{}, nil
}
