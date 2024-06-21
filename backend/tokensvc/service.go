package tokensvc

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"strings"
	"time"

	starkrpc "github.com/gofiles/internal/clients/stark_rpc"
	"github.com/gofiles/internal/contracts"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenDataService interface {
	CreateToken(ctx context.Context, tokenData CreateTokenReq) (CreateTokenRes, error)
	FetchToken(ctx context.Context, ticker string) (TokenData, error)
	UpdateToken(ctx context.Context, tokenData TokenData) error
	FetchAllToken(ctx context.Context, skip int) ([]TokenData, error)
	BuyToken(ctx context.Context, ticker string, buyData BuySellTokenReq) error
	SellToken(ctx context.Context, ticker string, sellData BuySellTokenReq) error
	FetchAllOrders(ctx context.Context, skip int) ([]OrderData, error)
	FetchOrdersByAddress(ctx context.Context, filter_address string) ([]OrderData, error)
	FetchOrdersByTicker(ctx context.Context, ticker string) ([]OrderData, error)
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
		Change24hr:         fmt.Sprintf("%f", 1+rand.Float64()*(30-1)),
		Change7day:         fmt.Sprintf("%f", 1+rand.Float64()*(9-1)),
	}
	err := svc.tokenRepo.Store(ctx, dbData)
	if err != nil {
		slog.Error("error while storing token ", tokenData.Ticker, "error:", err)
		return CreateTokenRes{}, err
	}

	res, err := svc.deployer.Launch(ctx, contracts.DeployerReq{
		Name:      tokenData.Name,
		Symbol:    tokenData.Ticker,
		Recipient: "0x003e1daE972977e8C7229074059cdF8C33B118B51FEeFd043d3b066551e2EC91",
	})

	if err != nil {
		slog.Error("unable to deploy contract, err: ", err)
		return CreateTokenRes{}, DeployError
	}

	txnStatusResp, statuserr := svc.client.GetTransactionStatus(ctx, res)
	if statuserr != nil {
		slog.Error("unable to get transaction status, err: ", statuserr)
		return CreateTokenRes{}, FetchTxnStatusError
	}
	slog.Info("FinalityStatus", txnStatusResp.FinalityStatus, "ExecutionStatus:", txnStatusResp.ExecutionStatus, "txn receipt", res)

	//updating status
	errUpdate := svc.tokenRepo.UpdateStatus(ctx, tokenData.Ticker, string(txnStatusResp.FinalityStatus))
	if errUpdate != nil {
		slog.Error("error in updating status", errUpdate)
	}
	//updating hash
	errtxn := svc.tokenRepo.UpdateTxnHash(ctx, tokenData.Ticker, fmt.Sprintf("%v", res))
	if errtxn != nil {
		slog.Error("error in updating status", errUpdate)
	}

	return CreateTokenRes{
		Name: tokenData.Name,
	}, nil
}

func (svc tokenDatasvc) FetchToken(ctx context.Context, ticker string) (TokenData, error) {
	res, err := svc.tokenRepo.Fetch(ctx, ticker)
	if err != nil && err == mongo.ErrNoDocuments {
		slog.Error("error fetching token with ticker", ticker, "error:", err)
		return TokenData{}, TokenNotFound
	} else if err != nil {
		return TokenData{}, err
	}
	return res, nil
}

func (svc tokenDatasvc) FetchAllToken(ctx context.Context, skip int) ([]TokenData, error) {
	//Todo
	res, err := svc.tokenRepo.FetchAll(ctx, skip)
	if err != nil {
		slog.Error("error fetching all token ", err)
		return []TokenData{}, err
	}
	return res, nil
}
func (svc tokenDatasvc) UpdateToken(ctx context.Context, tokenData TokenData) error {
	//Todo
	err := svc.tokenRepo.Update(ctx, tokenData.Ticker, tokenData)
	if err != nil {
		slog.Error("error updating token for ticker ", tokenData.Ticker, "error:", err)
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
		slog.Error("error in fetching token owner", err)
	}
	if strings.EqualFold(tokenOwner, buyDataReq.UserAccountAddress) {
		buyOrder.IsOwner = true
	}

	errbuy := svc.tokenRepo.Buy(ctx, buyOrder)
	if errbuy != nil {
		slog.Error("error updating token for ticker ", buyDataReq.Ticker, "error:", err)
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
		slog.Error("error in fetching token owner", err)
	}
	if strings.EqualFold(tokenOwner, sellDataReq.UserAccountAddress) {
		sellOrder.IsOwner = true
	}

	errsell := svc.tokenRepo.Sell(ctx, sellOrder)
	if errsell != nil {
		slog.Error("error updating token for ticker ", sellDataReq.Ticker, "error:", err)
		return err
	}
	return nil
}
func (svc tokenDatasvc) FetchAllOrders(ctx context.Context, skip int) ([]OrderData, error) {
	orderList, err := svc.tokenRepo.FetchAllOrders(ctx, skip)
	if err != nil {
		slog.Error("error fetchin all orders ", err)
		return []OrderData{}, err
	}

	return orderList, nil

}
func (svc tokenDatasvc) FetchOrdersByAddress(ctx context.Context, filter_address string) ([]OrderData, error) {
	ordersList, err := svc.tokenRepo.FetchOrderByAddress(ctx, filter_address)
	if err != nil {
		slog.Error("error fetching orders by address ", filter_address, "error:", err)
		return []OrderData{}, err
	}

	return ordersList, nil
}

func (svc tokenDatasvc) FetchOrdersByTicker(ctx context.Context, ticker string) ([]OrderData, error) {
	ordersList, err := svc.tokenRepo.FetchOrderByAddress(ctx, ticker)
	if err != nil {
		slog.Error("error fetching orders for ticker ", ticker, "error:", err)
		return []OrderData{}, err
	}

	return ordersList, nil
}
