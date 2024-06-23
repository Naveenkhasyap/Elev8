package starkrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/avast/retry-go"
)

type Provider struct {
	url string
	*rpc.Provider
}

func NewProvider(url string, client *rpc.Provider) *Provider {
	return &Provider{
		url:      url,
		Provider: client,
	}
}

func (i *Provider) WaitForTransaction(ctx context.Context, txnHash *felt.Felt) (TransactionReceipt, error) {
	var receipt TransactionReceipt
	err := retry.Do(func() error {
		status, err := i.GetTransactionStatus(ctx, txnHash)
		if err != nil {
			return fmt.Errorf("unable to fetch status, txn: %s, err: %v", txnHash.String(), err)
		}

		exeStatus := status.ExecutionStatus
		finalStatus := status.FinalityStatus

		if finalStatus == "" {
			return fmt.Errorf("invalid final status empty, txn: %s", txnHash.String())
		}

		if finalStatus == rpc.TxnStatus_Rejected {
			return fmt.Errorf("status rejected, txn: %s", txnHash.String())
		} else if exeStatus == rpc.TxnExecutionStatusSUCCEEDED || finalStatus == rpc.TxnStatus_Accepted_On_L1 || finalStatus == rpc.TxnStatus_Accepted_On_L2 {
			rcpt, err := i.waitForReceipt(ctx, txnHash)
			if err != nil {
				return fmt.Errorf("unable to get receipt, txn: %s", txnHash.String())
			}
			receipt = rcpt
		}
		return nil
	}, retry.Delay(1*time.Second), retry.Attempts(20))

	return receipt, err
}

func (i *Provider) waitForReceipt(ctx context.Context, txnHash *felt.Felt) (TransactionReceipt, error) {
	var receipt TransactionReceipt
	err := retry.Do(func() error {
		rcpt, err := i.GetTransactionReceipt(ctx, txnHash)
		if err != nil {
			return fmt.Errorf("unable to get receipt, txn: %s", txnHash.String())
		}
		receipt = rcpt
		return nil
	}, retry.Delay(1*time.Second), retry.Attempts(10))

	return receipt, err
}

func (i *Provider) GetTransactionReceipt(ctx context.Context, hash *felt.Felt) (TransactionReceipt, error) {
	client := &http.Client{}
	request := map[string]any{
		"jsonrpc": "2.0",
		"method":  "starknet_getTransactionReceipt",
		"params": map[string]any{
			"transaction_hash": hash.String(),
		},
		"id": 1,
	}
	b, _ := json.Marshal(request)
	payload := bytes.NewReader(b)
	req, err := http.NewRequest("POST", i.url, payload)
	req = req.WithContext(ctx)

	if err != nil {
		return TransactionReceipt{}, fmt.Errorf("unable to create http req, err: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return TransactionReceipt{}, fmt.Errorf("unable to call, err: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return TransactionReceipt{}, fmt.Errorf("unable to read response, err: %v", err)
	}

	receipt := TransactionReceipt{}
	err = json.Unmarshal(body, &receipt)
	if err != nil {
		return TransactionReceipt{}, fmt.Errorf("invalid response, err: %v", err)
	}
	return receipt, nil
}
