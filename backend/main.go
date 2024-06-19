package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	tokensvc "github.com/gofiles/tokensvc"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Error("error loading env")
		os.Exit(-1)
	}

	client, err := initMongoConnections(context.TODO())
	if err != nil {
		slog.Error("error mongo init")
		os.Exit(-1)
	}

	imageSvc := initTokenService(client)

	errs := make(chan error)
	sm := http.NewServeMux()
	sm.Handle("/token/create", imageSvc)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errs <- http.ListenAndServe(":8080", sm)
	}()

	select {
	case _, ok := <-errs:
		if ok {
			slog.Error("errs", "error from Server")
			os.Exit(-1)
		} else {
			slog.Error("errs", "exit channel empty")
		}
	}

}

func initMongoConnections(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		slog.Error("error in step 1")
		return nil, err
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		slog.Error("error in step 2")
		return nil, err
	}
	return client, nil
}

func initTokenService(client *mongo.Client) http.Handler {
	//Todo : move thos logic to service
	fn := func(w http.ResponseWriter, r *http.Request) {
		collection := client.Database("Assets").Collection("tokens")
		tokenData := tokensvc.TokenData{
			Name:        "Token",
			Ticker:      "TKN",
			Description: "some meme token",
			Image:       "fefref",
		}
		result, err := collection.InsertOne(context.TODO(), tokenData)
		if err != nil {
			slog.Error("some error ", err)
		}

		_, ok := result.InsertedID.(primitive.ObjectID)

		if !ok {
			slog.Error("NOT A VALID INSERT RESULT")
		}

		w.Write([]byte("Inserted token successfully "))
	}
	return http.HandlerFunc(fn)

}
