package wallet

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	"natchain/x/wallet/types"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: types.Query_serviceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "Wallet",
					Use:       "wallet [address]",
					Short:     "Query wallet by address",
				},
				{
					RpcMethod: "WalletAll",
					Use:       "wallets",
					Short:     "Query all wallets",
				},
				{
					RpcMethod: "Balance",
					Use:       "balance [address]",
					Short:     "Query wallet balance and threshold info",
				},
				{
					RpcMethod: "TransferRecords",
					Use:       "transfer-records [address]",
					Short:     "Query transfer records for an address",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: types.Msg_serviceDesc.ServiceName,
			EnhanceCustomCommand: true,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "SetSafetyLimit",
					Use:       "set-safety-limit [safety-limit]",
					Short:     "Set wallet safety limit",
				},
				{
					RpcMethod: "SetColdAddress",
					Use:       "set-cold-address [cold-address]",
					Short:     "Set cold wallet address",
				},
			},
		},
	}
}
