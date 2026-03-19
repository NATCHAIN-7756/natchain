package keeper

import (
	"fmt"
	"crypto/rand"
	"encoding/hex"
	
	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	
	sdk "github.com/cosmos/cosmos-sdk/types"
	"natchain/x/merchant/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	authority    []byte
	
	Schema collections.Schema
	Params collections.Item[types.Params]
	
	Merchants     collections.Map[string, types.Merchant]
	PaymentRecords collections.Map[string, types.PaymentRecord]
	PaymentCount  collections.Item[uint64]
	
	bankKeeper types.BankKeeper
	authKeeper types.AuthKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper types.BankKeeper,
	authKeeper types.AuthKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}
	
	sb := collections.NewSchemaBuilder(storeService)
	
	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		bankKeeper:   bankKeeper,
		authKeeper:   authKeeper,
		Params:       collections.NewItem(sb, collections.NewPrefix("params"), "params", codec.CollValue[types.Params](cdc)),
		Merchants:    collections.NewMap(sb, collections.NewPrefix(0), "merchants", collections.StringKey, codec.CollValue[types.Merchant](cdc)),
		PaymentRecords: collections.NewMap(sb, collections.NewPrefix(1), "payment_records", collections.StringKey, codec.CollValue[types.PaymentRecord](cdc)),
		PaymentCount: collections.NewItem(sb, collections.NewPrefix(2), "payment_count", collections.Uint64Value),
	}
	
	schema, err := sb.Build()
	if err != nil { panic(err) }
	k.Schema = schema
	
	return k
}

func (k Keeper) GetAuthority() []byte { return k.authority }

func (k Keeper) generateApiKey() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

// SendWebhookCallback 发送 Webhook 回调通知
// 通过事件机制通知链外服务发送 HTTP 回调
func (k Keeper) SendWebhookCallback(ctx sdk.Context, payment types.PaymentRecord) error {
	merchant, err := k.Merchants.Get(ctx, payment.MerchantAddress)
	if err != nil || merchant.CallbackUrl == "" {
		return nil
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"webhook_callback",
			sdk.NewAttribute("payment_id", payment.Id),
			sdk.NewAttribute("callback_url", merchant.CallbackUrl),
			sdk.NewAttribute("merchant_address", payment.MerchantAddress),
			sdk.NewAttribute("amount", payment.Amount),
			sdk.NewAttribute("status", payment.Status),
		),
	)
	return nil
}
