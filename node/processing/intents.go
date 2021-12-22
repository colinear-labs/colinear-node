package processing

import (
	"math/big"

	"github.com/colinear-labs/colinear-node/xutil"
)

// Node-side payment intent representation.
type PaymentIntentLocal struct {
	CurrencyId string
	Amount     *big.Float
	To         string
	Status     chan xutil.PaymentStatus
}

func NewPaymentIntentLocal(currencyId string, amount *big.Float, to string) *PaymentIntentLocal {
	return &PaymentIntentLocal{
		CurrencyId: currencyId,
		Amount:     amount,
		To:         to,
		Status:     make(chan xutil.PaymentStatus),
	}
}

type Processor interface {
	CurrencyId() string
	Process(*PaymentIntentLocal) chan xutil.PaymentStatus
}
