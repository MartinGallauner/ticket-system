package message

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/redis/go-redis/v9"
)

func NewHandler(
	receiptsClient ReceiptsService,
	spreadsheetsClient SpreadsheetsAPI,
	rdb *redis.Client,
	logger watermill.LoggerAdapter,
) {
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

	go processReceipt(receiptsSub, receiptsClient.IssueReceipt)
	go processSpreadsheet(spreadsheetSub, spreadsheetsClient.AppendRow)
}

type SpreadsheetsAPI interface {
	AppendRow(ctx context.Context, sheetName string, row []string) error
}

type ReceiptsService interface {
	IssueReceipt(ctx context.Context, ticketID string) error
}
