package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/curve"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	starkrpc "github.com/gofiles/internal/clients/stark_rpc"
)

type IAccount interface {
	Name(ctx context.Context) (*felt.Felt, error)
	Nonce(ctx context.Context) (*felt.Felt, error)
	SignV1Txn(ctx context.Context, txn *rpc.InvokeTxnV1) error
	AccountAddress() *felt.Felt
	SignAndInvokeV1Txn(ctx context.Context, txn rpc.InvokeTxnV1) (*rpc.AddInvokeTransactionResponse, error)
}

type localAccount struct {
	accountAddress string
	client         *starkrpc.Provider
	publicKey      *felt.Felt
	account        *account.Account
}

func NewAccount(client *starkrpc.Provider, accountAddress string, privateKey string) (IAccount, error) {
	account_addr, err := utils.HexToFelt(accountAddress)
	if err != nil {
		return nil, fmt.Errorf("invalid address, err: %v", err)
	}

	ks := account.NewMemKeystore()
	pvKey, ok := new(big.Int).SetString(privateKey, 0)
	if !ok {
		return nil, fmt.Errorf("invalid private key: %v", privateKey)
	}

	x, _, _ := curve.Curve.PrivateToPoint(pvKey)
	pub := utils.BigIntToFelt(x)
	ks.Put(pub.String(), pvKey)

	acc, err1 := account.NewAccount(client, account_addr, pub.String(), ks, 2)
	if err1 != nil {
		return nil, fmt.Errorf("cannot create account, err: %v", err1)
	}

	return &localAccount{
		accountAddress: accountAddress,
		client:         client,
		publicKey:      pub,
		account:        acc,
	}, nil
}

func (la *localAccount) Name(ctx context.Context) (*felt.Felt, error) {
	params := rpc.FunctionCall{
		ContractAddress:    la.AccountAddress(),
		EntryPointSelector: utils.GetSelectorFromNameFelt("getName"),
	}

	res, err := la.client.Call(context.Background(), params, rpc.BlockID{Tag: "latest"})
	if err != nil {
		return nil, fmt.Errorf("account name, err: %v", err)
	}

	return res[0], nil
}

func (la *localAccount) Nonce(ctx context.Context) (*felt.Felt, error) {
	val, err := la.account.Nonce(ctx, rpc.WithBlockTag("pending"), la.account.AccountAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch nonce, err: %v", err)
	}
	return val, nil
}

func (la *localAccount) SignV1Txn(ctx context.Context, txn *rpc.InvokeTxnV1) error {
	return la.account.SignInvokeTransaction(ctx, txn)
}

func (la *localAccount) AccountAddress() *felt.Felt {
	return la.account.AccountAddress
}

func (la *localAccount) SignAndInvokeV1Txn(ctx context.Context, txn rpc.InvokeTxnV1) (*rpc.AddInvokeTransactionResponse, error) {
	err := la.SignV1Txn(ctx, &txn)
	if err != nil {
		return nil, fmt.Errorf("unable to sign the transaction, err: %v", err)
	}

	b, _ := json.Marshal(txn)
	fmt.Println(string(b))

	res, txnErr := la.account.AddInvokeTransaction(ctx, txn)
	if txnErr != nil {
		return nil, fmt.Errorf("add invoke txn error, err: %v", err)
	}

	return res, nil
}
