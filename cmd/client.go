package cmd

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	rpcclient "github.com/cometbft/cometbft/rpc/client"
	rpchttp "github.com/cometbft/cometbft/rpc/client/http"
	libclient "github.com/cometbft/cometbft/rpc/jsonrpc/client"
	cosmosClient "github.com/cosmos/cosmos-sdk/client"
	tmservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	"google.golang.org/grpc"
)

const (
	sentryGRPCTimeoutSeconds = 5
	RPCTimeoutSeconds        = 5
)

func newClient(addr string) (rpcclient.Client, error) {
	httpClient, err := libclient.DefaultHTTPClient(addr)
	if err != nil {
		return nil, err
	}

	httpClient.Timeout = 10 * time.Second
	rpcClient, err := rpchttp.NewWithClient(addr, "/websocket", httpClient)
	if err != nil {
		return nil, err
	}

	return rpcClient, nil
}

func getCosmosClient(rpcAddress string, chainID string) (*cosmosClient.Context, error) {
	client, err := newClient(rpcAddress)
	if err != nil {
		return nil, err
	}
	return &cosmosClient.Context{
		Client:       client,
		ChainID:      chainID,
		Input:        os.Stdin,
		Output:       os.Stdout,
		OutputFormat: "json",
	}, nil
}

func getSlashingInfo(client *cosmosClient.Context) (*slashingtypes.QueryParamsResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*RPCTimeoutSeconds))
	defer cancel()
	return slashingtypes.NewQueryClient(client).Params(ctx, &slashingtypes.QueryParamsRequest{})
}

func getSigningInfo(client *cosmosClient.Context, address string) (*slashingtypes.QuerySigningInfoResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*RPCTimeoutSeconds))
	defer cancel()
	return slashingtypes.NewQueryClient(client).SigningInfo(ctx, &slashingtypes.QuerySigningInfoRequest{
		ConsAddress: address,
	})
}

func getSentryInfo(grpcAddr string) (*tmservice.GetNodeInfoResponse, *tmservice.GetLatestBlockResponse, error) {
	creds := grpc.WithTransportCredentials(insecure.NewCredentials())
	if strings.Contains(grpcAddr, ":443") {
		creds = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{MinVersion: tls.VersionTLS12}))
	}
	// ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	conn, err := grpc.NewClient(
		grpcAddr,
		creds)

	if err != nil {
		return nil, nil, err
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Could not close grpc connection: %v", err)
		}
	}(conn)
	serviceClient := tmservice.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*sentryGRPCTimeoutSeconds))
	defer cancel()
	nodeInfo, err := serviceClient.GetNodeInfo(ctx, &tmservice.GetNodeInfoRequest{})
	if err != nil {
		return nil, nil, err
	}
	syncingInfo, err := serviceClient.GetLatestBlock(ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, nil, err
	}
	return nodeInfo, syncingInfo, nil
}

func getWalletBalance(client *cosmosClient.Context, address string, denom string) (*banktypes.QueryBalanceResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*RPCTimeoutSeconds))
	defer cancel()
	return banktypes.NewQueryClient(client).Balance(ctx, &banktypes.QueryBalanceRequest{
		Address: address,
		Denom:   denom,
	})
}
