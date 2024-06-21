package tokensvc

import "errors"

type TokenData struct {
	Name            string `json:"name" bson:"name"`
	Ticker          string `json:"ticker" bson:"ticker"`
	Description     string `json:"description" bson:"description"`
	Image           string `json:"image" bson:"image"`
	Twitter         string `json:"twitter,omitempty" bson:"twitter,omitempty"`
	Telegram        string `json:"telegram,omitempty" bson:"telegram,omitempty"`
	Website         string `json:"website,omitempty" bson:"website,omitempty"`
	WalletAddress   string `json:"wallet" bson:"wallet"`
	Status          string `json:"status" bson:"status"`
	TransactionHash string `json:"txnHash" bson:"txnHash"`
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
)
