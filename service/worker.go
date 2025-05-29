package service

import (
	ticketsHttp "tickets/http"
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
	spreadsheetsAPI ticketsHttp.SpreadsheetsAPI
	receiptsService ticketsHttp.ReceiptsService
}

func NewWorker(api ticketsHttp.SpreadsheetsAPI, service ticketsHttp.ReceiptsService) *Worker {
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
			w.receiptsService.IssueReceipt()
		case TaskAppendToTracker:
		}
	}
}
