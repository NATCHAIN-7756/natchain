package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"natchain/x/wallet/keeper"
	"natchain/x/wallet/types"
)

func TestQueryParams(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewQueryServerImpl(f.keeper)
	resp, err := srv.Params(f.ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestQueryWalletAll(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewQueryServerImpl(f.keeper)
	f.keeper.Wallets.Set(f.ctx, "a1", types.Wallet{Address: "a1"})
	f.keeper.Wallets.Set(f.ctx, "a2", types.Wallet{Address: "a2"})

	resp, err := srv.WalletAll(f.ctx, &types.QueryWalletAllRequest{})
	require.NoError(t, err)
	require.Len(t, resp.Wallets, 2)
}
