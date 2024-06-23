package tokensvc

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/avast/retry-go"
	starkrpc "github.com/gofiles/internal/clients/stark_rpc"
	"github.com/gofiles/internal/contracts"
	"github.com/gofiles/internal/helpers"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenDataService interface {
	CreateToken(ctx context.Context, tokenData CreateTokenReq) (CreateTokenRes, error)
	FetchToken(ctx context.Context, ticker string) (TokenData, error)
	UpdateToken(ctx context.Context, tokenData TokenData) error
	FetchAllToken(ctx context.Context) ([]TokenData, error)
	BuyToken(ctx context.Context, ticker string, buyData BuySellTokenReq) error
	SellToken(ctx context.Context, ticker string, sellData BuySellTokenReq) error
	FetchAllOrders(ctx context.Context, skip int) ([]OrderData, error)
	FetchOrdersByAddress(ctx context.Context, filter_address string) ([]OrderData, error)
	FetchOrdersByTicker(ctx context.Context, ticker string) ([]OrderData, error)
	FetchTickerData(ctx context.Context) ([]DataPoint, error)
	FetchBalance(ctx context.Context, tokenAddress string) (string, error)
	FetchOwner(ctx context.Context, tokenAddress string) (string, error)
	FetchQuote(ctx context.Context, tokenAddress string, amount string) (string, error)
	FetchRecipt(ctx context.Context, txnHash string) (ReceiptResp, error)
}

type tokenDatasvc struct {
	tokenRepo TokenDatarepo
	deployer  contracts.Deployer
	client    *starkrpc.Provider
}

func NewTokenDataService(repo TokenDatarepo, deployer contracts.Deployer, client *starkrpc.Provider) TokenDataService {
	return &tokenDatasvc{
		tokenRepo: repo,
		deployer:  deployer,
		client:    client,
	}
}

// TODO:
//  1. Storing status for the transaction hash - in mongo ✅
//  2. Allowing creation of token if previous one is rejected ✅
//  3. Polling API for transaction status - Once success have to return with contract address which can be fetched from event logs
//  4. APIs for profiles -
//     a. Tokens created by the user
//     b. Tokens held by the user
//     b. Tokens purchase and sell history
func (svc tokenDatasvc) CreateToken(ctx context.Context, tokenData CreateTokenReq) (CreateTokenRes, error) {
	dbData := TokenData{
		Name:               tokenData.Name,
		Ticker:             tokenData.Ticker,
		Description:        tokenData.Description,
		Website:            tokenData.Website,
		Telegram:           tokenData.Telegram,
		Image:              tokenData.Image,
		Twitter:            tokenData.Telegram,
		UserAccountAddress: tokenData.UserAccountAddress,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		Change24hr:         fmt.Sprintf("%f", 1+rand.Float64()*(30-1)),
		Change7day:         fmt.Sprintf("%f", 1+rand.Float64()*(9-1)),
		Price:              fmt.Sprintf("%f", 1+rand.Float64()*(30-1)),
		MarketCap:          fmt.Sprintf("%d", 1+rand.Intn(1e9)),
	}

	err := svc.tokenRepo.Store(ctx, dbData)
	if err != nil {
		slog.Error("error storing token", "ticker", tokenData.Ticker, "error", err)
		return CreateTokenRes{}, &GenericError{
			Code:    500,
			Message: fmt.Sprintf("unable to store token in db, err: %v", err),
		}
	}

	res, err := svc.deployer.Launch(ctx, contracts.DeployerReq{
		Name:      tokenData.Name,
		Symbol:    tokenData.Ticker,
		Recipient: "0x003e1daE972977e8C7229074059cdF8C33B118B51FEeFd043d3b066551e2EC91",
	})

	if err != nil {
		slog.Error("unable to deploy contract", "err", err)
		return CreateTokenRes{}, &GenericError{
			Code:    500,
			Message: fmt.Sprintf("unable to deploy contract, err: %v", err),
		}
	}

	txnStatusResp, statuserr := svc.getTransactionStatusRetry(ctx, res)
	if statuserr != nil {
		slog.Error("unable to get transaction status", "txnHash", res.String(), "err", statuserr)
		return CreateTokenRes{}, &GenericError{
			Code:    500,
			Message: fmt.Sprintf("unable to fetch transaction status, err: %v", statuserr),
		}
	}

	err = svc.tokenRepo.UpdateToken(ctx, tokenData.Ticker, map[string]string{
		"status":  string(txnStatusResp.FinalityStatus),
		"txnHash": res.String(),
	})

	if err != nil {
		return CreateTokenRes{}, &GenericError{
			Code:    500,
			Message: fmt.Sprintf("unable to update db with staus and txnHash, err: %v", err),
		}
	}

	return CreateTokenRes{
		Name:               tokenData.Name,
		Ticker:             tokenData.Ticker,
		UserAccountAddress: tokenData.UserAccountAddress,
		Status:             string(txnStatusResp.FinalityStatus),
		TxnHash:            res.String(),
	}, nil
}

func (svc tokenDatasvc) FetchToken(ctx context.Context, ticker string) (TokenData, error) {
	res, err := svc.tokenRepo.Fetch(ctx, ticker)
	if err != nil && err == mongo.ErrNoDocuments {
		slog.Error("error fetching token with", " ticker", ticker, "error:", err)
		return TokenData{}, TokenNotFound
	} else if err != nil {
		return TokenData{}, err
	}
	return res, nil
}

func (svc tokenDatasvc) FetchAllToken(ctx context.Context) ([]TokenData, error) {
	res, err := svc.tokenRepo.FetchAllValid(ctx)
	if err != nil {
		slog.Error("error fetching all token ", "err", err)
		return []TokenData{}, &GenericError{
			Code:    500,
			Message: "unable to fetch tokens, please try again",
		}
	}
	return res, nil
}

func (svc tokenDatasvc) UpdateToken(ctx context.Context, tokenData TokenData) error {
	//Todo
	err := svc.tokenRepo.Update(ctx, tokenData.Ticker, tokenData)
	if err != nil {
		slog.Error("error updating token for ticker ", " tokenData.Ticker", tokenData.Ticker, "error:", err)
		return err
	}
	return nil
}

func (svc tokenDatasvc) BuyToken(ctx context.Context, ticker string, buyDataReq BuySellTokenReq) error {
	buyOrder := OrderData{
		Ticker:             buyDataReq.Ticker,
		Quantity:           buyDataReq.Quantity,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		OrderType:          buyDataReq.OrderType,
		UserAccountAddress: buyDataReq.UserAccountAddress,
		IsOwner:            false,
		TokenIn:            "ETH",
		TokenOut:           buyDataReq.Ticker,
	}

	tokenOwner, err := svc.tokenRepo.FetchOwnerofTicker(ctx, buyDataReq.Ticker)
	if err != nil {
		slog.Error("error in fetching token owner", "err", err)
	}
	if strings.EqualFold(tokenOwner, buyDataReq.UserAccountAddress) {
		buyOrder.IsOwner = true
	}

	errbuy := svc.tokenRepo.Buy(ctx, buyOrder)
	if errbuy != nil {
		slog.Error("error updating token for ", "ticker ", buyDataReq.Ticker, "error:", err)
		return err
	}
	return nil
}

func (svc tokenDatasvc) SellToken(ctx context.Context, ticker string, sellDataReq BuySellTokenReq) error {
	//Todo
	sellOrder := OrderData{
		Ticker:             sellDataReq.Ticker,
		Quantity:           sellDataReq.Quantity,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
		OrderType:          sellDataReq.OrderType,
		UserAccountAddress: sellDataReq.UserAccountAddress,
		IsOwner:            false,
		TokenIn:            sellDataReq.Ticker,
		TokenOut:           "ETH",
	}

	tokenOwner, err := svc.tokenRepo.FetchOwnerofTicker(ctx, sellDataReq.Ticker)
	if err != nil {
		slog.Error("error in fetching token owner", "err", err)
	}
	if strings.EqualFold(tokenOwner, sellDataReq.UserAccountAddress) {
		sellOrder.IsOwner = true
	}

	errsell := svc.tokenRepo.Sell(ctx, sellOrder)
	if errsell != nil {
		slog.Error("error updating token for", " ticker ", sellDataReq.Ticker, "error:", err)
		return err
	}
	return nil
}

func (svc tokenDatasvc) FetchAllOrders(ctx context.Context, skip int) ([]OrderData, error) {
	orderList, err := svc.tokenRepo.FetchAllOrders(ctx, skip)
	if err != nil {
		slog.Error("error fetchin all orders ", "err", err)
		return []OrderData{}, err
	}
	return orderList, nil
}

func (svc tokenDatasvc) FetchOrdersByAddress(ctx context.Context, filter_address string) ([]OrderData, error) {
	ordersList, err := svc.tokenRepo.FetchOrderByAddress(ctx, filter_address)
	if err != nil {
		slog.Error("error fetching orders by address ", "filter_address", filter_address, "error:", err)
		return []OrderData{}, err
	}
	return ordersList, nil
}

func (svc tokenDatasvc) FetchOrdersByTicker(ctx context.Context, ticker string) ([]OrderData, error) {
	ordersList, err := svc.tokenRepo.FetchOrderByAddress(ctx, ticker)
	if err != nil {
		slog.Error("error fetching orders for ", "ticker ", ticker, "error:", err)
		return []OrderData{}, err
	}
	return ordersList, nil
}

func (svc tokenDatasvc) FetchTickerData(ctx context.Context) ([]DataPoint, error) {
	var count = 100
	data := make([]DataPoint, count)
	baseTime := time.Now().AddDate(0, 0, -count).Unix()

	for i := 0; i < count; i++ {
		open := rand.Float64() * 100
		close := rand.Float64() * 100
		high := open + rand.Float64()*10
		low := open - rand.Float64()*10
		volume := rand.Float64() * 1000

		data[i] = DataPoint{
			Time:   baseTime + int64(i*3600),
			Open:   open,
			High:   high,
			Low:    low,
			Close:  close,
			Volume: volume,
		}
	}
	return data, nil
}

func (svc tokenDatasvc) FetchBalance(ctx context.Context, tokenAddress string) (string, error) {
	contractAddress, err := utils.HexToFelt(tokenAddress)
	var contractMethod = "get_balance"
	if err != nil {
		panic(err)
	}

	// Make read contract call
	tx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod),
	}

	fmt.Println("Making Call() request")
	callResp, _ := svc.client.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	resp := utils.FeltToBigInt(callResp[0])

	fmt.Println(fmt.Sprintf("Response to %s():%s ", contractMethod, resp.String()))
	return resp.String(), nil
}

func (svc tokenDatasvc) FetchQuote(ctx context.Context, tokenAddress string, amount string) (string, error) {
	var contractMethod = "quote"
	contractAddress, err := utils.HexToFelt(tokenAddress)
	if err != nil {
		panic(err)
	}
	amountf, err := utils.HexToFelt(amount)
	if err != nil {
		return "", fmt.Errorf("unable to generate calldata for amount bytearray: %s, err: %v", amount, err)
	}
	stylef, err := utils.HexToFelt("1")
	if err != nil {
		return "", fmt.Errorf("unable to generate calldata for style bytearray: %s, err: %v", err)
	}

	sa, err := utils.HexToFelt("0")
	if err != nil {
		return "", fmt.Errorf("unable to generate calldata for sa bytearray:, err: %v", err)
	}

	// Make read contract call
	tx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod),
		Calldata:           []*felt.Felt{amountf, sa, stylef},
	}

	fmt.Println("Making Call() request")
	callResp, _ := svc.client.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})
	resp := utils.FeltToBigInt(callResp[0])

	fmt.Println(fmt.Sprintf("Response to %s():%s ", contractMethod, resp.String()))
	return resp.String(), nil
}

func (svc tokenDatasvc) FetchOwner(ctx context.Context, tokenAddress string) (string, error) {
	contractAddress, err := utils.HexToFelt(tokenAddress)
	var contractMethod = "owner"
	if err != nil {
		panic(err)
	}

	// Make read contract call
	tx := rpc.FunctionCall{
		ContractAddress:    contractAddress,
		EntryPointSelector: utils.GetSelectorFromNameFelt(contractMethod),
	}

	fmt.Println("Making Call() request")
	callResp, _ := svc.client.Call(context.Background(), tx, rpc.BlockID{Tag: "latest"})

	fmt.Println(fmt.Sprintf("Response to %s():%s ", contractMethod, callResp[0]))
	return fmt.Sprintf("%s", callResp[0]), nil
}

func (svc tokenDatasvc) getTransactionStatusRetry(ctx context.Context, txnHash *felt.Felt) (*rpc.TxnStatusResp, error) {
	var out *rpc.TxnStatusResp
	retry.Do(func() error {
		res, err := svc.client.GetTransactionStatus(ctx, txnHash)
		if err != nil {
			return err
		}
		out = res
		return nil
	}, retry.Delay(1*time.Second), retry.Attempts(5))
	return out, nil
}

func (svc tokenDatasvc) FetchRecipt(ctx context.Context, txnHash string) (ReceiptResp, error) {
	var receiptResp = ReceiptResp{}
	hash, err := utils.HexToFelt(txnHash)
	if err != nil {
		panic(err)
	}
	resp, err := svc.client.WaitForTransaction(ctx, hash)
	if err != nil {
		return ReceiptResp{}, err
	}

	for _, val := range resp.Result.Events {
		if len(val.Keys) == 2 && len(val.Data) == 1 {
			receiptResp.TokenAddress = fmt.Sprintf("%s", val.Data[0])
		}
		if len(val.Data) == 16 {
			receiptResp.Ticker = helpers.HexToString(fmt.Sprintf("%s", val.Data[6]))
		}
	}

	return receiptResp, nil
}
