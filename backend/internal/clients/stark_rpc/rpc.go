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
	retryCount := 0
	for retryCount < 5 {
		status, err := i.GetTransactionStatus(ctx, txnHash)
		if err != nil {
			return TransactionReceipt{}, fmt.Errorf("unable to fetch status, txn: %s, err: %v", txnHash.String(), err)
		}

		exeStatus := status.ExecutionStatus
		finalStatus := status.FinalityStatus

		if finalStatus == "" {
			return TransactionReceipt{}, fmt.Errorf("invalid final status empty, txn: %s", txnHash.String())
		}

		if finalStatus == rpc.TxnStatus_Rejected {
			return TransactionReceipt{}, fmt.Errorf("status rejected, txn: %s", txnHash.String())
		} else if exeStatus == rpc.TxnExecutionStatusSUCCEEDED || finalStatus == rpc.TxnStatus_Accepted_On_L1 || finalStatus == rpc.TxnStatus_Accepted_On_L2 {
			return i.GetTransactionReceipt(ctx, txnHash)
		}

		retryCount += 1
		time.Sleep(5 * time.Second)
	}
	return TransactionReceipt{}, fmt.Errorf("retry count exceeded, txn: %s", txnHash.String())
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
