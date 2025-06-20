package wallet

import (
	"sort"
	"sync"
	"time"
)

// Record represents a deposit record with timestamp and amount.
type Record struct {
	Time   time.Time
	Amount float64
}

// Wallet manages deposit records and base balance.
type Wallet struct {
	mu         sync.RWMutex
	baseAmount float64
	records    []Record
}

// New returns a new Wallet with the given base amount.
func New(base float64) *Wallet {
	return &Wallet{baseAmount: base}
}

// AddRecord adds a deposit record.
func (w *Wallet) AddRecord(r Record) {
	w.mu.Lock()
	defer w.mu.Unlock()
	// insert while keeping slice sorted by time
	idx := sort.Search(len(w.records), func(i int) bool {
		return w.records[i].Time.After(r.Time) || w.records[i].Time.Equal(r.Time)
	})
	w.records = append(w.records, Record{})
	copy(w.records[idx+1:], w.records[idx:])
	w.records[idx] = r
}

// BalanceAt returns the wallet balance up to and including the specified time.
func (w *Wallet) BalanceAt(t time.Time) float64 {
	w.mu.RLock()
	defer w.mu.RUnlock()
	sum := w.baseAmount
	for _, r := range w.records {
		if !r.Time.After(t) {
			sum += r.Amount
		} else {
			break
		}
	}
	return sum
}

// History returns the wallet balance at the end of each hour between start and end.
func (w *Wallet) History(start, end time.Time) []Record {
	if end.Before(start) {
		start, end = end, start
	}
	// align start and end to hour
	start = start.Truncate(time.Hour)
	end = end.Truncate(time.Hour)

	var result []Record
	for t := start; !t.After(end); t = t.Add(time.Hour) {
		balance := w.BalanceAt(t.Add(time.Hour - time.Nanosecond))
		result = append(result, Record{Time: t, Amount: balance})
	}
	return result
}
