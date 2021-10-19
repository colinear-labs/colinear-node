package intents

import (
	"fmt"
	"math/big"
	"xnode/blockdata"
)

type CurrencyType int

const (
	Coin = iota
	CoinObfuscated
	Erc20
)

type PaymentIntent struct {
	Type       CurrencyType
	CurrencyId string
	Amount     *big.Float
	To         string
}

func FindPendingPayment(intent PaymentIntent) bool {
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

func FindVerifiedPayment(intent PaymentIntent) bool {
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
