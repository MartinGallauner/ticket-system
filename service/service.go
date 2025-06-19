package service

import (
	"context"
	"errors"
	"log/slog"
	stdHTTP "net/http"
	"os"
	"os/signal"
	ticketsHttp "tickets/http"
	"tickets/message"

	"github.com/ThreeDotsLabs/watermill"
	m2 "github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
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

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error { //start watermill router
		return s.router.Run(ctx)
	})

	group.Go(func() error { //start http server
		err := s.echoRouter.Start(":8080")
		if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
			return err
		}
		return nil
	})

	group.Go(func() error {
		<- ctx.Done()
		return s.echoRouter.Shutdown(ctx)
	})

	err := group.Wait()
	if err != nil {
		return err
	}
	return nil

}
