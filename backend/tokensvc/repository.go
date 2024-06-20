package tokensvc

import (
	"context"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TokenDatarepo interface {
	Store(ctx context.Context, tokenData CreateTokenReq) error
	Fetch(ctx context.Context, ticker string) (TokenData, error)
	Update(ctx context.Context, ticker string, tokenData TokenData) error
	FetchAll(ctx context.Context, skip int) ([]TokenData, error)

}
type repo struct {
	dbClient *mongo.Client
}

func NewTokenDatarepo(client *mongo.Client) TokenDatarepo {
	return &repo{
		dbClient: client,
	}
}


func (r repo) Store(ctx context.Context, tokenData CreateTokenReq) error {
	collection := r.dbClient.Database("Assets").Collection("tokens")

	err := collection.FindOne(ctx, bson.M{"ticker": tokenData.Ticker}).Decode(&tokenData)
	if err == nil {
		return TokenExists
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
func (r repo) Fetch(ctx context.Context, ticker string) (TokenData, error) {
	collection := r.dbClient.Database("Assets").Collection("tokens")
	var tokenData TokenData
	err := collection.FindOne(ctx, bson.M{"ticker": ticker}).Decode(&tokenData)
	if err != nil {
		return TokenData{}, err
	}
	return tokenData, err
}


func (r repo) FetchAll(ctx context.Context, skip int) ([]TokenData, error) {
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
	for cursor.Next(ctx) {
		var mode TokenData
		err1 := cursor.Decode(&mode)
		if err1 != nil {
			continue
		}
		tokenList = append(tokenList, mode)
	}
	return tokenList, err
}


func (r repo) Update(ctx context.Context, ticker string, tokenData TokenData) error {
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
