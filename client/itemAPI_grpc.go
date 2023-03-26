package client

import (
	"log"

	pb "github.com/MarketScrapperAPI/ItemAPI/proto/gen"
	"google.golang.org/grpc"
)

type MarketAPIGRPCClient struct {
	Client *pb.ItemApiClient
}

func NewMarketAPIGRPCClient(apiUrl string) MarketAPIGRPCClient {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(apiUrl+":9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := pb.NewItemApiClient(conn)

	return MarketAPIGRPCClient{
		Client: &client,
	}
}
