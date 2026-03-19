package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"natchain/x/merchant/keeper"
	"natchain/x/merchant/types"
)

func TestQueryParams(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewQueryServerImpl(f.keeper)
	resp, err := srv.Params(f.ctx, &types.QueryParamsRequest{})
	require.NoError(t, err)
	require.Equal(t, int64(1000), resp.Params.MaxMerchants)
}

func TestQueryMerchant(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewQueryServerImpl(f.keeper)
	f.keeper.Merchants.Set(f.ctx, "a1", types.Merchant{Address: "a1", Name: "Shop"})

	resp, err := srv.Merchant(f.ctx, &types.QueryMerchantRequest{Address: "a1"})
	require.NoError(t, err)
	require.Equal(t, "Shop", resp.Merchant.Name)
}

func TestQueryMerchantNotFound(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewQueryServerImpl(f.keeper)
	_, err := srv.Merchant(f.ctx, &types.QueryMerchantRequest{Address: "nope"})
	require.Error(t, err)
}
