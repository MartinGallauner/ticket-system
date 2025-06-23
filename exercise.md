# Project: Append more details to the spreadsheet

As mentioned in the previous exercise, the Operations team has asked us to add the customer's email and the price to the spreadsheet.
The spreadsheet format they want looks like this:

|          |                   |        |     |
|----------|-------------------|--------|-----|
| ticket-1 | user@example.com  | 50.00  | EUR |
| ticket-2 | user2@example.com | 100.00 | USD |


First, we need these details in the messages we publish to the Pub/Sub topics.
We still want to publish two messages per ticket: one for issuing a receipt and one for appending it to the spreadsheet.

Let's update the message sent to the `append-to-tracker` topic to include the customer's email and the price.

The data sent to the `append-to-tracker` topic should have the following format:

```json
{
  "ticket_id": "ticket-1",
  "customer_email": "user@example.com",
  "price": {
    "amount": "50.00",
    "currency": "EUR"
  }
}
```

The Go struct for the payload can look like this:

```go
type AppendToTrackerPayload struct {
    TicketID      string `json:"ticket_id"`
    CustomerEmail string `json:"customer_email"`
    Price         Money  `json:"price"`
}

type Money struct {
    Amount   string `json:"amount"`
    Currency string `json:"currency"`
}
```

**You can keep this struct in the `entities` package.**
You will use it in both the HTTP handler and the message handler.

{{tip}}

For simplicity, in the example solution we reuse the `Money` entity for both HTTP requests and message payloads.

But it's not always the best idea.
It can lead to tight coupling between different parts of the system.

You can read more about it in [When to avoid DRY in Go](https://threedots.tech/post/things-to-know-about-dry/).

{{endtip}}

{{.Exercise}}

1. **Extend the message published to the `append-to-tracker` topic to include the customer's email and the price.**

In the HTTP handler, marshal the struct to JSON and publish it as the message payload.

2. **Update the message handler so that it unmarshalls the payload into the struct.**

Add extra fields to the spreadsheet:

```diff
 	return spreadsheetsAPI.AppendRow(
		ctx,
 		"tickets-to-print",
-		[]string{string(msg.Payload)},
+		[]string{payload.TicketID, payload.CustomerEmail, payload.Price.Amount, payload.Price.Currency},
 	)
 },
```

For now, keep the message sent to `issue-receipt` as-is.
We'll update it in the next exercise.
