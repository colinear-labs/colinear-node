// Processor for BTC and all BTC forks, including:
//
// BCH, LTC, DOGE

package btc

import (
	"xnode/processing"
	"xnode/processing/basechain"
)

type BtcProcessor struct {
	Id             string
	Port           uint
	Chain          *basechain.BaseChain
	NewBlockEvents chan string
}

func NewBtcProcessor(id string, port uint, historyLen uint) *BtcProcessor {

	p := &BtcProcessor{
		Id:    id,
		Port:  port,
		Chain: basechain.NewChain((int)(historyLen)),
	}

	go func() {
		JsonRpcListenMempool(p)
	}()

	go func() {
		JsonRpcListenBlocks(p)
	}()

	return p
}

func (p *BtcProcessor) CurrencyId() string {
	return p.Id
}

func (p *BtcProcessor) Process(intent *processing.PaymentIntentNS) chan processing.PaymentStatus {
	status := make(chan processing.PaymentStatus)

	return status
}
