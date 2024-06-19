package tokensvc

import "context"

type TokenDataService interface {
	CreateToken(ctx context.Context, tokenData TokenData) error
	FetchToken(ctx context.Context, ticker string) (TokenData, error)
	UpdateToken(ctx context.Context, tokenData TokenData) error
	SellToken(ctx context.Context, ticker string, owner_address string) error
}
type tokenDatasvc struct {
	tokenRepo TokenDatarepo
}

func NewTokenDataService(repo TokenDatarepo) TokenDataService {
	return &tokenDatasvc{
		tokenRepo: repo,
	}
}

func (svc tokenDatasvc) CreateToken(ctx context.Context, tokenData TokenData) error {
	//Todo
	return nil
}
func (svc tokenDatasvc) FetchToken(ctx context.Context, ticker string) (TokenData, error) {
	//Todo
	return TokenData{}, nil
}
func (svc tokenDatasvc) UpdateToken(ctx context.Context, tokenData TokenData) error {
	//Todo
	return nil
}
func (svc tokenDatasvc) SellToken(ctx context.Context, ticker string, owner_address string) error {
	//Todo
	return nil
}
