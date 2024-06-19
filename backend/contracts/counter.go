package contracts

import (
	"fmt"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
)

const CONTRACT_ADDRESS = "0x017a103c65198abd3ff9777d2b1127fd4a2ce1ff90b6eb0b66d521d22d906d17"
const INCREMENT_COUNTER = "increase_counter"
const DECEREMENT_COUNTER = "decrement"
const READ_COUNTER = "get_counter"

func CreateIncrementTransaction(nonce *felt.Felt, senderAddress *felt.Felt) (rpc.InvokeTxnV1, error) {
	cn, err := utils.HexToFelt(CONTRACT_ADDRESS)
	if err != nil {
		return rpc.InvokeTxnV1{}, fmt.Errorf("invalid contract address, err: %v", err)
	}

	fnCall := rpc.FunctionCall{
		ContractAddress:    cn,
		EntryPointSelector: utils.GetSelectorFromNameFelt(INCREMENT_COUNTER),
	}

	callData := account.FmtCallDataCairo2([]rpc.FunctionCall{fnCall})
	invokeTxn := rpc.InvokeTxnV1{
		MaxFee:        new(felt.Felt).SetUint64(0),
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: senderAddress,
		Calldata:      callData,
	}

	return invokeTxn, nil
}

func CreateDecrementTransaction(nonce *felt.Felt, senderAddress *felt.Felt) (rpc.InvokeTxnV1, error) {
	cn, err := utils.HexToFelt(CONTRACT_ADDRESS)
	if err != nil {
		return rpc.InvokeTxnV1{}, fmt.Errorf("invalid contract address, err: %v", err)
	}

	fnCall := rpc.FunctionCall{
		ContractAddress:    cn,
		EntryPointSelector: utils.GetSelectorFromNameFelt(DECEREMENT_COUNTER),
	}

	callData := account.FmtCallDataCairo2([]rpc.FunctionCall{fnCall})
	invokeTxn := rpc.InvokeTxnV1{
		MaxFee:        new(felt.Felt).SetUint64(0),
		Version:       rpc.TransactionV1,
		Nonce:         nonce,
		Type:          rpc.TransactionType_Invoke,
		SenderAddress: senderAddress,
		Calldata:      callData,
	}

	return invokeTxn, nil
}

func ReadCounter() (rpc.FunctionCall, error) {
	cn, err := utils.HexToFelt(CONTRACT_ADDRESS)
	if err != nil {
		return rpc.FunctionCall{}, fmt.Errorf("invalid contract address, err: %v", err)
	}

	fnCall := rpc.FunctionCall{
		ContractAddress:    cn,
		EntryPointSelector: utils.GetSelectorFromNameFelt(READ_COUNTER),
	}

	return fnCall, nil
}
