package main

import (
	"context"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"tickets/message"

	"github.com/ThreeDotsLabs/go-event-driven/v2/common/clients"
	"github.com/ThreeDotsLabs/go-event-driven/v2/common/log"

	"tickets/adapters"
	"tickets/service"
)

func main() {
	log.Init(slog.LevelInfo)

	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("Error loading .env file")
	}

	apiClients, err := clients.NewClients(os.Getenv("GATEWAY_ADDR"), nil)
	if err != nil {
		panic(err)
	}

	redisClient := message.NewRedisClient(os.Getenv("REDIS_ADDR"))
	defer redisClient.Close()

	spreadsheetsService := adapters.NewSpreadsheetsAPIClient(apiClients)
	receiptsService := adapters.NewReceiptsServiceClient(apiClients)

	//go message.NewHandler(receiptsService, spreadsheetsService)

	err = service.New(redisClient, spreadsheetsService, receiptsService).Run(context.Background())
	if err != nil {
		panic(err)
	}
}
