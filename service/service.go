package service

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	stdHTTP "net/http"
	ticketsHttp "tickets/http"
	"tickets/worker"
)

type Service struct {
	echoRouter *echo.Echo
	worker     *worker.Worker
}

func New(
	spreadsheetsService worker.SpreadsheetsAPI,
	receiptsService worker.ReceiptsService,
) Service {
	w := worker.NewWorker(spreadsheetsService, receiptsService)
	echoRouter := ticketsHttp.NewHttpRouter(w)
	return Service{
		echoRouter: echoRouter,
		worker:     w,
	}
}

func (s Service) Run(ctx context.Context) error {
	go s.worker.Run()
	err := s.echoRouter.Start(":8080")
	if err != nil && !errors.Is(err, stdHTTP.ErrServerClosed) {
		return err
	}
	return nil
}
