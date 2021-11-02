package intents

import (
	"math/big"
)

type PaymentStatus int

const (
	Pending = iota
	Verified
)

type PaymentIntent struct {
	CurrencyId string
	Amount     *big.Float
	To         string
	Status     chan PaymentStatus
}

type Processor interface {
	CurrencyId() string
	Process(*PaymentIntent) chan PaymentStatus
}
