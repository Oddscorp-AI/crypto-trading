package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"crypto-trading/internal/exchange"
	"crypto-trading/internal/wallet"
)

var (
	w  = wallet.New(1000)
	ob = exchange.NewBook()
)

func main() {
	http.HandleFunc("/records", addRecordHandler)
	http.HandleFunc("/history", historyHandler)
	http.HandleFunc("/orders", orderHandler)
	http.HandleFunc("/orderbook", bookHandler)
	http.HandleFunc("/trades", tradesHandler)

	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func addRecordHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Datetime string  `json:"datetime"`
		Amount   float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	t, err := time.Parse(time.RFC3339, req.Datetime)
	if err != nil {
		http.Error(rw, "invalid datetime", http.StatusBadRequest)
		return
	}
	w.AddRecord(wallet.Record{Time: t, Amount: req.Amount})
	rw.WriteHeader(http.StatusOK)
}

func historyHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		StartDatetime string `json:"startDatetime"`
		EndDatetime   string `json:"endDatetime"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	start, err := time.Parse(time.RFC3339, req.StartDatetime)
	if err != nil {
		http.Error(rw, "invalid startDatetime", http.StatusBadRequest)
		return
	}
	end, err := time.Parse(time.RFC3339, req.EndDatetime)
	if err != nil {
		http.Error(rw, "invalid endDatetime", http.StatusBadRequest)
		return
	}
	history := w.History(start, end)
	var resp []struct {
		Datetime string  `json:"datetime"`
		Amount   float64 `json:"amount"`
	}
	for _, h := range history {
		resp = append(resp, struct {
			Datetime string  `json:"datetime"`
			Amount   float64 `json:"amount"`
		}{Datetime: h.Time.Format(time.RFC3339), Amount: h.Amount})
	}
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(resp)
}

func orderHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req struct {
		Type     string  `json:"type"`
		Price    float64 `json:"price"`
		Quantity float64 `json:"quantity"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	var ot exchange.OrderType
	switch req.Type {
	case "buy":
		ot = exchange.Buy
	case "sell":
		ot = exchange.Sell
	default:
		http.Error(rw, "invalid type", http.StatusBadRequest)
		return
	}
	trades, order := ob.PlaceOrder(exchange.Order{Type: ot, Price: req.Price, Quantity: req.Quantity})
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(struct {
		Order  exchange.Order   `json:"order"`
		Trades []exchange.Trade `json:"trades"`
	}{Order: order, Trades: trades})
}

func bookHandler(rw http.ResponseWriter, r *http.Request) {
	buys, sells := ob.Book()
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(struct {
		Buys  []exchange.Order `json:"buys"`
		Sells []exchange.Order `json:"sells"`
	}{Buys: buys, Sells: sells})
}

func tradesHandler(rw http.ResponseWriter, r *http.Request) {
	trades := ob.Trades()
	rw.Header().Set("Content-Type", "application/json")
	json.NewEncoder(rw).Encode(trades)
}
