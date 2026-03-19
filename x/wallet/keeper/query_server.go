package keeper

import (
	"context"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"natchain/x/wallet/types"
)

type queryServer struct {
	Keeper
}

func NewQueryServerImpl(k Keeper) types.QueryServer {
	return &queryServer{Keeper: k}
}

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params, err := q.Keeper.Params.Get(sdkCtx)
	if err != nil {
		return nil, err
	}
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) Wallet(ctx context.Context, req *types.QueryWalletRequest) (*types.QueryWalletResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	wallet, err := q.Wallets.Get(sdkCtx, req.Address)
	if err != nil {
		return nil, types.ErrWalletNotFound
	}
	return &types.QueryWalletResponse{Wallet: wallet}, nil
}

func (q queryServer) WalletAll(ctx context.Context, req *types.QueryWalletAllRequest) (*types.QueryWalletAllResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var wallets []types.Wallet
	if err := q.Wallets.Walk(sdkCtx, nil, func(key string, value types.Wallet) (bool, error) {
		wallets = append(wallets, value)
		return false, nil
	}); err != nil {
		return nil, err
	}
	return &types.QueryWalletAllResponse{Wallets: wallets}, nil
}

func (q queryServer) Balance(ctx context.Context, req *types.QueryBalanceRequest) (*types.QueryBalanceResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	addr, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	balance := q.bankKeeper.SpendableCoins(sdkCtx, addr)
	
	// Check wallet for safety limit
	wallet, walletErr := q.Wallets.Get(sdkCtx, req.Address)
	threshold := ""
	excess := ""
	autoTransferTriggered := false
	
	if walletErr == nil && wallet.SafetyLimit > 0 {
		// Calculate threshold info
		balanceInt := balance.AmountOf("untnc") // Assuming native token is untnc
		thresholdInt := math.NewInt(wallet.SafetyLimit)
		threshold = thresholdInt.String()
		
		if balanceInt.GT(thresholdInt) {
			excess = balanceInt.Sub(thresholdInt).String()
			autoTransferTriggered = true
		}
	}
	
	return &types.QueryBalanceResponse{
		Address:              req.Address,
		Balance:              balance.String(),
		Threshold:            threshold,
		Excess:               excess,
		AutoTransferTriggered: autoTransferTriggered,
	}, nil
}

func (q queryServer) TransferRecords(ctx context.Context, req *types.QueryTransferRecordsRequest) (*types.QueryTransferRecordsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	records, err := q.GetTransferRecordsByAddress(sdkCtx, req.Address)
	if err != nil {
		return nil, err
	}
	
	return &types.QueryTransferRecordsResponse{Records: records}, nil
}

func (q queryServer) Partners(ctx context.Context, req *types.QueryPartnersRequest) (*types.QueryPartnersResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	partners, err := q.GetPartnersByOwner(sdkCtx, req.Owner)
	if err != nil {
		return nil, err
	}
	
	return &types.QueryPartnersResponse{Partners: partners}, nil
}

func (q queryServer) Partner(ctx context.Context, req *types.QueryPartnerRequest) (*types.QueryPartnerResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	
	key := req.Owner + "/" + req.Address
	partner, err := q.PartnerStore.Get(sdkCtx, key)
	if err != nil {
		return nil, types.ErrPartnerNotFound
	}
	
	return &types.QueryPartnerResponse{Partner: partner}, nil
}
