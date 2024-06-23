package contracts

import (
	"context"
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/gofiles/internal/accounts"
	starkrpc "github.com/gofiles/internal/clients/stark_rpc"
	"github.com/gofiles/internal/gas"
	"github.com/holiman/uint256"
)

type IERC20 interface {
	BalanceOf(ctx context.Context, walletAccount string) (*felt.Felt, error)
	Decimals(ctx context.Context) (*felt.Felt, error)
	Name(ctx context.Context) (*felt.Felt, error)
	Symbol(ctx context.Context) (*felt.Felt, error)
	Allowance(ctx context.Context, owner string, spender string) (*felt.Felt, error)

	Transfer(ctx context.Context, la accounts.IAccount, recipient string, amount *uint256.Int) (*felt.Felt, error)
}

type erc20 struct {
	client       *starkrpc.Provider
	tokenAccount *felt.Felt
}

func NewERC20(client *starkrpc.Provider, tokenAccount string) (IERC20, error) {
	token, err := utils.HexToFelt(tokenAccount)
	if err != nil {
		return &erc20{}, fmt.Errorf("invalid token address, err: %v", err)
	}
	return &erc20{
		client:       client,
		tokenAccount: token,
	}, nil
}

func (e *erc20) BalanceOf(ctx context.Context, walletAccount string) (*felt.Felt, error) {
	account, err := utils.HexToFelt(walletAccount)
	if err != nil {
		return nil, fmt.Errorf("invalid account address, err: %v", err)
	}

	params := rpc.FunctionCall{
		ContractAddress:    e.tokenAccount,
		EntryPointSelector: utils.GetSelectorFromNameFelt("balanceOf"),
		Calldata:           []*felt.Felt{account},
	}

	res, err1 := e.client.Call(context.Background(), params, rpc.BlockID{Tag: "latest"})
	if err1 != nil {
		return nil, fmt.Errorf("RPC error, err: %v", err)
	}

	return res[0], nil
}

func (e *erc20) Decimals(ctx context.Context) (*felt.Felt, error) {
	params := rpc.FunctionCall{
		ContractAddress:    e.tokenAccount,
		EntryPointSelector: utils.GetSelectorFromNameFelt("decimals"),
	}

	res, err := e.client.Call(context.Background(), params, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return nil, fmt.Errorf("RPC error, err: %v", err)
	}

	return res[0], nil
}

func (e *erc20) Name(ctx context.Context) (*felt.Felt, error) {
	params := rpc.FunctionCall{
		ContractAddress:    e.tokenAccount,
		EntryPointSelector: utils.GetSelectorFromNameFelt("name"),
	}

	res, err := e.client.Call(context.Background(), params, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return nil, fmt.Errorf("RPC error, err: %v", err)
	}

	return res[0], nil
}

func (e *erc20) Symbol(ctx context.Context) (*felt.Felt, error) {
	params := rpc.FunctionCall{
		ContractAddress:    e.tokenAccount,
		EntryPointSelector: utils.GetSelectorFromNameFelt("symbol"),
	}

	res, err := e.client.Call(context.Background(), params, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return nil, fmt.Errorf("RPC error, err: %v", err)
	}

	return res[0], nil
}

func (e *erc20) Allowance(ctx context.Context, owner string, spender string) (*felt.Felt, error) {
	o, err := utils.HexToFelt(owner)
	if err != nil {
		return nil, fmt.Errorf("invalid owner contract: %s, err: %v", owner, err)
	}

	s, err1 := utils.HexToFelt(spender)
	if err1 != nil {
		return nil, fmt.Errorf("invalid spender contract: %s, err: %v", spender, err1)
	}

	params := rpc.FunctionCall{
		ContractAddress:    e.tokenAccount,
		EntryPointSelector: utils.GetSelectorFromNameFelt("allowance"),
		Calldata:           []*felt.Felt{o, s},
	}

	res, err2 := e.client.Call(context.Background(), params, rpc.BlockID{Tag: "latest"})
	if err2 != nil {
		return nil, fmt.Errorf("RPC error, err: %v", err2)
	}

	return res[0], nil
}

func (e *erc20) Transfer(ctx context.Context, la accounts.IAccount, recipient string, amount *uint256.Int) (*felt.Felt, error) {
	r, err := utils.HexToFelt(recipient)
	if err != nil {
		return nil, fmt.Errorf("invalid recepient address: %s, err: %v", recipient, err)
	}

	amnt, err1 := utils.HexToFelt(amount.Hex())
	if err1 != nil {
		return nil, fmt.Errorf("invalid amount: %s, err: %v", amount.String(), err1)
	}

	fnCall := rpc.FunctionCall{
		ContractAddress:    e.tokenAccount,
		EntryPointSelector: utils.GetSelectorFromNameFelt("transfer"),
		Calldata:           []*felt.Felt{r, amnt, new(felt.Felt).SetUint64(0)},
	}

	txn, err := createInvokeTxnV1(ctx, la, new(felt.Felt).SetUint64(0), []rpc.FunctionCall{fnCall})
	if err != nil {
		return nil, err
	}

	fee, err := e.gasEstimate(ctx, la, txn)
	if err != nil {
		return nil, err
	}

	extraFee := fee.Add(fee, new(felt.Felt).SetUint64(10000000000))
	txn1, err := createInvokeTxnV1(ctx, la, extraFee, []rpc.FunctionCall{fnCall})
	if err != nil {
		return nil, err
	}

	res, err := la.SignAndInvokeV1Txn(ctx, txn1)
	if err != nil {
		return nil, err
	}
	return res.TransactionHash, nil
}

func createInvokeTxnV1(ctx context.Context, la accounts.IAccount, maxfee *felt.Felt, fnCalls []rpc.FunctionCall) (rpc.InvokeTxnV1, error) {
	callData := account.FmtCallDataCairo2(fnCalls)
	nonce, err := la.Nonce(ctx)
	if err != nil {
		return rpc.InvokeTxnV1{}, fmt.Errorf("no nonce, err: %v", err)
	}

	fee := new(felt.Felt).SetUint64(0)
	if maxfee != nil && !maxfee.IsZero() {
		fee = maxfee
	}

	return rpc.InvokeTxnV1{
		MaxFee:        fee,
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: la.AccountAddress(),
		Calldata:      callData,
	}, nil
}

func (e *erc20) gasEstimate(ctx context.Context, la accounts.IAccount, txn rpc.InvokeTxnV1) (*felt.Felt, error) {
	err := la.SignV1Txn(ctx, &txn)
	if err != nil {
		return nil, fmt.Errorf("gas estimate sign, error: %v", err)
	}

	estimate, esterr := gas.EstimateGas(ctx, e.client, &txn)
	if esterr != nil {
		return nil, fmt.Errorf("gas estimate, error: %v", esterr)
	}
	return estimate[0].OverallFee, nil
}
