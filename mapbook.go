package mapbook

import (
	"sort"

	"github.com/shopspring/decimal"
)

// [price, volume]

// Ask
func NewAskBook(IsEvent bool) *AskBook {
	return &AskBook{
		OrderKey: []string{},
		Book:     map[string]string{},

		isEvent: IsEvent,
		event:   [][]string{},
	}
}

func (book *AskBook) Snapshot(snapshot [][]string) {
	orderkey := make([]string, len(snapshot))
	book.Lock()
	defer book.Unlock()

	// from small to big
	sort.Slice(snapshot, func(i, j int) bool {
		oi, _ := decimal.NewFromString(snapshot[i][0])
		oj, _ := decimal.NewFromString(snapshot[j][0])
		return oi.LessThan(oj)
	})

	for i := range snapshot {
		order := snapshot[i]
		price := order[0]
		volume := order[1]

		book.Book[price] = volume
		orderkey[i] = price
	}

	book.OrderKey = orderkey
}

func (book *AskBook) Update(update [][]string) {

	sort.Slice(update, func(i, j int) bool {
		oi, _ := decimal.NewFromString(update[i][0])
		oj, _ := decimal.NewFromString(update[j][0])
		return oi.LessThan(oj)
	})

	book.Lock()
	defer book.Unlock()

	idx := 0
	delKeyIdx := 0
	for i := range update {
		order := update[i]
		price := order[0]
		volume := order[1]

		alreadyIn := false
		switch {
		case volume == "0":
			delete(book.Book, price)
			alreadyIn = true

			for delKeyIdx < len(book.OrderKey) {
				if book.OrderKey[delKeyIdx] == price {
					book.OrderKey = append(book.OrderKey[:delKeyIdx], book.OrderKey[delKeyIdx+1:]...)
					break
				}
				delKeyIdx++
			}

		case volume != "0":
			if _, ok := book.Book[price]; ok {
				alreadyIn = true
				book.Book[price] = volume
			} else {
				book.Book[price] = volume
			}
		}

		if !alreadyIn {
			decPrice, _ := decimal.NewFromString(price)
			for j := idx; j < len(book.OrderKey); j++ {
				decKey, _ := decimal.NewFromString(book.OrderKey[j])
				if decPrice.LessThan(decKey) {
					book.OrderKey = append(append(book.OrderKey[:j], price), book.OrderKey[j:]...)
					idx = j
					break
				}
			}
		}

	}
}

func (book *AskBook) GetTop(n int) ([][]string, bool) {
	asks := make([][]string, 0, n)
	vanishKey := []string{}

	book.RLock()
	orderkey := book.OrderKey
	book.RUnlock()
	for i := range orderkey {
		price := orderkey[i]
		if _, ok := book.Book[price]; ok {
			volume := book.Book[price]
			if volume != "0" {
				asks = append(asks, []string{price, volume})
			}
		} else {
			vanishKey = append(vanishKey, price)
		}

		if len(asks) >= n {
			break
		}
	}

	k := 0
	i := 0

	book.Lock()
	for i < len(book.OrderKey) {
		for j := k; j < len(vanishKey); j++ {
			if book.OrderKey[i] == vanishKey[j] {
				book.OrderKey = append(book.OrderKey[:i], book.OrderKey[i+1:]...)
				i--
				k++
			}
		}
		if k == len(vanishKey)-1 {
			break
		}
		i++
	}
	book.Unlock()

	if len(asks) >= n {
		return asks, true
	} else {
		return asks, false
	}
}

func (book *AskBook) GetAll() ([][]string, bool) {
	var asks [][]string
	vanishKey := []string{}

	book.RLock()
	orderkey := book.OrderKey
	book.RUnlock()

	for i := range orderkey {
		price := orderkey[i]
		if _, ok := book.Book[price]; ok {
			volume := book.Book[price]
			if volume != "0" {
				asks = append(asks, []string{price, volume})
			}
		} else {
			vanishKey = append(vanishKey, price)
		}
	}

	k := 0
	i := 0

	book.Lock()
	for i < len(book.OrderKey) {
		for j := k; j < len(vanishKey); j++ {
			if book.OrderKey[i] == vanishKey[j] {
				book.OrderKey = append(book.OrderKey[:i], book.OrderKey[i+1:]...)
				i--
				k++
			}
		}
		if k == len(vanishKey)-1 {
			break
		}
		i++
	}
	book.Unlock()

	if len(asks) != 0 {
		return asks, true
	} else {
		return asks, false
	}
}

// bid
func NewBidBook(IsEvent bool) *BidBook {
	return &BidBook{
		OrderKey: []string{},
		Book:     map[string]string{},

		isEvent: IsEvent,
		event:   [][]string{},
	}
}

func (book *BidBook) Snapshot(snapshot [][]string) {
	orderkey := make([]string, len(snapshot))
	book.Lock()
	defer book.Unlock()

	// from small to big
	sort.Slice(snapshot, func(i, j int) bool {
		oi, _ := decimal.NewFromString(snapshot[i][0])
		oj, _ := decimal.NewFromString(snapshot[j][0])
		return oi.GreaterThan(oj)
	})

	for i := range snapshot {
		order := snapshot[i]
		price := order[0]
		volume := order[1]

		book.Book[price] = volume
		orderkey[i] = price
	}

	book.OrderKey = orderkey
}

func (book *BidBook) Update(update [][]string) {

	sort.Slice(update, func(i, j int) bool {
		oi, _ := decimal.NewFromString(update[i][0])
		oj, _ := decimal.NewFromString(update[j][0])
		return oi.GreaterThan(oj)
	})

	book.Lock()
	defer book.Unlock()

	idx := 0
	delKeyIdx := 0
	for i := range update {
		order := update[i]
		price := order[0]
		volume := order[1]

		alreadyIn := false
		switch {
		case volume == "0":
			delete(book.Book, price)
			alreadyIn = true

			for delKeyIdx < len(book.OrderKey) {
				if book.OrderKey[delKeyIdx] == price {
					book.OrderKey = append(book.OrderKey[:delKeyIdx], book.OrderKey[delKeyIdx+1:]...)
					break
				}
				delKeyIdx++
			}

		case volume != "0":
			if _, ok := book.Book[price]; ok {
				alreadyIn = true
				book.Book[price] = volume
			} else {
				book.Book[price] = volume
			}
		}

		if !alreadyIn {
			decPrice, _ := decimal.NewFromString(price)
			for j := idx; j < len(book.OrderKey); j++ {
				decKey, _ := decimal.NewFromString(book.OrderKey[j])
				if decPrice.GreaterThan(decKey) {
					book.OrderKey = append(append(book.OrderKey[:j], price), book.OrderKey[j:]...)
					idx = j
					break
				}
			}
		}

	}
}

func (book *BidBook) GetTop(n int) ([][]string, bool) {
	bids := make([][]string, 0, n)
	vanishKey := []string{}

	book.RLock()
	orderkey := book.OrderKey
	book.RUnlock()
	for i := range orderkey {
		price := orderkey[i]
		if _, ok := book.Book[price]; ok {
			volume := book.Book[price]
			if volume != "0" {
				bids = append(bids, []string{price, volume})
			}
		} else {
			vanishKey = append(vanishKey, price)
		}

		if len(bids) >= n {
			break
		}
	}

	k := 0
	i := 0

	book.Lock()
	for i < len(book.OrderKey) {
		for j := k; j < len(vanishKey); j++ {
			if book.OrderKey[i] == vanishKey[j] {
				book.OrderKey = append(book.OrderKey[:i], book.OrderKey[i+1:]...)
				i--
				k++
			}
		}
		if k == len(vanishKey)-1 {
			break
		}
		i++
	}
	book.Unlock()

	if len(bids) >= n {
		return bids, true
	} else {
		return bids, false
	}
}

func (book *BidBook) GetAll() ([][]string, bool) {
	var bids [][]string
	vanishKey := []string{}

	book.RLock()
	orderkey := book.OrderKey
	book.RUnlock()

	for i := range orderkey {
		price := orderkey[i]
		if _, ok := book.Book[price]; ok {
			volume := book.Book[price]
			if volume != "0" {
				bids = append(bids, []string{price, volume})
			}
		} else {
			vanishKey = append(vanishKey, price)
		}
	}

	k := 0
	i := 0

	book.Lock()
	for i < len(book.OrderKey) {
		for j := k; j < len(vanishKey); j++ {
			if book.OrderKey[i] == vanishKey[j] {
				book.OrderKey = append(book.OrderKey[:i], book.OrderKey[i+1:]...)
				i--
				k++
			}
		}
		if k == len(vanishKey)-1 {
			break
		}
		i++
	}
	book.Unlock()

	if len(bids) != 0 {
		return bids, true
	} else {
		return bids, false
	}
}
