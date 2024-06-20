package contracts

import (
	"encoding/hex"
	"fmt"
	"regexp"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/utils"
)

const SHORT_LENGTH = 31

func GenerateCallDataforByteArray(s string) ([]*felt.Felt, error) {
	arr, err := SplitLongString(s)
	if err != nil {
		return nil, err
	}

	hexarr := []string{}
	count := 0
	size := 0
	for _, val := range arr {
		if len(val) == 31 {
			count += 1
		}
		size = len(val)
		hexarr = append(hexarr, StringToHex(val))
	}

	harr, err := utils.HexArrToFelt(hexarr)
	if err != nil {
		return nil, err
	}

	res := append([]*felt.Felt{new(felt.Felt).SetUint64(uint64(count))}, harr...)
	return append(res, new(felt.Felt).SetUint64(uint64(size))), nil
}

func SplitLongString(s string) ([]string, error) {
	exp := fmt.Sprintf(".{1,%d}", SHORT_LENGTH)
	r, err := regexp.Compile(exp)
	if err != nil {
		return []string{}, fmt.Errorf("invalid regex, err: %v", err)
	}
	res := r.FindAllString(s, -1)
	if len(res) == 0 {
		return []string{}, fmt.Errorf("invalid string no regex matches found, s: %s", s)
	}
	return res, nil
}

func StringToHex(s string) string {
	bytes := []byte(s)
	hexString := hex.EncodeToString(bytes)
	return hexString
}
