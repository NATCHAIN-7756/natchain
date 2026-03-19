package keeper

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"natchain/x/miner/types"
)

// SetMiner 保存矿机
func (k Keeper) SetMiner(ctx sdk.Context, miner types.Miner) {
	store := k.storeService.OpenKVStore(ctx)
	bz, _ := json.Marshal(miner)
	store.Set([]byte(types.MinerKey+miner.Index), bz)
}

// GetMiner 获取矿机
func (k Keeper) GetMiner(ctx sdk.Context, index string) (types.Miner, bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get([]byte(types.MinerKey + index))
	if err != nil || bz == nil {
		return types.Miner{}, false
	}
	var miner types.Miner
	json.Unmarshal(bz, &miner)
	return miner, true
}

// GetAllMiner 获取所有矿机
func (k Keeper) GetAllMiner(ctx sdk.Context) []types.Miner {
	var miners []types.Miner
	store := k.storeService.OpenKVStore(ctx)
	iter, _ := store.Iterator([]byte(types.MinerKey), []byte(types.MinerKey+"\xff"))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var miner types.Miner
		json.Unmarshal(iter.Value(), &miner)
		miners = append(miners, miner)
	}
	return miners
}

// GetMinerCount 获取矿机数量
func (k Keeper) GetMinerCount(ctx sdk.Context) uint64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, _ := store.Get([]byte(types.MinerCountKey))
	if bz == nil {
		return 0
	}
	count, _ := strconv.ParseUint(string(bz), 10, 64)
	return count
}

// SetMinerCount 设置矿机数量
func (k Keeper) SetMinerCount(ctx sdk.Context, count uint64) {
	store := k.storeService.OpenKVStore(ctx)
	store.Set([]byte(types.MinerCountKey), []byte(strconv.FormatUint(count, 10)))
}

// HasMinerByCpuSn 检查 CPU SN 是否已注册
func (k Keeper) HasMinerByCpuSn(ctx sdk.Context, cpuSn string) bool {
	miners := k.GetAllMiner(ctx)
	for _, m := range miners {
		if m.CpuSn == cpuSn {
			return true
		}
	}
	return false
}

// CheckOfflineMiners 检查离线矿机并自动解绑
// 超过 180 天无心跳的矿机自动标记为非活跃
func (k Keeper) CheckOfflineMiners(ctx sdk.Context) int {
	miners := k.GetAllMiner(ctx)
	currentTime := ctx.BlockTime().Unix()
	offlineThreshold := int64(180 * 24 * 60 * 60) // 180 天（秒）
	unboundCount := 0

	for _, miner := range miners {
		if miner.Active && (currentTime-miner.LastHeartbeat) > offlineThreshold {
			miner.Active = false
			k.SetMiner(ctx, miner)
			unboundCount++
		}
	}

	return unboundCount
}

// GetActiveMiners 获取所有活跃矿机
func (k Keeper) GetActiveMiners(ctx sdk.Context) []types.Miner {
	var activeMiners []types.Miner
	for _, m := range k.GetAllMiner(ctx) {
		if m.Active {
			activeMiners = append(activeMiners, m)
		}
	}
	return activeMiners
}
