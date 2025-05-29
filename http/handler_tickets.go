package http

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tickets/worker"
)

type ticketsConfirmationRequest struct {
	Tickets []string `json:"tickets"`
}

func (h Handler) PostTicketsConfirmation(c echo.Context) error {
	var request ticketsConfirmationRequest
	err := c.Bind(&request)
	if err != nil {
		return err
	}

	for _, ticket := range request.Tickets {
		h.worker.Send(
			worker.Message{
				Task:     0,
				TicketID: ticket,
			},
			worker.Message{
				Task:     1,
				TicketID: ticket,
			},
		)

	}

	return c.NoContent(http.StatusOK)
}
