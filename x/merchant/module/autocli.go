package merchant

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
	"natchain/x/merchant/types"
)

func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "Params", Use: "params", Short: "Query module params"},
				{RpcMethod: "Merchant", Use: "merchant [address]", Short: "Query merchant"},
				{RpcMethod: "MerchantAll", Use: "merchants", Short: "Query all merchants"},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: types.Msg_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{RpcMethod: "RegisterMerchant", Use: "register [name] [callback_url] [commission_rate]", Short: "Register as merchant"},
				{RpcMethod: "CreatePayment", Use: "create-payment [merchant] [amount]", Short: "Create payment"},
			},
		},
	}
}
