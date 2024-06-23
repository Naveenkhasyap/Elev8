package contracts

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/gofiles/internal/accounts"
	starkrpc "github.com/gofiles/internal/clients/stark_rpc"
	"github.com/holiman/uint256"
)

type Pool interface {
	Quote(ctx context.Context, amountIn *uint256.Int, side uint8) (string, error)
	Swap(ctx context.Context, amountIn *uint256.Int, side uint8) error
}

type pool struct {
	contractAddress *felt.Felt
	client          *starkrpc.Provider
	localAccount    accounts.IAccount
}

func NewPool(contractAddress string, client *starkrpc.Provider, localAccount accounts.IAccount) (*pool, error) {
	ca, err := utils.HexToFelt(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid pool contract address: %s, err: %v", contractAddress, err)
	}
	return &pool{
		contractAddress: ca,
		client:          client,
		localAccount:    localAccount,
	}, nil
}

func (p *pool) Quote(ctx context.Context, amountIn *uint256.Int, side uint8) (string, error) {
	val, _ := utils.HexToFelt("0x214e8348c4f0000")
	callData := []*felt.Felt{
		val,
		new(felt.Felt).SetUint64(uint64(side)),
		new(felt.Felt).SetUint64(uint64(side)),
	}
	fnCall := rpc.FunctionCall{
		ContractAddress:    p.contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt("quote"),
		Calldata:           callData,
	}

	res, err2 := p.client.Call(context.Background(), fnCall, rpc.BlockID{Tag: "latest"})
	if err2 != nil {
		return "", fmt.Errorf("RPC error, err: %v", err2)
	}

	return res[0].String(), nil
}

func (p *pool) Swap(ctx context.Context, amountIn *uint256.Int, side uint8) error {
	return nil
}
