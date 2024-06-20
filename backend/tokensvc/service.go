package tokensvc

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
)

type TokenDataService interface {
	CreateToken(ctx context.Context, tokenData CreateTokenReq) error
	FetchToken(ctx context.Context, ticker string) (TokenData, error)
	UpdateToken(ctx context.Context, tokenData TokenData) error
	SellToken(ctx context.Context, ticker string, owner_address string) error
	FetchAllToken(ctx context.Context, skip int) ([]TokenData, error)

}
type tokenDatasvc struct {
	tokenRepo TokenDatarepo
}

func NewTokenDataService(repo TokenDatarepo) TokenDataService {
	return &tokenDatasvc{
		tokenRepo: repo,
	}
}

func (svc tokenDatasvc) CreateToken(ctx context.Context, tokenData CreateTokenReq) error {
	//Todo
	err := svc.tokenRepo.Store(ctx, tokenData)
	if err != nil {
		slog.Error("error while storing token ", tokenData.Ticker, "error:", err)
		return err
	}

	return nil
}

func (svc tokenDatasvc) FetchToken(ctx context.Context, ticker string) (TokenData, error) {
	//Todo
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
