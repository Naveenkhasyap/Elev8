package tokensvc

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenDatarepo interface {
	Store(ctx context.Context, tokenData TokenData) error
	Fetch(ctx context.Context, ticker string) (TokenData, error)
	Update(ctx context.Context, ticker string, tokenData TokenData) error
}
type repo struct {
	dbClient *mongo.Client
}

func NewTokenDatarepo(client *mongo.Client) TokenDatarepo {
	return &repo{
		dbClient: client,
	}
}

func (r repo) Store(ctx context.Context, tokenData TokenData) error {
	collection := r.dbClient.Database("Assets").Collection("tokens")
	res, err := collection.InsertOne(ctx, tokenData)
	if err != nil {
		return err
	}
	_, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return errors.New("error inserting")
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