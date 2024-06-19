package tokensvc

import (
	"context"
	"encoding/json"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func NewHTTPServer(srv TokenDataService) http.Handler {
	endpoints := newEndpoints(srv)
	r := mux.NewRouter()
	serverOptions := []httptransport.ServerOption{httptransport.ServerErrorEncoder(encodeError)}
	validate := validator.New()
	r.Methods("POST").Path("/token/v1/create").Handler(httptransport.NewServer(
		endpoints.createTokenEndpoint,
		DecodeRequest[CreateTokenReq],
		httptransport.EncodeJSONResponse,
		serverOptions...,
	))

	r.Methods("POST").Path("/token/v1/fetch").Handler(httptransport.NewServer(
		endpoints.fetchTokenEndpoint,
		DecodeRequest[TickerReq],
		httptransport.EncodeJSONResponse,
		serverOptions...,
	))

	r.Methods("GET").Path("/token/v1/fetch/all/{skip}").Handler(httptransport.NewServer(
		endpoints.fetchAllTokenEndpoint,
		DecodePathParams(validate, decodeTokensListRequest),
		httptransport.EncodeJSONResponse,
		serverOptions...,
	))
	return r
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	errInfo := ErrorInfo{
		Code:    500,
		Message: "Internal Server Error",
	}

	if err == TokenExists {
		w.WriteHeader(http.StatusUnauthorized)
		errInfo = ErrorInfo{
			Code:    403,
			Message: TokenExists.Error(),
		}
	} else if err == InsertError {
		w.WriteHeader(http.StatusNotFound)
		errInfo = ErrorInfo{
			Code:    400,
			Message: InsertError.Error(),
		}
	} else if err == BadRequest {
		w.WriteHeader(http.StatusBadRequest)
		errInfo = ErrorInfo{
			Code:    400,
			Message: BadRequest.Error(),
		}
	} else if err == TokenNotFound {
		w.WriteHeader(http.StatusNotFound)
		errInfo = ErrorInfo{
			Code:    404,
			Message: TokenNotFound.Error(),
		}
	}

	switch e := err.(type) {
	case *GenericError:
		w.WriteHeader(e.Code)
		errInfo = ErrorInfo{
			Code:    e.Code,
			Message: e.Error(),
		}
	}

	errResp := Response{
		Success: false,
		Error:   errInfo,
	}
	json.NewEncoder(w).Encode(&errResp)
}
