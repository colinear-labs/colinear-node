package intents

import (
	"fmt"
	"math/big"
	"time"
	"xnode/blockdata"
	"xnode/blockdata/tokens"
)

type CurrencyType int

const (
	Coin = iota
	CoinObfuscated
	Erc20Eth
)

type PaymentStatus int

const (
	Incomplete = iota
	Pending
	Verified
)

type PaymentIntent struct {
	CurrencyId string
	Amount     *big.Float
	To         string
	Status     PaymentStatus
}

func (intent PaymentIntent) WaitForPayment() {

	res := make(chan string, 1)

	go func() {
		switch determineIntentType(intent) {
		case Coin:
			timeinterval := 1 // seconds
			found := false
			for {
				// scan for pending / verified payments here
				if intent.Status == Pending {
					found = findVerifiedCoinPayment(intent)
					if found {
						intent.Status = Verified
						res <- "Found"
						return
					}
				} else if intent.Status == Incomplete {
					found = findPendingCoinPayment(intent)
					if found {
						intent.Status = Pending
					}
				}
				time.After(time.Duration(timeinterval) * time.Second)
			}
		case Erc20Eth:
			tokens.ListenForERC20EthPayment(
				blockdata.Eth.Client,
				blockdata.ERC20EthDict[intent.CurrencyId].ContractAddress,
				intent.To,
				intent.Amount,
			)
		}
	}()

	// get chain timeout
	timeout := 100

	// wait for result or timeout
	select {
	case <-res:
		return
	case <-time.After(time.Duration(timeout) * time.Second):
		return
	}
}

// Check if total in mempool amounts to price
//
// NOTE: We COULD total across both mempool and
// recent blocks. idk if good idea tho cuz extra
// work for computer :(
func findPendingCoinPayment(intent PaymentIntent) bool {
	if blockdata.ChainDict[intent.CurrencyId] != nil {
		total := big.NewFloat(0.0)
		for _, tx := range blockdata.ChainDict[intent.CurrencyId].PendingTxs {
			if tx.To == intent.To {
				total.Add(total, tx.Amount)
				compareAmounts := total.Cmp(intent.Amount)
				if compareAmounts == 1 || compareAmounts == 0 {
					return true
				}
			}
		}
	} else {
		fmt.Printf("%s doesn't match any running currencies.", intent.CurrencyId)
	}
	return false
}

// Check if total in last 10 blocks amounts to price
func findVerifiedCoinPayment(intent PaymentIntent) bool {
	if blockdata.ChainDict[intent.CurrencyId] != nil {
		total := big.NewFloat(0.0)
		for _, block := range blockdata.ChainDict[intent.CurrencyId].Blocks10 {
			for _, tx := range block.Txs {
				if tx.To == intent.To {
					total.Add(total, tx.Amount)
					compareAmounts := total.Cmp(intent.Amount)
					if compareAmounts == 1 || compareAmounts == 0 {
						return true
					}
				}
			}
		}
	} else {
		// maybe switch this to whatever logging system noise uses
		fmt.Printf("%s doesn't match any running currencies.", intent.CurrencyId)
	}
	return false
}

func determineIntentType(intent PaymentIntent) CurrencyType {
	var res CurrencyType
	if blockdata.ChainDict[intent.CurrencyId] != nil {
		res = Coin
	} else if blockdata.ERC20EthDict[intent.CurrencyId] != nil {
		res = Erc20Eth
	}
	return res
}
