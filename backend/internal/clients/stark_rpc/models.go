package starkrpc

type TransactionReceipt struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		Type            string `json:"type"`
		TransactionHash string `json:"transaction_hash"`
		ActualFee       struct {
			Amount string `json:"amount"`
			Unit   string `json:"unit"`
		} `json:"actual_fee"`
		ExecutionStatus string        `json:"execution_status"`
		FinalityStatus  string        `json:"finality_status"`
		BlockHash       string        `json:"block_hash"`
		BlockNumber     int           `json:"block_number"`
		MessagesSent    []interface{} `json:"messages_sent"`
		Events          []struct {
			FromAddress string        `json:"from_address"`
			Keys        []string      `json:"keys"`
			Data        []interface{} `json:"data"`
		} `json:"events"`
		ExecutionResources struct {
			Steps                         int `json:"steps"`
			PedersenBuiltinApplications   int `json:"pedersen_builtin_applications"`
			RangeCheckBuiltinApplications int `json:"range_check_builtin_applications"`
			EcOpBuiltinApplications       int `json:"ec_op_builtin_applications"`
			PoseidonBuiltinApplications   int `json:"poseidon_builtin_applications"`
			DataAvailability              struct {
				L1Gas     int `json:"l1_gas"`
				L1DataGas int `json:"l1_data_gas"`
			} `json:"data_availability"`
		} `json:"execution_resources"`
	} `json:"result"`
	ID int `json:"id"`
}
