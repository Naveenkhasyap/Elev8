package gas

import (
	"context"
	"fmt"

	"github.com/NethermindEth/starknet.go/rpc"
)

func EstimateGas(ctx context.Context, client *rpc.Provider, signedTxn *rpc.InvokeTxnV1) ([]rpc.FeeEstimate, error) {
	res, err := client.EstimateFee(
		ctx,
		[]rpc.BroadcastTxn{signedTxn},
		[]rpc.SimulationFlag{rpc.SKIP_VALIDATE},
		rpc.WithBlockTag("latest"),
	)
	if err != nil {
		return []rpc.FeeEstimate{}, fmt.Errorf("unable to get gas estimate, err: %v", err)
	}
	return res, nil
}
