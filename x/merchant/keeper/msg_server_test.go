package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"natchain/x/merchant/keeper"
	"natchain/x/merchant/types"
)

func TestRegisterMerchant(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)

	resp, err := srv.RegisterMerchant(f.ctx, &types.MsgRegisterMerchant{
		Creator: "addr1", Name: "TestShop", CallbackUrl: "https://example.com", CommissionRate: 500,
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.ApiKey)

	m, err := f.keeper.Merchants.Get(f.ctx, "addr1")
	require.NoError(t, err)
	require.True(t, m.Active)
}

func TestDuplicateMerchant(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	srv.RegisterMerchant(f.ctx, &types.MsgRegisterMerchant{Creator: "a1", Name: "S1"})

	_, err := srv.RegisterMerchant(f.ctx, &types.MsgRegisterMerchant{Creator: "a1", Name: "S2"})
	require.Error(t, err)
}

func TestCreatePayment(t *testing.T) {
	f := initFixture(t)
	srv := keeper.NewMsgServerImpl(f.keeper)
	srv.RegisterMerchant(f.ctx, &types.MsgRegisterMerchant{Creator: "m1", Name: "Shop"})

	resp, err := srv.CreatePayment(f.ctx, &types.MsgCreatePayment{
		Creator: "u1", MerchantAddress: "m1", Amount: "100stake",
	})
	require.NoError(t, err)
	require.Contains(t, resp.PaymentId, "PAY-")
}
