package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
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

	//c := client.NewMarketAPIGRPCClient(itemAPIHost + ":" + itemAPIPort)

	for {
		msg, err := subscriber.ReceiveMessage(ctx)
		if err != nil {
			panic(err)
		}

		//c.Client.CreateItem(ctx, &pb.ItemRequest{})

		fmt.Println(msg.Payload)
	}
}
