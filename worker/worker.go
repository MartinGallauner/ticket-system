package worker

import (
	"context"
)

type Task int

const (
	TaskIssueReceipt Task = iota
	TaskAppendToTracker
)

type Message struct {
	Task     Task
	TicketID string
}

type Worker struct {
	queue           chan Message
	spreadsheetsAPI SpreadsheetsAPI
	receiptsService ReceiptsService
}

func NewWorker(api SpreadsheetsAPI, service ReceiptsService) *Worker {
	return &Worker{
		spreadsheetsAPI: api,
		receiptsService: service,
		queue:           make(chan Message, 100)}
}

func (w *Worker) Send(msg ...Message) {
	for _, m := range msg {
		w.queue <- m
	}
}

func (w *Worker) Run() {
	for msg := range w.queue {
		switch msg.Task {
		case TaskIssueReceipt:
			err := w.receiptsService.IssueReceipt(context.Background(), msg.TicketID)
			if err != nil {
				w.Send(msg)
			}
		case TaskAppendToTracker:
			err := w.spreadsheetsAPI.AppendRow(context.Background(), "tickets-to-print", []string{msg.TicketID})
			if err != nil {
				w.Send(msg)
			}
		}
	}
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}
