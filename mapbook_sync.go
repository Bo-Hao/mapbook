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
	var price []decimal.Decimal
	// iterate the map
	book.Book.Range(func(k, v interface{}) bool {
		p, _ := decimal.NewFromString(k.(string))
		price = append(price, p)
		return true
	})

	// from small to big
	sort.Slice(price, func(i, j int) bool {
		return price[i].LessThan(price[j])
	})

	var asks [][]string
	for i := range price {
		k := price[i].String()
		v, ok := book.Book.Load(k)
		if ok {
			asks = append(asks, []string{k, v.(string)})
		}
	}

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
	var price []decimal.Decimal
	// iterate the map
	book.Book.Range(func(k, v interface{}) bool {
		p, _ := decimal.NewFromString(k.(string))
		price = append(price, p)
		return true
	})

	// from small to big
	sort.Slice(price, func(i, j int) bool {
		return price[i].GreaterThan(price[j])
	})

	var bids [][]string
	for i := range price {
		k := price[i].String()
		v, ok := book.Book.Load(k)
		if ok {
			bids = append(bids, []string{k, v.(string)})
		}
	}

	if len(bids) != 0 {
		return bids, true
	} else {
		return bids, false
	}
}
