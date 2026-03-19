package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"natchain/x/merchant/types"
)

type queryServer struct{ Keeper }
func NewQueryServerImpl(k Keeper) types.QueryServer { return &queryServer{k} }

func (q queryServer) Params(ctx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params, err := q.Keeper.Params.Get(sdkCtx)
	if err != nil { return nil, err }
	return &types.QueryParamsResponse{Params: params}, nil
}

func (q queryServer) Merchant(ctx context.Context, req *types.QueryMerchantRequest) (*types.QueryMerchantResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	m, err := q.Merchants.Get(sdkCtx, req.Address)
	if err != nil { return nil, types.ErrMerchantNotFound }
	return &types.QueryMerchantResponse{Merchant: m}, nil
}

func (q queryServer) MerchantAll(ctx context.Context, req *types.QueryMerchantAllRequest) (*types.QueryMerchantAllResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var merchants []types.Merchant
	q.Merchants.Walk(sdkCtx, nil, func(key string, value types.Merchant) (bool, error) {
		merchants = append(merchants, value)
		return false, nil
	})
	return &types.QueryMerchantAllResponse{Merchants: merchants}, nil
}

func (q queryServer) PaymentRecord(ctx context.Context, req *types.QueryPaymentRecordRequest) (*types.QueryPaymentRecordResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	p, err := q.PaymentRecords.Get(sdkCtx, req.Id)
	if err != nil { return nil, types.ErrPaymentNotFound }
	return &types.QueryPaymentRecordResponse{PaymentRecord: p}, nil
}

func (q queryServer) PaymentRecordsByMerchant(ctx context.Context, req *types.QueryPaymentRecordsByMerchantRequest) (*types.QueryPaymentRecordsByMerchantResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	var payments []types.PaymentRecord
	q.PaymentRecords.Walk(sdkCtx, nil, func(key string, value types.PaymentRecord) (bool, error) {
		if value.MerchantAddress == req.MerchantAddress {
			payments = append(payments, value)
		}
		return false, nil
	})
	return &types.QueryPaymentRecordsByMerchantResponse{PaymentRecords: payments}, nil
}
