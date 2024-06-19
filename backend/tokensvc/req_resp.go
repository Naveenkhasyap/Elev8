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
	Name          string `json:"name" bson:"name"`
	Ticker        string `json:"ticker" bson:"ticker"`
	Description   string `json:"description" bson:"description"`
	Image         string `json:"image" bson:"image"`
	Twitter       string `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Telegram      string `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Website       string `json:"website,omitempty" bson:"website,omitempty"`
	WalletAddress string `json:"wallet" bson:"wallet"`
}

type CreateTokenResp struct {
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

type DecodeTypes interface {
	TickerReq | CreateTokenReq | TokensListRequest
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

func DecodePathParams(validate *validator.Validate, f func(data map[string]string) (any, error)) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		vars := mux.Vars(r)
		request, err := f(vars)
		if err != nil {
			slog.Error("err", err, "msg", "error in processing path params")
			return nil, BadRequest
		}
		err = validate.Struct(request)
		if err != nil {
			if _, ok := err.(*validator.InvalidValidationError); ok {
				slog.Error("err", err, "msg", "invalid validation error")
			}

			return nil, BadRequest
		}
		return request, nil
	}
}
