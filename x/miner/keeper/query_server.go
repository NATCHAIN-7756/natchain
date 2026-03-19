package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	
	"natchain/x/miner/types"
)

// Querier is used as QueryServer
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// NewQueryServerImpl returns an implementation of the QueryServer interface
func NewQueryServerImpl(k Keeper) types.QueryServer {
	return Querier{Keeper: k}
}

// Params queries the parameters of the module
func (q Querier) Params(goCtx context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	params, err := q.GetParams(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &types.QueryParamsResponse{Params: params}, nil
}

// Miner queries a single miner by index
func (q Querier) Miner(goCtx context.Context, req *types.QueryMinerRequest) (*types.QueryMinerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.Index == "" {
		return nil, status.Error(codes.InvalidArgument, "index cannot be empty")
	}
	
	ctx := sdk.UnwrapSDKContext(goCtx)
	miner, found := q.GetMiner(ctx, req.Index)
	if !found {
		return nil, status.Error(codes.NotFound, "miner not found")
	}
	
	return &types.QueryMinerResponse{Miner: miner}, nil
}

// MinerList queries all miners
func (q Querier) MinerList(goCtx context.Context, req *types.QueryMinerListRequest) (*types.QueryMinerListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	ctx := sdk.UnwrapSDKContext(goCtx)
	miners := q.GetAllMiner(ctx)
	
	return &types.QueryMinerListResponse{Miners: miners}, nil
}

// MinerCount queries the total miner count
func (q Querier) MinerCount(goCtx context.Context, req *types.QueryMinerCountRequest) (*types.QueryMinerCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	
	ctx := sdk.UnwrapSDKContext(goCtx)
	count := q.GetMinerCount(ctx)
	
	return &types.QueryMinerCountResponse{Count: count}, nil
}
