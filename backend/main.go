package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/NethermindEth/juno/core/felt"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/gofiles/internal/accounts"
	starkrpc "github.com/gofiles/internal/clients/stark_rpc"
	"github.com/gofiles/internal/contracts"
	tokensvc "github.com/gofiles/tokensvc"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Error("error loading env", "err", err)
		os.Exit(-1)
	}

	rpcx, rpcerr := rpc.NewProvider(os.Getenv("RPC_URL"))
	if rpcerr != nil {
		slog.Error("error rpc init", "err", err)
		os.Exit(-1)
	}

	rpcClient := initStarkRPC(os.Getenv("RPC_URL"), rpcx)
	client, err := initMongoConnections(context.TODO())
	if err != nil {
		slog.Error("error mongo init", "err", err)
		os.Exit(-1)
	}
	fmt.Println("mongo connection successful")

	account, err := initAccount(rpcClient)
	if err != nil {
		slog.Error("error account init", "err", err)
		os.Exit(-1)
	}

	deployer, err := initDeployer(rpcClient, account)
	if err != nil {
		slog.Error("error deployer init", "err", err)
		os.Exit(-1)
	}

	tokenService := initTokenService(client, deployer, rpcClient) //init service
	tokenHandler := tokensvc.NewHTTPServer(tokenService)          //init handler

	errs := make(chan error)
	sm := http.NewServeMux()
	sm.Handle("/token/v1/", tokenHandler)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		errs <- http.ListenAndServe(":8080", sm)
	}()

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	//poll orders
	go func() {
		for {
			select {
			case <-ticker.C:
				pollOrderforStatus(context.TODO(), client, rpcClient)
			case <-quit:
				ticker.Stop()
				return
			}
		}

	}()

	err, ok := <-errs
	if ok {
		slog.Error("error from Server", "err", err)
		os.Exit(-1)
	} else {
		slog.Error("exit channel empty", "err", err)
	}

}

func pollOrderforStatus(ctx context.Context, client *mongo.Client, rpcClient *starkrpc.Provider) {
	slog.Info("msg", "info", "polling orders")
	tokenCollection := client.Database("Assets").Collection("tokens")
	cursor, err := tokenCollection.Find(ctx, bson.M{"status": bson.M{"$eq": ""}})
	if err != nil {
		slog.Error("err in finding poll data", "err", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		slog.Info("msg", "info", "looping now")
		var mode tokensvc.TokenData
		err1 := cursor.Decode(&mode)
		if err1 != nil {
			continue
		}
		if mode.TransactionHash != "" {
			hash, err := new(felt.Felt).SetString(mode.TransactionHash)
			if err != nil {
				slog.Error("error converting hash to felt", "err", err)
			}
			resp, _ := rpcClient.GetTransactionStatus(ctx, hash)
			slog.Info("Transaction status", "resp", resp)
			if resp == nil {
				continue
			}
			if !strings.EqualFold(string(resp.FinalityStatus), mode.Status) {
				_, err := tokenCollection.UpdateOne(ctx, bson.M{"ticker": bson.M{"$eq": mode.Ticker}},
					bson.M{"$set": bson.M{
						"status": string(resp.FinalityStatus),
					}})
				if err != nil {
					slog.Error("error updating status in go routine", "err", err)
				}
			}
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

func initTokenService(client *mongo.Client, deployer contracts.Deployer, rpc *starkrpc.Provider) tokensvc.TokenDataService {
	repo := tokensvc.NewTokenDatarepo(client)
	return tokensvc.NewTokenDataService(repo, deployer, rpc)
}

func initStarkRPC(url string, client *rpc.Provider) *starkrpc.Provider {
	return starkrpc.NewProvider(url, client)
}

func initAccount(client *starkrpc.Provider) (accounts.IAccount, error) {
	accountAddress := os.Getenv("ACCOUNT_ADDRESS")
	privateKey := os.Getenv("PRIVATE_KEY")
	return accounts.NewAccount(client, accountAddress, privateKey)
}

func initDeployer(client *starkrpc.Provider, la accounts.IAccount) (contracts.Deployer, error) {
	contractAddress := os.Getenv("DEPLOYER_CONTRACT_ADDRESS")
	return contracts.NewDeployer(contractAddress, client, la)
}
