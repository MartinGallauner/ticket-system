package service

import (
	"context"
	"errors"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"log/slog"
	stdHTTP "net/http"
	ticketsHttp "tickets/http"
	"tickets/message"
)

type Service struct {
	echoRouter *echo.Echo
}

func New(redisClient *redis.Client, spreadsheetservcie message.SpreadsheetsAPI, receiptService message.ReceiptsService) Service {
	logger := watermill.NewSlogLogger(slog.Default())
	publisher := message.NewPublisher(redisClient, logger)
	echoRouter := ticketsHttp.NewHttpRouter(publisher)

	message.NewHandler(receiptService, spreadsheetservcie, redisClient, logger)

	return Service{
		echoRouter: echoRouter,
	}
}

func (s Service) Run(ctx context.Context) error {

	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}
	return nil
}
