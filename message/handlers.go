package message

import (
	"context"
	"encoding/json"
	"tickets/entities"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
)

func NewRouter(
	receiptsClient ReceiptsService,
	spreadsheetClient SpreadsheetsAPI,
	rdb *redis.Client,
	logger watermill.LoggerAdapter,
) *message.Router {

	router := message.NewDefaultRouter(logger)

	receiptsSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "issue-receipt",
	}, logger)
	if err != nil {
		panic(err)
	}

	spreadsheetSub, err := redisstream.NewSubscriber(redisstream.SubscriberConfig{
		Client:        rdb,
		ConsumerGroup: "append-to-tracker",
	}, logger)
	if err != nil {
		panic(err)
	}

	router.AddNoPublisherHandler(
		"issue-receipt-handler",
		"issue-receipt",
		receiptsSub,
		func(msg *message.Message) error {
			var payload entities.IssueReceiptPayload
			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				return err
			}
			return receiptsClient.IssueReceipt(msg.Context(), payload)
		})

	router.AddNoPublisherHandler(
		"append-to-tracker-handler",
		"append-to-tracker",
		spreadsheetSub,
		func(msg *message.Message) error {

			var payload entities.AppendToTrackerPayload
			err := json.Unmarshal(msg.Payload, &payload)
			if err != nil {
				return err
			}

			return spreadsheetClient.AppendRow(
				msg.Context(),
				"tickets-to-print",
				[]string{string(payload.TicketID), payload.CustomerEmail, payload.Price.Amount, payload.Price.Currency},
			)
		})

	return router
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, request entities.IssueReceiptPayload) error
}
