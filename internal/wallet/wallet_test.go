package wallet

import (
	"reflect"
	"testing"
	"time"
)

func TestHistory(t *testing.T) {
	w := New(1000)
	w.AddRecord(Record{Time: mustParse("2019-10-05T14:45:05+07:00"), Amount: 10})

	start := mustParse("2019-10-05T10:00:00+07:00")
	end := mustParse("2019-10-05T15:00:00+07:00")

	history := w.History(start, end)

	expected := []Record{
		{Time: mustParse("2019-10-05T10:00:00+07:00"), Amount: 1000},
		{Time: mustParse("2019-10-05T11:00:00+07:00"), Amount: 1000},
		{Time: mustParse("2019-10-05T12:00:00+07:00"), Amount: 1000},
		{Time: mustParse("2019-10-05T13:00:00+07:00"), Amount: 1000},
		{Time: mustParse("2019-10-05T14:00:00+07:00"), Amount: 1010},
		{Time: mustParse("2019-10-05T15:00:00+07:00"), Amount: 1010},
	}

	if !reflect.DeepEqual(history, expected) {
		t.Fatalf("expected %v, got %v", expected, history)
	}
}

func mustParse(v string) time.Time {
	t, err := time.Parse(time.RFC3339, v)
	if err != nil {
		panic(err)
	}
	return t
}
