package keeper

import (
	"fmt"

	"cosmossdk.io/collections"
	"cosmossdk.io/core/address"
	corestore "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"cosmossdk.io/math"

	"natchain/x/wallet/types"
)

type Keeper struct {
	storeService corestore.KVStoreService
	cdc          codec.Codec
	addressCodec address.Codec
	authority    []byte

	Schema collections.Schema
	Params collections.Item[types.Params]

	// Wallet storage
	Wallets collections.Map[string, types.Wallet]
	
	// Transfer records storage
	TransferRecords collections.Map[uint64, types.TransferRecord]
	NextTransferId  collections.Item[uint64]

	// Partners storage (key: owner/address)
	PartnerStore collections.Map[string, types.Partner]

	bankKeeper types.BankKeeper
	authKeeper types.AuthKeeper
}

func NewKeeper(
	storeService corestore.KVStoreService,
	cdc codec.Codec,
	addressCodec address.Codec,
	authority []byte,
	bankKeeper types.BankKeeper,
	authKeeper types.AuthKeeper,
) Keeper {
	if _, err := addressCodec.BytesToString(authority); err != nil {
		panic(fmt.Sprintf("invalid authority address %s: %s", authority, err))
	}

	sb := collections.NewSchemaBuilder(storeService)

	k := Keeper{
		storeService: storeService,
		cdc:          cdc,
		addressCodec: addressCodec,
		authority:    authority,
		bankKeeper:   bankKeeper,
		authKeeper:   authKeeper,
		Params:       collections.NewItem(sb, collections.NewPrefix("params"), "params", codec.CollValue[types.Params](cdc)),
		Wallets:      collections.NewMap(sb, collections.NewPrefix(0), "wallets", collections.StringKey, codec.CollValue[types.Wallet](cdc)),
		TransferRecords: collections.NewMap(sb, collections.NewPrefix(1), "transfer_records", collections.Uint64Key, codec.CollValue[types.TransferRecord](cdc)),
		NextTransferId:  collections.NewItem(sb, collections.NewPrefix(2), "next_transfer_id", collections.Uint64Value),
		PartnerStore:    collections.NewMap(sb, collections.NewPrefix(3), "partners", collections.StringKey, codec.CollValue[types.Partner](cdc)),
	}

	schema, err := sb.Build()
	if err != nil {
		panic(err)
	}
	k.Schema = schema

	return k
}

func (k Keeper) GetAuthority() []byte {
	return k.authority
}

// AddTransferRecord adds a new transfer record and returns its ID
func (k Keeper) AddTransferRecord(sdkCtx sdk.Context, record types.TransferRecord) (uint64, error) {
	nextId, err := k.NextTransferId.Get(sdkCtx)
	if err != nil {
		nextId = 1
	}
	
	record.Id = nextId
	
	if err := k.TransferRecords.Set(sdkCtx, nextId, record); err != nil {
		return 0, err
	}
	
	if err := k.NextTransferId.Set(sdkCtx, nextId+1); err != nil {
		return 0, err
	}
	
	return nextId, nil
}

// GetTransferRecordsByAddress returns all transfer records for a given address
func (k Keeper) GetTransferRecordsByAddress(sdkCtx sdk.Context, address string) ([]types.TransferRecord, error) {
	var records []types.TransferRecord
	
	err := k.TransferRecords.Walk(sdkCtx, nil, func(key uint64, value types.TransferRecord) (bool, error) {
		if value.FromAddress == address || value.ToAddress == address {
			records = append(records, value)
		}
		return false, nil
	})
	
	return records, err
}

// GetPartnersByOwner returns all partners for a given owner
func (k Keeper) GetPartnersByOwner(sdkCtx sdk.Context, owner string) ([]types.Partner, error) {
	var partners []types.Partner
	
	err := k.PartnerStore.Walk(sdkCtx, nil, func(key string, value types.Partner) (bool, error) {
		if value.Owner == owner {
			partners = append(partners, value)
		}
		return false, nil
	})
	
	return partners, err
}

// CheckAutoTransfers 检查所有钱包的自动转出条件并执行
func (k Keeper) CheckAutoTransfers(ctx sdk.Context) int {
	transferred := 0
	
	err := k.Wallets.Walk(ctx, nil, func(key string, wallet types.Wallet) (bool, error) {
		// 检查是否启用自动转出
		if !wallet.AutoTransferEnabled || wallet.AutoTransferThreshold <= 0 {
			return false, nil
		}
		
		// 检查冷钱包地址是否设置
		if wallet.ColdAddress == "" {
			return false, nil
		}
		
		// 获取余额
		addr, err := sdk.AccAddressFromBech32(wallet.Address)
		if err != nil {
			return false, nil
		}
		
		balance := k.bankKeeper.SpendableCoins(ctx, addr)
		balanceInt := balance.AmountOf("unat") // 假设原生代币是 unat
		
		// 检查是否超过阈值
		if balanceInt.GTE(math.NewInt(wallet.AutoTransferThreshold)) {
			// 执行自动转出
			amount := wallet.AutoTransferAmount
			if amount == "" {
				// 默认转出超过阈值的部分
				excess := balanceInt.Sub(math.NewInt(wallet.AutoTransferThreshold))
				amount = excess.String()
			}
			
			// 创建转账
			coins, err := sdk.ParseCoinsNormalized(amount + "unat")
			if err != nil {
				return false, nil
			}
			
			coldAddr, err := sdk.AccAddressFromBech32(wallet.ColdAddress)
			if err != nil {
				return false, nil
			}
			
			if err := k.bankKeeper.SendCoins(ctx, addr, coldAddr, coins); err != nil {
				return false, nil
			}
			
			transferred++
		}
		
		return false, nil
	})
	
	if err != nil {
		ctx.Logger().Error("auto transfer check failed", "error", err)
	}
	
	return transferred
}
