package keeper_test

import (
	"testing"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	"natchain/x/wallet/types"
)

func TestInitExportGenesis(t *testing.T) {
	f := initFixture(t)
	genState := types.GenesisState{
		Params:  types.DefaultParams(),
		Wallets: []types.Wallet{{Address: "addr1", SafetyLimit: 100}},
	}
	err := f.keeper.InitGenesis(f.ctx.(sdk.Context), genState)
	require.NoError(t, err)

	exported, err := f.keeper.ExportGenesis(f.ctx.(sdk.Context))
	require.NoError(t, err)
	require.Len(t, exported.Wallets, 1)
}
