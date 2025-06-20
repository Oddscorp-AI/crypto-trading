# Crypto Trading Wallet API

This project provides a simple REST API to track incoming Bitcoin deposits,
place buy/sell orders and retrieve the wallet balance history at hourly
intervals.

## Building and Running

1. **Install Go 1.20+**
2. Fetch the dependencies and build:

   ```bash
   go build ./cmd/server
   ```

3. Run the server:

   ```bash
   go run ./cmd/server
   ```

The server listens on `:8080` by default.

## API

### POST /records

Add a new deposit record.

Request body:

```json
{
  "datetime": "2019-10-05T14:45:05+07:00",
  "amount": 10
}
```

### POST /history

Return the wallet balance at the end of each hour between two datetimes.

Request body:

```json
{
  "startDatetime": "2019-10-05T10:48:01+00:00",
  "endDatetime": "2019-10-05T18:48:02+00:00"
}
```

Example response:

```json
[
  { "datetime": "2019-10-05T10:00:00+00:00", "amount": 1000 },
  { "datetime": "2019-10-05T11:00:00+00:00", "amount": 1000 }
]
```

## Testing

Run unit tests with:

```bash
go test ./...
```

### POST /orders

Place a buy or sell limit order.

Request body:

```json
{
  "type": "buy",
  "price": 100,
  "quantity": 1
}
```

### GET /orderbook

Return current buy and sell orders.

### GET /trades

Return all executed trades.
