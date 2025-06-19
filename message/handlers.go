package message

import (
	"context"

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
			return receiptsClient.IssueReceipt(msg.Context(), string(msg.Payload))
	})

	router.AddNoPublisherHandler(
		"append-to-tracker-handler",
		"append-to-tracker",
		spreadsheetSub,
		func(msg *message.Message) error {
			return spreadsheetClient.AppendRow(msg.Context(), "tickets-to-print", []string{string(msg.Payload)})
		})

	return router
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}
