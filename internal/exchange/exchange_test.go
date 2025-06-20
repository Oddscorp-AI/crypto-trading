package exchange

import "testing"

func TestMatch(t *testing.T) {
	ob := NewBook()
	// place sell order first
	trades, sell := ob.PlaceOrder(Order{Type: Sell, Price: 100, Quantity: 1})
	if len(trades) != 0 || sell.ID == 0 {
		t.Fatalf("unexpected trades or id")
	}

	// place buy order that matches
	trades, buy := ob.PlaceOrder(Order{Type: Buy, Price: 120, Quantity: 1})
	if buy.Quantity != 0 {
		t.Fatalf("buy order should be fully filled")
	}
	if len(trades) != 1 {
		t.Fatalf("expected 1 trade, got %d", len(trades))
	}
	if trades[0].Price != 100 {
		t.Fatalf("trade price %v", trades[0].Price)
	}
}

func TestCancelOrder(t *testing.T) {
	ob := NewBook()
	_, o := ob.PlaceOrder(Order{Type: Buy, Price: 100, Quantity: 1})
	if o.ID == 0 {
		t.Fatal("missing id")
	}
	if !ob.CancelOrder(o.ID) {
		t.Fatal("cancel failed")
	}
	buys, sells := ob.Book()
	if len(buys) != 0 || len(sells) != 0 {
		t.Fatal("order not removed")
	}
	if ob.CancelOrder(o.ID) {
		t.Fatal("should not cancel twice")
	}
}
