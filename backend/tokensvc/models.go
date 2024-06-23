package tokensvc

import (
	"errors"
	"time"
)

type TokenData struct {
	Name               string    `json:"name" bson:"name"`
	Ticker             string    `json:"ticker" bson:"ticker"`
	Description        string    `json:"description" bson:"description"`
	Image              string    `json:"image" bson:"image"`
	Twitter            string    `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Telegram           string    `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Website            string    `json:"website,omitempty" bson:"website,omitempty"`
	UserAccountAddress string    `json:"userAccountAddress" bson:"userAccountAddress"`
	Status             string    `json:"status" bson:"status"`
	TransactionHash    string    `json:"txnHash" bson:"txnHash"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" bson:"updated_at"`
	Change24hr         string    `json:"change24hr" bson:"change24hr"`
	Change7day         string    `json:"change7day" bson:"change7day"`
	Price              string    `json:"price" bson:"price"`
	MarketCap          string    `json:"marketCap" bson:"marketcap"`
}

type OrderData struct {
	Ticker             string    `json:"ticker" bson:"ticker"`
	Quantity           string    `json:"quantity" bson:"quantity"`
	UserAccountAddress string    `json:"userAccountAddress" bson:"userAccountAddress"`
	OrderType          string    `json:"orderType" bson:"orderType"`
	IsOwner            bool      `json:"isOwner" bson:"isOwner"`
	CreatedAt          time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" bson:"created_at"`
	TransactionHash    string    `json:"txnHash" bson:"txnHash"`
	Status             string    `json:"status" bson:"status"`
	TokenIn            string    `json:"tokenIn" bson:"tokenIn"`
	TokenOut           string    `json:"tokenOut" bson:"tokenOut"`
}
type DataPoint struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}
type ReceiptEvent struct {
	From string   `json:"from_address" bson:"from_address"`
	Keys []string `json:"keys" bson:"keys"`
	Data []string `json:"data" bson:"data"`
}
type GenericError struct {
	Code    int
	Message string
}

func (e *GenericError) Error() string {
	return e.Message
}

var (
	ErrDuplicateBankAc    = errors.New("Account already verified")
	BadRequest            = errors.New("Bad Request")
	ErrInvalidReferenceId = errors.New("Invalid reference Id")
	ErrInvalidHash        = errors.New("Invalid transaction hash")
	ErrUnauthorized       = errors.New("Unauthorized request")
	TokenNotFound         = errors.New("NO DATA FOUND")
	ErrForbidden          = errors.New("User is Forbidden")
	TokenExists           = errors.New("token already exists")
	InsertError           = errors.New("error inserting")
	DeployError           = errors.New("unable to deploy contract")
	FetchTxnStatusError   = errors.New("unable to get transaction status")
	UpdateTxnStatusError  = errors.New("db update transaction status error")
)
