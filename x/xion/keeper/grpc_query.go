package keeper

import (
	"context"
	"github.com/burnt-labs/xion/x/xion/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) PlatformPercentage(c context.Context, _ *types.QueryPlatformPercentageRequest) (*types.QueryPlatformPercentageResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	percentage := k.GetPlatformPercentage(ctx)
	resp := types.QueryPlatformPercentageResponse{PlatformPercentage: uint32(percentage.Uint64())}
	return &resp, nil
}
