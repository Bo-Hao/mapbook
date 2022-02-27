package mapbook

import (
	"fmt"
	"log"
	"sort"

	"github.com/shopspring/decimal"
)

// [price, volume]

// Ask
func NewAskBook(IsEvent bool) *AskBook {
	return &AskBook{
		Book: map[string]string{},

		isEvent: IsEvent,
		event:   [][]string{},
	}
}

func (book *AskBook) Snapshot(snapshot [][]string) {
	book.Lock()
	defer book.Unlock()

	for i := range snapshot {
		order := snapshot[i]
		price := order[0]
		volume := order[1]

		book.Book[price] = volume
	}
}

func (book *AskBook) Update(update [][]string) {
	book.Lock()
	defer book.Unlock()

	for i := range update {
		order := update[i]
		price := order[0]
		volume := order[1]

		switch {
		case volume == "0":
			delete(book.Book, price)

		case volume != "0":
			if _, ok := book.Book[price]; ok {
				book.Book[price] = volume
			} else {
				book.Book[price] = volume
			}
		}
	}
}

func (book *AskBook) GetAll() ([][]string, bool) {
	book.RLock()
	Book := book.Book
	book.RUnlock()

	prices := make([]decimal.Decimal, 0, len(book.Book))
	for k := range Book {
		p, err := decimal.NewFromString(k)
		if err != nil {
			log.Println(err)
		}
		prices = append(prices, p)
	}

	// from small to big
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].LessThan(prices[j])
	})

	asks := make([][]string, 0, len(prices))
	for i := range prices {
		p := prices[i].String()
		v := Book[p]
		asks = append(asks, []string{p, v})
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
		Book: map[string]string{},

		isEvent: IsEvent,
		event:   [][]string{},
	}
}

func (book *BidBook) Snapshot(snapshot [][]string) {
	book.Lock()
	defer book.Unlock()

	for i := range snapshot {
		order := snapshot[i]
		price := order[0]
		volume := order[1]

		book.Book[price] = volume
	}
}

func (book *BidBook) Update(update [][]string) {
	book.Lock()
	defer book.Unlock()

	for i := range update {
		order := update[i]
		price := order[0]
		volume := order[1]

		switch {
		case volume == "0":
			delete(book.Book, price)

		case volume != "0":
			if _, ok := book.Book[price]; ok {
				book.Book[price] = volume
			} else {
				book.Book[price] = volume
			}
		}
	}
}

func (book *BidBook) GetAll() ([][]string, bool) {
	book.RLock()
	Book := book.Book
	book.RUnlock()

	prices := make([]decimal.Decimal, 0, len(book.Book))
	for k := range Book {
		p, err := decimal.NewFromString(k)
		if err != nil {
			log.Println(err)
		}
		prices = append(prices, p)
	}

	// from small to big
	sort.Slice(prices, func(i, j int) bool {
		return prices[i].GreaterThan(prices[j])
	})

	bids := make([][]string, 0, len(prices))
	for i := range prices {
		p := prices[i].String()
		v := Book[p]
		bids = append(bids, []string{p, v})
	}

	if len(bids) != 0 {
		return bids, true
	} else {
		fmt.Println(book)
		return bids, false
	}
}
