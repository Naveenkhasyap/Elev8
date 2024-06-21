package tokensvc

import (
	"context"
	"fmt"
	"log/slog"

	starkrpc "github.com/gofiles/internal/clients/stark_rpc"
	"github.com/gofiles/internal/contracts"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenDataService interface {
	CreateToken(ctx context.Context, tokenData CreateTokenReq) (CreateTokenRes, error)
	FetchToken(ctx context.Context, ticker string) (TokenData, error)
	UpdateToken(ctx context.Context, tokenData TokenData) error
	SellToken(ctx context.Context, ticker string, owner_address string) error
	FetchAllToken(ctx context.Context, skip int) ([]TokenData, error)
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
//  1. Storing status for the transaction hash - in mongo
//  2. Allowing creation of token if previous one is rejected
//  3. Polling API for transaction status - Once success have to return with contract address which can be fetched from event logs
//  4. APIs for profiles -
//     a. Tokens created by the user
//     b. Tokens held by the user
//     b. Tokens purchase and sell history
func (svc tokenDatasvc) CreateToken(ctx context.Context, tokenData CreateTokenReq) (CreateTokenRes, error) {
	err := svc.tokenRepo.Store(ctx, tokenData)
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
		return CreateTokenRes{}, fmt.Errorf("unable to deploy contract, err: %v", err)
	}

	_, statuserr := svc.client.GetTransactionStatus(ctx, res)
	if statuserr != nil {
		return CreateTokenRes{}, fmt.Errorf("unable to get transaction status, err: %v", statuserr)
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
func (svc tokenDatasvc) SellToken(ctx context.Context, ticker string, owner_address string) error {
	//Todo
	return nil
}
