package intents

import (
	"fmt"
	"math/big"
	"time"
	"xnode/blockdata"
)

type CurrencyType int

const (
	Coin = iota
	CoinObfuscated
	Erc20
)

type PaymentStatus int

const (
	Incomplete = iota
	Pending
	Verified
)

type PaymentIntent struct {
	Type       CurrencyType
	CurrencyId string
	Amount     *big.Float
	To         string
	Status     PaymentStatus
}

func (intent PaymentIntent) WaitForPayment() {

	timeinterval := 5 // seconds
	res := make(chan string, 1)

	go func() {
		found := false
		for {
			// scan for pending / verified payments here
			if intent.Status == Pending {
				found = findVerifiedPayment(intent)
				if found {
					intent.Status = Verified
					res <- "Found"
					return
				}
			} else if intent.Status == Incomplete {
				found = findPendingPayment(intent)
				if found {
					intent.Status = Pending
				}
			}
			time.After(time.Duration(timeinterval) * time.Second)
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

func findPendingPayment(intent PaymentIntent) bool {
	switch intent.Type {
	case Coin:
		if blockdata.ChainDict[intent.CurrencyId] != nil {
			for _, tx := range blockdata.ChainDict[intent.CurrencyId].PendingTxs {
				compareAmounts := tx.Amount.Cmp(intent.Amount)
				if tx.To == intent.To && (compareAmounts == 1 || compareAmounts == 0) {
					return true
				}
			}
		} else {
			fmt.Printf("%s doesn't match any blockchains.", intent.CurrencyId)
		}
	}
	return false
}

func findVerifiedPayment(intent PaymentIntent) bool {
	switch intent.Type {
	case Coin:
		if blockdata.ChainDict[intent.CurrencyId] != nil {
			for _, block := range blockdata.ChainDict[intent.CurrencyId].Blocks10 {
				for _, tx := range block.Txs {
					compareAmounts := tx.Amount.Cmp(intent.Amount)
					if tx.To == intent.To && (compareAmounts == 1 || compareAmounts == 0) {
						return true
					}
				}
			}
		} else {
			// maybe switch this to whatever logging system noise uses
			fmt.Printf("%s doesn't match any blockchains.", intent.CurrencyId)
		}
	}
	return false
}
