package keeper

import (
	"fmt"
	"context"
	"time"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"natchain/x/merchant/types"
)

type msgServer struct{ Keeper }
func NewMsgServerImpl(k Keeper) types.MsgServer { return &msgServer{k} }

func (m msgServer) RegisterMerchant(ctx context.Context, req *types.MsgRegisterMerchant) (*types.MsgRegisterMerchantResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	if _, err := m.Merchants.Get(sdkCtx, req.Creator); err == nil {
		return nil, types.ErrMerchantAlreadyExists
	}
	merchant := types.Merchant{
		Address:        req.Creator,
		Name:           req.Name,
		ApiKey:         m.generateApiKey(),
		CallbackUrl:    req.CallbackUrl,
		CreatedAt:      time.Now().Unix(),
		CommissionRate: req.CommissionRate,
		Active:         true,
	}
	if err := m.Merchants.Set(sdkCtx, req.Creator, merchant); err != nil { return nil, err }
	return &types.MsgRegisterMerchantResponse{ApiKey: merchant.ApiKey}, nil
}

func (m msgServer) UpdateMerchant(ctx context.Context, req *types.MsgUpdateMerchant) (*types.MsgUpdateMerchantResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	merchant, err := m.Merchants.Get(sdkCtx, req.Creator)
	if err != nil { return nil, types.ErrMerchantNotFound }
	merchant.Name = req.Name
	merchant.CallbackUrl = req.CallbackUrl
	merchant.CommissionRate = req.CommissionRate
	return &types.MsgUpdateMerchantResponse{}, m.Merchants.Set(sdkCtx, req.Creator, merchant)
}

func (m msgServer) RegenerateApiKey(ctx context.Context, req *types.MsgRegenerateApiKey) (*types.MsgRegenerateApiKeyResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	merchant, err := m.Merchants.Get(sdkCtx, req.Creator)
	if err != nil { return nil, types.ErrMerchantNotFound }
	merchant.ApiKey = m.generateApiKey()
	if err := m.Merchants.Set(sdkCtx, req.Creator, merchant); err != nil { return nil, err }
	return &types.MsgRegenerateApiKeyResponse{ApiKey: merchant.ApiKey}, nil
}

func (m msgServer) CreatePayment(ctx context.Context, req *types.MsgCreatePayment) (*types.MsgCreatePaymentResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	count, _ := m.PaymentCount.Get(sdkCtx)
	paymentId := fmt.Sprintf("PAY-%d", count)
	payment := types.PaymentRecord{
		Id:              paymentId,
		MerchantAddress: req.MerchantAddress,
		UserAddress:     req.Creator,
		Amount:          req.Amount,
		Status:          "pending",
		Timestamp:       time.Now().Unix(),
		CallbackData:    req.CallbackData,
	}
	if err := m.PaymentRecords.Set(sdkCtx, paymentId, payment); err != nil { return nil, err }
	m.PaymentCount.Set(sdkCtx, count+1)
	return &types.MsgCreatePaymentResponse{PaymentId: paymentId}, nil
}

func (m msgServer) ConfirmPayment(ctx context.Context, req *types.MsgConfirmPayment) (*types.MsgConfirmPaymentResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	payment, err := m.PaymentRecords.Get(sdkCtx, req.PaymentId)
	if err != nil { return nil, types.ErrPaymentNotFound }
	if payment.Status != "pending" { return nil, types.ErrInvalidPaymentStatus }
	
	// 更新支付状态
	payment.Status = req.Status
	if err := m.PaymentRecords.Set(sdkCtx, req.PaymentId, payment); err != nil {
		return nil, err
	}
	
	// 如果支付成功，触发 Webhook 回调
	if req.Status == "confirmed" {
		m.SendWebhookCallback(sdkCtx, payment)
	}
	
	return &types.MsgConfirmPaymentResponse{}, nil
}
