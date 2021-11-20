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
type PaymentIntentLocal struct {
	CurrencyId string
	Amount     *big.Float
	To         string
	Status     chan PaymentStatus
}

func NewPaymentIntentLocal(currencyId string, amount *big.Float, to string) *PaymentIntentLocal {
	return &PaymentIntentLocal{
		CurrencyId: currencyId,
		Amount:     amount,
		To:         to,
		Status:     make(chan PaymentStatus),
	}
}

type Processor interface {
	CurrencyId() string
	Process(*PaymentIntentLocal) chan PaymentStatus
}
