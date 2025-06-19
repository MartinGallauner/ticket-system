package service

import (
	"context"
	"errors"
	"log/slog"
	stdHTTP "net/http"
	ticketsHttp "tickets/http"
	"tickets/message"

	"github.com/ThreeDotsLabs/watermill"
	m2 "github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	echoRouter *echo.Echo
	router     *m2.Router
}

func New(redisClient *redis.Client, spreadsheetservcie message.SpreadsheetsAPI, receiptService message.ReceiptsService) Service {
	logger := watermill.NewSlogLogger(slog.Default())
	publisher := message.NewPublisher(redisClient, logger)
	echoRouter := ticketsHttp.NewHttpRouter(publisher)

	r := message.NewRouter(receiptService, spreadsheetservcie, redisClient, logger)

	return Service{
		echoRouter: echoRouter,
		router:     r,
	}
}

func (s Service) Run(ctx context.Context) error {

	go func() {
		if err := s.router.Run(ctx); err != nil {
			panic(err)
		}
	}()

	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}
	return nil
}
