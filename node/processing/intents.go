package processing

import (
	"math/big"
)

type PaymentStatus int

const (
	Pending = iota
	Verified
)

// Node-side payment intent representation.
type PaymentIntentNS struct {
	CurrencyId string
	Amount     *big.Float
	To         string
	Status     chan PaymentStatus
}

type Processor interface {
	CurrencyId() string
	Process(*PaymentIntentNS) chan PaymentStatus
}
