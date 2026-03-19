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

	"natchain/x/wallet/keeper"
	module "natchain/x/wallet/module"
	"natchain/x/wallet/types"
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

func TestSetAndGetWallet(t *testing.T) {
	f := initFixture(t)
	wallet := types.Wallet{Address: "natchain1test", ColdAddress: "natchain1cold", SafetyLimit: 1000}
	err := f.keeper.Wallets.Set(f.ctx, wallet.Address, wallet)
	require.NoError(t, err)
	got, err := f.keeper.Wallets.Get(f.ctx, "natchain1test")
	require.NoError(t, err)
	require.Equal(t, int64(1000), got.SafetyLimit)
}

func TestWalletNotFound(t *testing.T) {
	f := initFixture(t)
	_, err := f.keeper.Wallets.Get(f.ctx, "nonexistent")
	require.Error(t, err)
}
