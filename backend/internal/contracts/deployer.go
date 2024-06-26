package contracts

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/gofiles/internal/accounts"
	starkrpc "github.com/gofiles/internal/clients/stark_rpc"
	"github.com/gofiles/internal/gas"
	"github.com/gofiles/internal/helpers"
)

type DeployerReq struct {
	Name      string
	Symbol    string
	Recipient string
}

type Deployer interface {
	Launch(ctx context.Context, req DeployerReq) (*felt.Felt, error)
}

type deployer struct {
	contractAddress *felt.Felt
	client          *starkrpc.Provider
	localAccount    accounts.IAccount
}

func NewDeployer(contractAddress string, client *starkrpc.Provider, localAccount accounts.IAccount) (Deployer, error) {
	ca, err := utils.HexToFelt(contractAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid deployer contract address: %s, err: %v", contractAddress, err)
	}
	return &deployer{
		contractAddress: ca,
		client:          client,
		localAccount:    localAccount,
	}, nil
}

func (d *deployer) Launch(ctx context.Context, req DeployerReq) (*felt.Felt, error) {
	receiver, rcerr := utils.HexToFelt(req.Recipient)
	if rcerr != nil {
		return nil, fmt.Errorf("invalid recipient address: %s, err: %v", req.Recipient, rcerr)
	}

	nf, err := helpers.GenerateCallDataforByteArray(req.Name)
	if err != nil {
		return nil, fmt.Errorf("unable to generate calldata for name bytearray: %s, err: %v", req.Name, err)
	}

	sf, err := helpers.GenerateCallDataforByteArray(req.Symbol)
	if err != nil {
		return nil, fmt.Errorf("unable to generate calldata for symbol bytearray: %s, err: %v", req.Symbol, err)
	}

	callData := append(append(nf, sf...), receiver)
	fnCall := rpc.FunctionCall{
		ContractAddress:    d.contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt("deployERC20"),
		Calldata:           callData,
	}

	txn, err := createInvokeTxnV1(ctx, d.localAccount, new(felt.Felt).SetUint64(0), []rpc.FunctionCall{fnCall})
	if err != nil {
		return nil, fmt.Errorf("unable to create txn for gas estimate, err: %v", err)
	}

	err = d.localAccount.SignV1Txn(ctx, &txn)
	if err != nil {
		return nil, fmt.Errorf("gas estimate sign, error: %v", err)
	}

	b, _ := json.Marshal(txn)
	slog.Info("txn objext for gas fee", "txn", string(b))

	ge, esterr := gas.EstimateGas(ctx, d.client, &txn)
	if esterr != nil {
		return nil, esterr
	}

	fee := ge[0].OverallFee
	extraFee := fee.Double(fee)
	txn1, err := createInvokeTxnV1(ctx, d.localAccount, extraFee, []rpc.FunctionCall{fnCall})
	if err != nil {
		return nil, err
	}

	c, _ := json.Marshal(txn1)
	slog.Info("txn for deploying smart contract", "txn", string(c))
	res, err := d.localAccount.SignAndInvokeV1Txn(ctx, txn1)
	if err != nil {
		return nil, err
	}

	return res.TransactionHash, nil
}
