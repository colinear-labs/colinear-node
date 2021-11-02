package btc

import "xnode/intents"

type BtcProcessor struct {
	Id   string
	Port uint
}

func NewBtcProcessor(id string, port uint) *BtcProcessor {
	return &BtcProcessor{Id: id, Port: port}
}

func (p *BtcProcessor) CurrencyId() string {
	return p.Id
}

func (p *BtcProcessor) Process(intent *intents.PaymentIntent) chan intents.PaymentStatus {
	status := make(chan intents.PaymentStatus)

	go func() {

	}()

	return status
}
