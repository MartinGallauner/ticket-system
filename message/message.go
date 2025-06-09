package message

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-redisstream/pkg/redisstream"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewPublisher(rdb *redis.Client, logger watermill.LoggerAdapter) message.Publisher {
	var pub message.Publisher
	pub, err := redisstream.NewPublisher(redisstream.PublisherConfig{
		Client: rdb,
	}, logger)
	if err != nil {
		log.Fatal(err)
	}
	return pub
}

func processReceipt(subscriber message.Subscriber, action func(ctx context.Context, orderID string) error) {
	messages, err := subscriber.Subscribe(context.Background(), "issue-receipt")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		orderID := string(msg.Payload)

		err := action(context.Background(), orderID)
		if err != nil {
			msg.Nack()
		} else {
			msg.Ack()
		}
	}
}

func processSpreadsheet(subscriber message.Subscriber, action func(ctx context.Context, sheetName string, row []string) error) {
	messages, err := subscriber.Subscribe(context.Background(), "append-to-tracker")
	if err != nil {
		panic(err)
	}

	for msg := range messages {
		orderID := string(msg.Payload)

		err := action(context.Background(), "tickets-to-print", []string{orderID})
		if err != nil {
			msg.Nack()
		} else {
			msg.Ack()
		}
	}
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: addr,
	})
}
