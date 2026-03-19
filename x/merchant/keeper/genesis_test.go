package keeper_test

import (
	"testing"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	"natchain/x/merchant/types"
)

func TestMerchantInitExportGenesis(t *testing.T) {
	f := initFixture(t)
	genState := types.GenesisState{
		Params:    types.DefaultParams(),
		Merchants: []types.Merchant{{Address: "a1", Name: "Shop"}},
	}
	err := f.keeper.InitGenesis(f.ctx.(sdk.Context), genState)
	require.NoError(t, err)

	exported, err := f.keeper.ExportGenesis(f.ctx.(sdk.Context))
	require.NoError(t, err)
	require.Len(t, exported.Merchants, 1)
}
