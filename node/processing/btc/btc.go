// Processor for BTC and all BTC forks, including:
//
// BCH, LTC, DOGE

package btc

import (
	"math/big"
	"time"

	"github.com/colinear-labs/colinear-node/processing"
	"github.com/colinear-labs/colinear-node/processing/basechain"
	"github.com/colinear-labs/colinear-node/xutil"
)

type BtcProcessor struct {
	Id                  string
	Port                uint
	Chain               *basechain.BaseChain
	NewBlockRpcEvents   chan string
	NewBlockLocalEvents chan string
}

func NewBtcProcessor(id string, port uint, historyLen uint) *BtcProcessor {

	p := &BtcProcessor{
		Id:    id,
		Port:  port,
		Chain: &basechain.BaseChain{},
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

func (p *BtcProcessor) Process(intent *processing.PaymentIntentLocal) chan xutil.PaymentStatus {
	status := make(chan xutil.PaymentStatus)

	go func() {

		// Pending txs checking loop

		pSum := big.NewFloat(0)
	pendingLoop:
		for {

			for _, p := range p.Chain.PendingTxs {
				if p.To == intent.To {
					pSum = pSum.Add(pSum, p.Amount)
					compare := intent.Amount.Cmp(pSum)
					if compare > 0 {
						status <- xutil.Verified
						break pendingLoop
					}
				}
			}

			time.Sleep(1 * time.Second)
		}

		// Verified txs checking loop

		vSum := big.NewFloat(0)
	verifiedLoop:
		for {
			<-p.NewBlockLocalEvents // wait for new local block event

			for _, x := range p.Chain.LatestBlock.Txs {
				if x.To == intent.To {
					vSum = vSum.Add(vSum, x.Amount)
					compare := intent.Amount.Cmp(vSum)
					if compare > 0 {
						status <- xutil.Verified
						break verifiedLoop
					}
				}
			}
		}
	}()

	return status
}
