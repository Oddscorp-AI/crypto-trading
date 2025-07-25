package exchange

import (
	"sort"
	"sync"
	"time"
)

// OrderType represents a buy or sell order.
// We won't generate stringer; manual implementation.
//
//go:generate stringer -type=OrderType
type OrderType int

const (
	Buy OrderType = iota
	Sell
)

// Order represents a limit order in the order book.
type Order struct {
	ID       int64
	Type     OrderType
	Price    float64
	Quantity float64
	Time     time.Time
}

// Trade represents an executed trade between a buy and sell order.
type Trade struct {
	BuyOrderID  int64
	SellOrderID int64
	Price       float64
	Quantity    float64
	Time        time.Time
}

// OrderBook manages buy and sell limit orders and executed trades.
type OrderBook struct {
	mu     sync.Mutex
	nextID int64
	buys   []Order
	sells  []Order
	trades []Trade
}

// NewBook creates an empty OrderBook.
func NewBook() *OrderBook {
	return &OrderBook{}
}

// PlaceOrder adds an order to the book and attempts to match it.
// Returned trades are the executions generated by this order, and the final
// state of the order (with remaining quantity, if any).
func (ob *OrderBook) PlaceOrder(o Order) ([]Trade, Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()

	ob.nextID++
	o.ID = ob.nextID
	o.Time = time.Now()

	var trades []Trade
	if o.Type == Buy {
		trades, o = ob.matchBuy(o)
		if o.Quantity > 0 {
			ob.insertBuy(o)
		}
	} else {
		trades, o = ob.matchSell(o)
		if o.Quantity > 0 {
			ob.insertSell(o)
		}
	}
	ob.trades = append(ob.trades, trades...)
	return trades, o
}

func (ob *OrderBook) matchBuy(buy Order) ([]Trade, Order) {
	var trades []Trade
	for buy.Quantity > 0 && len(ob.sells) > 0 {
		best := &ob.sells[0]
		if best.Price > buy.Price {
			break
		}
		qty := min(buy.Quantity, best.Quantity)
		trade := Trade{
			BuyOrderID:  buy.ID,
			SellOrderID: best.ID,
			Price:       best.Price,
			Quantity:    qty,
			Time:        time.Now(),
		}
		trades = append(trades, trade)
		buy.Quantity -= qty
		best.Quantity -= qty
		if best.Quantity == 0 {
			ob.sells = ob.sells[1:]
		}
	}
	return trades, buy
}

func (ob *OrderBook) matchSell(sell Order) ([]Trade, Order) {
	var trades []Trade
	for sell.Quantity > 0 && len(ob.buys) > 0 {
		best := &ob.buys[0]
		if best.Price < sell.Price {
			break
		}
		qty := min(sell.Quantity, best.Quantity)
		trade := Trade{
			BuyOrderID:  best.ID,
			SellOrderID: sell.ID,
			Price:       best.Price,
			Quantity:    qty,
			Time:        time.Now(),
		}
		trades = append(trades, trade)
		sell.Quantity -= qty
		best.Quantity -= qty
		if best.Quantity == 0 {
			ob.buys = ob.buys[1:]
		}
	}
	return trades, sell
}

func (ob *OrderBook) insertBuy(o Order) {
	i := sort.Search(len(ob.buys), func(i int) bool {
		if ob.buys[i].Price == o.Price {
			return ob.buys[i].Time.After(o.Time)
		}
		return ob.buys[i].Price < o.Price
	})
	ob.buys = append(ob.buys, Order{})
	copy(ob.buys[i+1:], ob.buys[i:])
	ob.buys[i] = o
}

func (ob *OrderBook) insertSell(o Order) {
	i := sort.Search(len(ob.sells), func(i int) bool {
		if ob.sells[i].Price == o.Price {
			return ob.sells[i].Time.After(o.Time)
		}
		return ob.sells[i].Price > o.Price
	})
	ob.sells = append(ob.sells, Order{})
	copy(ob.sells[i+1:], ob.sells[i:])
	ob.sells[i] = o
}

// Book returns a snapshot of current buy and sell orders.
func (ob *OrderBook) Book() (buys, sells []Order) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	buys = append([]Order(nil), ob.buys...)
	sells = append([]Order(nil), ob.sells...)
	return
}

// Trades returns all executed trades.
func (ob *OrderBook) Trades() []Trade {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return append([]Trade(nil), ob.trades...)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
