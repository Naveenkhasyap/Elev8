package tokensvc

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	createTokenEndpoint   endpoint.Endpoint
	fetchTokenEndpoint    endpoint.Endpoint
	fetchAllTokenEndpoint endpoint.Endpoint
}

func newEndpoints(s TokenDataService) Endpoints {
	return Endpoints{
		createTokenEndpoint:   makecreateTokenEndpoint(s),
		fetchTokenEndpoint:    makefetchTokenEndpoint(s),
		fetchAllTokenEndpoint: makefetchAllTokenEndpoint(s),
	}

}

func makecreateTokenEndpoint(s TokenDataService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(CreateTokenReq)
		if !ok {
			return nil, &GenericError{
				Code:    400,
				Message: "Bad Request",
			}
		}

		err := s.CreateToken(ctx, req)
		success := err == nil
		return Response{
			Success: success,
			Data:    map[string]string{},
		}, err
	}
}

func makefetchAllTokenEndpoint(s TokenDataService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok1 := request.(TokensListRequest)
		if !ok1 {
			return nil, errors.New("failed to create req fetchCustomerOrdersEndpoint")
		}

		tokenList, err := s.FetchAllToken(ctx, req.Skip)
		success := err == nil
		return Response{
			Success: success,
			Data:    tokenList,
		}, err
	}
}

func makefetchTokenEndpoint(s TokenDataService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(TickerReq)
		if !ok {
			return nil, &GenericError{
				Code:    400,
				Message: "Bad Request",
			}
		}
		tokenList, err := s.FetchToken(ctx, req.Ticker)
		success := err == nil
		return Response{
			Success: success,
			Data:    tokenList,
		}, err
	}
}
