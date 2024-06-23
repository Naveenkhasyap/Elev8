package tokensvc

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

type CreateTokenReq struct {
	Name               string `json:"name" bson:"name"`
	Ticker             string `json:"ticker" bson:"ticker"`
	Description        string `json:"description" bson:"description"`
	Image              string `json:"image" bson:"image"`
	Twitter            string `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Telegram           string `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Website            string `json:"website,omitempty" bson:"website,omitempty"`
	UserAccountAddress string `json:"userAccountAddress" bson:"userAccountAddress"`
	Status             string `json:"status" bson:"status"`
}

type CreateTokenRes struct {
	Name               string `json:"name"`
	Ticker             string `json:"ticker"`
	UserAccountAddress string `json:"userAccountAddress"`
	TxnHash            string `json:"txnHash"`
	Status             string `json:"status"`
}

type CreateTokenResp struct {
	Success bool `json:"sucess" bson:"sucess"`
}

type OrderDataReq struct {
	Ticker       string `json:"ticker" bson:"ticker"`
	Quantity     string `json:"quantity" bson:"quantity"`
	BuyerAddress string `json:"buyerAddress" bson:"buyerAddress"`
	OrderType    string `json:"orderType" bson:"orderType"`
	IsOwner      bool   `json:"isOwner" bson:"isOwner"`
}
type OrderDataRes struct {
	Success bool `json:"sucess" bson:"sucess"`
}
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type TickerReq struct {
	Ticker string `json:"ticker"`
}

type FetchOrderByUserReq struct {
	UserAccountAddress string `json:"userAccountAddress" bson:"userAccountAddress"`
}

type BuySellTokenReq struct {
	Ticker             string `json:"ticker" bson:"ticker"`
	Quantity           string `json:"quantity" bson:"quantity"`
	UserAccountAddress string `json:"userAccountAddress" bson:"userAccountAddress"`
	OrderType          string `json:"orderType" bson:"orderType"`
}

type BuyTokenRes struct {
	Success bool `json:"sucess" bson:"sucess"`
}

type QuoteReq struct {
	TokenAddress string `json:"tokenAddress" bson:"tokenAddress"`
	Amount       string `json:"amount" bson:"amount"`
}

type OwnerReq struct {
	TokenAddress string `json:"tokenAddress" bson:"tokenAddress"`
}

type BalaceReq struct {
	TokenAddress string `json:"tokenAddress" bson:"tokenAddress"`
}

type DecodeTypes interface {
	TickerReq | CreateTokenReq | TokensListRequest | BuySellTokenReq | OrderDataReq | QuoteReq | BalaceReq | OwnerReq
}

func decodeTokensListRequest(vars map[string]string) (any, error) {
	var req TokensListRequest
	skip, err := strconv.Atoi(vars["skip"])
	if err != nil {
		return nil, &GenericError{
			Code:    400,
			Message: "Invalid Request",
		}
	}
	req.Skip = skip
	return req, nil
}

type ErrorInfo struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func DecodeRequest[T DecodeTypes](ctx context.Context, r *http.Request) (interface{}, error) {
	var request T
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return nil, err
	}
	return request, nil
}

type TokensListRequest struct {
	Skip int
}

func DecodeEmptyreq() httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		return nil, nil
	}
}

func DecodePathParams(validate *validator.Validate, f func(data map[string]string) (any, error)) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		vars := mux.Vars(r)
		request, err := f(vars)
		if err != nil {
			slog.Error("error in processing path params", "err", err)
			return nil, BadRequest
		}
		err = validate.Struct(request)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				slog.Error("invalid validation error", "err", err)
			}

			return nil, BadRequest
		}
		return request, nil
	}
}
