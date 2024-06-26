package tokensvc

import (
	"context"

	"github.com/NethermindEth/starknet.go/rpc"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenDatarepo interface {
	Store(ctx context.Context, tokenData TokenData) error
	Fetch(ctx context.Context, ticker string) (TokenData, error)
	FetchByTxnHash(ctx context.Context, txnHash string) (TokenData, error)
	Update(ctx context.Context, ticker string, tokenData TokenData) error
	UpdateToken(ctx context.Context, ticker string, body map[string]string) error
	FetchAll(ctx context.Context, skip int) ([]TokenData, error)
	FetchAllValid(ctx context.Context) ([]TokenData, error)
	Buy(ctx context.Context, orderData OrderData) error
	Sell(ctx context.Context, orderData OrderData) error
	FetchAllOrders(ctx context.Context, skip int) ([]OrderData, error)
	FetchOwnerofTicker(ctx context.Context, ticker string) (string, error)
	FetchOrderByAddress(ctx context.Context, address string) ([]OrderData, error)
	FetchOrderByTicker(ctx context.Context, ticker string) ([]OrderData, error)
}

type repo struct {
	dbClient *mongo.Client
}

func NewTokenDatarepo(client *mongo.Client) TokenDatarepo {
	return &repo{
		dbClient: client,
	}
}

func (r *repo) Store(ctx context.Context, tokenData TokenData) error {
	collection := r.dbClient.Database("Assets").Collection("tokens")
	var existingData TokenData
	err := collection.FindOne(ctx, bson.M{"ticker": tokenData.Ticker}).Decode(&existingData)
	if err == nil {
		if existingData.Status != "REJECTED" {
			return TokenExists
		}
	}

	res, err := collection.InsertOne(ctx, tokenData)
	if err != nil {
		return err
	}

	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return InsertError
	}
	return nil
}

func (r *repo) Fetch(ctx context.Context, ticker string) (TokenData, error) {
	collection := r.dbClient.Database("Assets").Collection("tokens")
	var tokenData TokenData
	err := collection.FindOne(ctx, bson.M{"ticker": ticker}).Decode(&tokenData)
	return tokenData, err
}

func (r *repo) FetchAll(ctx context.Context, skip int) ([]TokenData, error) {
	var tokenList = []TokenData{}

	opt := options.Find()
	opt = opt.SetLimit(int64(10))
	opt = opt.SetSkip(int64(skip * 10))
	opt = opt.SetSort(bson.M{"_id": -1})

	collection := r.dbClient.Database("Assets").Collection("tokens")

	cursor, err := collection.Find(ctx, bson.M{}, opt)
	if err != nil {
		return []TokenData{}, err
	}

	err = cursor.All(ctx, &tokenList)
	return tokenList, err
}

func (r *repo) Update(ctx context.Context, ticker string, tokenData TokenData) error {
	collection := r.dbClient.Database("Assets").Collection("tokens")
	_, err := collection.UpdateOne(ctx, bson.M{"ticker": bson.M{"$eq": ticker}},
		bson.M{"$set": bson.M{
			"name":        tokenData.Name,
			"descriptio ": tokenData.Description,
			"image":       tokenData.Image,
			"twitter":     tokenData.Twitter,
			"telegram":    tokenData.Telegram,
			"website":     tokenData.Website,
		}})
	return err
}

func (r *repo) UpdateToken(ctx context.Context, ticker string, body map[string]string) error {
	collection := r.dbClient.Database("Assets").Collection("tokens")
	updateBody := bson.M{}
	for k, v := range body {
		updateBody[k] = v
	}
	_, err := collection.UpdateOne(ctx, bson.M{"ticker": bson.M{"$eq": ticker}}, bson.M{"$set": updateBody})
	return err
}

func (r *repo) FetchOwnerofTicker(ctx context.Context, ticker string) (string, error) {
	collection := r.dbClient.Database("Assets").Collection("tokens")
	var tokenData TokenData
	err := collection.FindOne(ctx, bson.M{"ticker": ticker}).Decode(&tokenData)
	if err != nil {
		return "", err
	}
	return tokenData.UserAccountAddress, err
}

func (r *repo) Buy(ctx context.Context, orderData OrderData) error {
	collection := r.dbClient.Database("Assets").Collection("orders")
	res, err := collection.InsertOne(ctx, orderData)
	if err != nil {
		return err
	}
	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return InsertError
	}
	return nil
}

func (r *repo) Sell(ctx context.Context, orderData OrderData) error {
	collection := r.dbClient.Database("Assets").Collection("orders")
	res, err := collection.InsertOne(ctx, orderData)
	if err != nil {
		return err
	}
	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return InsertError
	}
	return nil
}

func (r *repo) FetchOrderByTicker(ctx context.Context, ticker string) ([]OrderData, error) {
	collection := r.dbClient.Database("Assets").Collection("orders")
	var orderData []OrderData
	err := collection.FindOne(ctx, bson.M{"ticker": ticker}).Decode(&orderData)
	return orderData, err
}

func (r *repo) FetchOrderByAddress(ctx context.Context, address string) ([]OrderData, error) {
	collection := r.dbClient.Database("Assets").Collection("orders")
	var orderData []OrderData
	err := collection.FindOne(ctx, bson.M{"userAccountAddress": address}).Decode(&orderData)
	return orderData, err
}

func (r *repo) FetchAllOrders(ctx context.Context, skip int) ([]OrderData, error) {
	collection := r.dbClient.Database("Assets").Collection("orders")
	var orderList = []OrderData{}

	opt := options.Find()
	opt = opt.SetLimit(int64(10))
	opt = opt.SetSkip(int64(skip * 10))
	opt = opt.SetSort(bson.M{"_id": -1})

	cursor, err := collection.Find(ctx, bson.M{}, opt)
	if err != nil {
		return []OrderData{}, err
	}

	err = cursor.All(ctx, &orderList)
	return orderList, err
}

func (r *repo) FetchAllValid(ctx context.Context) ([]TokenData, error) {
	collection := r.dbClient.Database("Assets").Collection("tokens")
	var tokenData []TokenData

	opt := options.Find()
	opt = opt.SetSort(bson.M{"_id": -1})
	cursor, err := collection.Find(ctx, bson.M{"$or": bson.A{
		bson.M{"status": string(rpc.BlockStatus_AcceptedOnL1)},
		bson.M{"status": string(rpc.BlockStatus_AcceptedOnL2)},
		bson.M{"status": "RECEIVED"},
	}}, opt)

	if err != nil {
		return []TokenData{}, err
	}

	err = cursor.All(ctx, &tokenData)
	return tokenData, err
}

func (r *repo) FetchByTxnHash(ctx context.Context, txnHash string) (TokenData, error) {
	collection := r.dbClient.Database("Assets").Collection("tokens")
	var tokenData TokenData
	err := collection.FindOne(ctx, bson.M{"txnHash": txnHash}).Decode(&tokenData)
	if err != nil {
		return TokenData{}, err
	}
	return tokenData, err
}
