package http

import (
	"net/http"

	libHttp "github.com/ThreeDotsLabs/go-event-driven/v2/common/http"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/labstack/echo/v4"
)

func NewHttpRouter(p message.Publisher) *echo.Echo {
	e := libHttp.NewEcho()

	handler := Handler{
		p,
	}

	e.POST("/tickets-status", handler.PostTicketsStatus)
	e.GET("/health", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	return e
}
