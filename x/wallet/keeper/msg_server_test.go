package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"natchain/x/wallet/keeper"
	"natchain/x/wallet/types"
)

func TestSetSafetyLimit(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	_, err := ms.SetSafetyLimit(f.ctx, &types.MsgSetSafetyLimit{Creator: "natchain1test", SafetyLimit: 5000})
	require.NoError(t, err)

	w, err := f.keeper.Wallets.Get(f.ctx, "natchain1test")
	require.NoError(t, err)
	require.Equal(t, int64(5000), w.SafetyLimit)
}

func TestSetColdAddress(t *testing.T) {
	f := initFixture(t)
	ms := keeper.NewMsgServerImpl(f.keeper)

	_, err := ms.SetColdAddress(f.ctx, &types.MsgSetColdAddress{Creator: "natchain1test", ColdAddress: "natchain1cold"})
	require.NoError(t, err)

	w, err := f.keeper.Wallets.Get(f.ctx, "natchain1test")
	require.NoError(t, err)
	require.Equal(t, "natchain1cold", w.ColdAddress)
}
