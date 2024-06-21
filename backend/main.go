package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/gofiles/internal/accounts"
	starkrpc "github.com/gofiles/internal/clients/stark_rpc"
	"github.com/gofiles/internal/contracts"
	tokensvc "github.com/gofiles/tokensvc"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Error("error loading env")
		os.Exit(-1)
	}

	rpcx, rpcerr := rpc.NewProvider(os.Getenv("RPC_URL"))
	if rpcerr != nil {
		slog.Error("error rpc init")
		os.Exit(-1)
	}

	rpcClient := initStarkRPC(os.Getenv("RPC_URL"), rpcx)
	client, err := initMongoConnections(context.TODO())
	if err != nil {
		slog.Error("error mongo init, ", err)
		os.Exit(-1)
	}
	fmt.Println("mongo connection successful")

	account, err := initAccount(rpcClient)
	if err != nil {
		slog.Error("error account init")
		os.Exit(-1)
	}

	deployer, err := initDeployer(rpcClient, account)
	if err != nil {
		slog.Error("error deployer init")
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
