package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strconv"

	pb "github.com/MarketScrapperAPI/ItemAPI/proto/gen"
	"github.com/MarketScrapperAPI/QueueWorker/models"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const ENV_REDIS_HOST = "REDIS_HOST"
const ENV_REDIS_PORT = "REDIS_PORT"
const ENV_REDIS_DB = "REDIS_DB"
const ENV_REDIS_PASSWORD = "REDIS_PASSWORD"

const ENV_ITEM_API_HOST = "ITEM_API_HOST"
const ENV_ITEM_API_PORT = "ITEM_API_PORT"

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	redisHost := os.Getenv(ENV_REDIS_HOST)
	if redisHost == "" {
		panic("missing env var: " + ENV_REDIS_HOST)
	}
	redisPort := os.Getenv(ENV_REDIS_PORT)
	if redisPort == "" {
		panic("missing env var: " + ENV_REDIS_PORT)
	}
	redisDB := os.Getenv(ENV_REDIS_DB)
	if redisDB == "" {
		panic("missing env var: " + ENV_REDIS_DB)
	}

	dBNumber, err := strconv.Atoi(redisDB)
	if err != nil {
		panic("invalid Redis DB number: " + redisDB)
	}

	itemAPIHost := os.Getenv(ENV_ITEM_API_HOST)
	if itemAPIHost == "" {
		panic("missing env var: " + ENV_ITEM_API_HOST)
	}
	itemAPIPort := os.Getenv(ENV_ITEM_API_PORT)
	if itemAPIPort == "" {
		panic("missing env var: " + ENV_ITEM_API_PORT)
	}

	// Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "",
		DB:       dBNumber,
	})

	ctx := context.Background()

	subscriber := rdb.Subscribe(ctx, "items")

	conn, err := grpc.Dial(itemAPIHost+":"+itemAPIPort, grpc.WithTransportCredentials(insecure.NewCredentials()))

	client := pb.NewItemApiClient(conn)

	for {
		msgStr, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		msg := models.Message{}

		err = json.Unmarshal([]byte(msgStr.Payload), &msg)
		if err != nil {
			panic(err)
		}

		log.Printf("msg: %+v\n", msg)

		req := pb.CreateItemRequest{
			Name:             msg.Item.Name,
			Brand:            msg.Item.Brand,
			Package:          msg.Item.Package,
			PricePerItem:     msg.Item.PricePerItem,
			PricePerQuantity: *msg.Item.PricePerQuantity,
			QuantityUnit:     *msg.Item.QuantityUnit,
			Url:              msg.Item.Url,
			ImageUrl:         msg.Item.ImageUrl,
			MarketName:       msg.Market.Name,
			MarketLocation:   msg.Market.Location,
		}
		log.Printf("msg: %+v\n", req)

		response, err := client.CreateItem(context.Background(), &req)
		if err != nil {
			log.Println(err)
		}
		log.Printf("payload: %+v\n", msgStr.Payload)
		log.Printf("resp: %+v\n", response)
	}
}
