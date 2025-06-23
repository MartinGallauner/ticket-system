# Project: Send price to the Receipts service

{{background}}

Since we are already appending the ticket price to the spreadsheet, 
the Finance team has asked us to also use it when issuing receipts.

How did they know the price before? Good question...
But it makes sense to send it since we have it now.

{{endbackground}}

We need to extend the message sent to the `issue-receipt` topic to include the price.
It should look like this:

```json
{
    "ticket_id": "ticket-1",
    "price": {
      "amount": "50.00",
      "currency": "EUR"
  }
}
```

You can use a Go struct like this:

```go
type IssueReceiptPayload struct {
    TicketID string `json:"ticket_id"`
    Price    Money  `json:"price"`
}
```

{{.Exercise}}

1. Update the HTTP handler, to publish JSON payload on the `issue-receipt` topic in the format above.

2. Update the message handler so that it unmarshals the payload on the struct.

3. Modify the receipts client so that it can accept a price. It might look like this:

```go
func (c ReceiptsClient) IssueReceipt(ctx context.Context, request entities.IssueReceiptRequest) error {
```

It's a good practice to not import the `adapters` package in other packages.
You can keep the `IssueReceiptRequest` struct in `entities`:

```go
type IssueReceiptRequest struct {
    TicketID string
    Price    Money
}
```

Next, include the price in the `Receipts.PutReceiptsWithResponse`:

```diff
resp, err := c.clients.Receipts.PutReceiptsWithResponse(ctx, receipts.CreateReceipt{
-	TicketId: ticketID,
+	TicketId: request.TicketID,
+	Price: receipts.Money{
+		MoneyAmount:   request.Price.Amount,
+		MoneyCurrency: request.Price.Currency,
+	},
})
```

{{tip}}

Note that there are two `Money` types: one comes from the HTTP request, and the other from the `receipts` package.
This is correct. It's a common situation when integrating with external systems that the data formats are slightly different.

{{endtip}}
