package tokensvc

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

func NewHTTPServer(srv TokenDataService) http.Handler {
	endpoints := newEndpoints(srv)
	r := mux.NewRouter()
	validate := validator.New()
	r.Methods("POST").Path("/token/v1/create").Handler(httptransport.NewServer(
		endpoints.createTokenEndpoint,
		DecodeRequest[CreateTokenReq],
		httptransport.EncodeJSONResponse,
	))

	r.Methods("POST").Path("/token/v1/fetch").Handler(httptransport.NewServer(
		endpoints.fetchTokenEndpoint,
		DecodeRequest[TickerReq],
		httptransport.EncodeJSONResponse,
	))

	r.Methods("GET").Path("/token/v1/fetch/all/{skip}").Handler(httptransport.NewServer(
		endpoints.fetchAllTokenEndpoint,
		DecodePathParams(validate, decodeTokensListRequest),
		httptransport.EncodeJSONResponse,
	))
	return r
}
