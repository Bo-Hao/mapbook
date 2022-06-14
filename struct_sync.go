package mapbook

import "sync"

type AskBook struct {
	Book sync.Map

	isEvent bool
	event   [][]string
}

type BidBook struct {
	Book sync.Map

	isEvent bool
	event   [][]string
}
