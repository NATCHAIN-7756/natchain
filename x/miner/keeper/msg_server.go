package keeper

import (
	"bytes"
	"context"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"natchain/x/miner/types"
)

type msgServer struct {
	Keeper
}

func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

// UpdateParams 更新模块参数
func (k msgServer) UpdateParams(goCtx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	
	// 验证权限
	authorityBytes, err := k.addressCodec.StringToBytes(msg.Authority)
	if err != nil {
		return nil, sdkerrors.ErrUnauthorized.Wrap("invalid authority address")
	}
	
	if !bytes.Equal(authorityBytes, k.authority) {
		return nil, sdkerrors.ErrUnauthorized.Wrap("invalid authority")
	}
	
	// 保存参数
	k.SetParams(ctx, msg.Params)
	
	return &types.MsgUpdateParamsResponse{}, nil
}

// RegisterMiner 注册新矿机
func (k msgServer) RegisterMiner(goCtx context.Context, msg *types.MsgRegisterMiner) (*types.MsgRegisterMinerResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 检查矿机数量上限
	count := k.GetMinerCount(ctx)
	if count >= 1000 {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("max miners reached (1000)")
	}

	// 检查 CPU SN 是否已注册
	if k.HasMinerByCpuSn(ctx, msg.CpuSn) {
		return nil, sdkerrors.ErrInvalidRequest.Wrap("miner with this CPU SN already exists")
	}

	// 创建矿机
	miner := types.Miner{
		Index:         strconv.FormatUint(count+1, 10),
		Owner:         msg.Creator,
		CpuSn:         msg.CpuSn,
		RegisteredAt:  time.Now().Unix(),
		LastHeartbeat: time.Now().Unix(),
		Active:        true,
		TotalMined:    0,
	}

	// 保存
	k.SetMiner(ctx, miner)
	k.SetMinerCount(ctx, count+1)

	return &types.MsgRegisterMinerResponse{Index: miner.Index}, nil
}

// Heartbeat 矿机心跳
func (k msgServer) Heartbeat(goCtx context.Context, msg *types.MsgHeartbeat) (*types.MsgHeartbeatResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// 获取矿机
	miner, found := k.GetMiner(ctx, msg.Index)
	if !found {
		return nil, sdkerrors.ErrNotFound.Wrap("miner not found")
	}

	// 验证所有者
	if miner.Owner != msg.Creator {
		return nil, sdkerrors.ErrUnauthorized.Wrap("not miner owner")
	}

	// 更新心跳
	miner.LastHeartbeat = time.Now().Unix()
	miner.Active = true
	k.SetMiner(ctx, miner)

	return &types.MsgHeartbeatResponse{Success: true}, nil
}
