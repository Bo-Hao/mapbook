package mapbook

import (
	"sort"
	"sync"

	"github.com/shopspring/decimal"
)

// [price, volume]

// Ask
func NewAskBook(IsEvent bool) *AskBook {
	return &AskBook{
		Book: sync.Map{},

		isEvent: IsEvent,
		event:   [][]string{},
	}
}

func (book *AskBook) Snapshot(snapshot [][]string) {
	for i := range snapshot {
		order := snapshot[i]
		price := order[0]
		volume := order[1]

		book.Book.Store(price, volume)
	}
}

func (book *AskBook) Update(update [][]string) {
	for i := range update {
		order := update[i]
		price := order[0]
		volume := order[1]

		switch {
		case volume == "0":
			book.Book.Delete(price)

		case volume != "0":
			book.Book.Store(price, volume)
		}
	}
}

func (book *AskBook) GetAll() ([][]string, bool) {
	var asks [][]string
	// iterate the map
	book.Book.Range(func(k, v interface{}) bool {
		asks = append(asks, []string{k.(string), v.(string)})
		return true
	})

	// from small to big
	sort.Slice(asks, func(i, j int) bool {
		pi, _ := decimal.NewFromString(asks[i][0])
		pj, _ := decimal.NewFromString(asks[j][0])
		return pi.LessThan(pj)
	})

	if len(asks) != 0 {
		return asks, true
	} else {
		return asks, false
	}
}

// bid
func NewBidBook(IsEvent bool) *BidBook {
	return &BidBook{
		Book: sync.Map{},

		isEvent: IsEvent,
		event:   [][]string{},
	}
}

func (book *BidBook) Snapshot(snapshot [][]string) {
	for i := range snapshot {
		order := snapshot[i]
		price := order[0]
		volume := order[1]

		book.Book.Store(price, volume)
	}
}

func (book *BidBook) Update(update [][]string) {
	for i := range update {
		order := update[i]
		price := order[0]
		volume := order[1]

		switch {
		case volume == "0":
			book.Book.Delete(price)

		case volume != "0":
			book.Book.Store(price, volume)
		}
	}
}

func (book *BidBook) GetAll() ([][]string, bool) {
	var bids [][]string
	// iterate the map
	book.Book.Range(func(k, v interface{}) bool {
		bids = append(bids, []string{k.(string), v.(string)})
		return true
	})

	// from small to big
	sort.Slice(bids, func(i, j int) bool {
		pi, _ := decimal.NewFromString(bids[i][0])
		pj, _ := decimal.NewFromString(bids[j][0])
		return pi.GreaterThan(pj)
	})

	if len(bids) != 0 {
		return bids, true
	} else {
		return bids, false
	}
}
