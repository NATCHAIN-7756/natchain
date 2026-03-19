package keeper_test

import (
	"context"
	"testing"

	"cosmossdk.io/core/address"
	storetypes "cosmossdk.io/store/types"
	addresscodec "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"natchain/x/merchant/keeper"
	module "natchain/x/merchant/module"
	"natchain/x/merchant/types"
)

type fixture struct {
	ctx          context.Context
	keeper       keeper.Keeper
	addressCodec address.Codec
}

func initFixture(t *testing.T) *fixture {
	t.Helper()
	encCfg := moduletestutil.MakeTestEncodingConfig(module.AppModule{})
	addressCodec := addresscodec.NewBech32Codec(sdk.GetConfig().GetBech32AccountAddrPrefix())
	storeKey := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(storeKey)
	ctx := testutil.DefaultContextWithDB(t, storeKey, storetypes.NewTransientStoreKey("transient_test")).Ctx
	authority := authtypes.NewModuleAddress(types.GovModuleName)
	k := keeper.NewKeeper(storeService, encCfg.Codec, addressCodec, authority, nil, nil)
	k.Params.Set(ctx, types.DefaultParams())
	return &fixture{ctx: ctx, keeper: k, addressCodec: addressCodec}
}

func TestSetAndGetMerchant(t *testing.T) {
	f := initFixture(t)
	m := types.Merchant{Address: "addr1", Name: "TestShop", Active: true}
	err := f.keeper.Merchants.Set(f.ctx, m.Address, m)
	require.NoError(t, err)
	got, err := f.keeper.Merchants.Get(f.ctx, "addr1")
	require.NoError(t, err)
	require.Equal(t, "TestShop", got.Name)
}

func TestMerchantNotFound(t *testing.T) {
	f := initFixture(t)
	_, err := f.keeper.Merchants.Get(f.ctx, "nope")
	require.Error(t, err)
}
