package contracts

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/gofiles/accounts"
)

type DeployerReq struct {
	Name      string
	Symbol    string
	Recipient string
	Account   accounts.IAccount
}

type Deployer interface {
	Launch(ctx context.Context, req DeployerReq) error
}

type deployer struct {
	contractAddress *felt.Felt
	client          *rpc.Provider
}

func NewDeployer(contractAddress string, client *rpc.Provider) (Deployer, error) {
	ca, err := utils.HexToFelt(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid deployer contract address: %s, err: %v", contractAddress, err)
	}
	return &deployer{
		contractAddress: ca,
		client:          client,
	}, nil
}

func (d *deployer) Launch(ctx context.Context, req DeployerReq) error {
	receiver, rcerr := utils.HexToFelt(req.Recipient)
	if rcerr != nil {
		return fmt.Errorf("invalid recipient address: %s, err: %v", req.Recipient, rcerr)
	}

	nf, err := GenerateCallDataforByteArray(req.Name)
	if err != nil {
		return fmt.Errorf("unable to generate calldata for name bytearray: %s, err: %v", req.Name, err)
	}

	sf, err := GenerateCallDataforByteArray(req.Symbol)
	if err != nil {
		return fmt.Errorf("unable to generate calldata for symbol bytearray: %s, err: %v", req.Symbol, err)
	}

	callData := append(append(nf, sf...), receiver)
	fnCall := rpc.FunctionCall{
		ContractAddress:    d.contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt("deployERC20"),
		Calldata:           callData,
	}

	_, err = createInvokeTxnV1(ctx, req.Account, new(felt.Felt).SetUint64(0), []rpc.FunctionCall{fnCall})
	if err != nil {
		return fmt.Errorf("unable to create txn for gas estimate, err: %v", req.Symbol, err)
	}

	return nil
}
