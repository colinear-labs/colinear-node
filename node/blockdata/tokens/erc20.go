package tokens

import (
	"fmt"
	"math/big"
	"xnode/blockdata/chains"
)

type ERC20 struct {
	Chain           string
	Id              string
	ContractAddress string
	Txs             []chains.Tx
}

func ValidChainERC20(chain string) bool {
	validChains := []string{"eth"}
	for _, refchain := range validChains {
		if chain == refchain {
			return true
		}
	}
	return false
}

func NewERC20Eth(id string, contractAddress string) *ERC20 {
	return &ERC20{Id: id, ContractAddress: contractAddress}
}

func (token *ERC20) FindPendingPayment(address string, amount big.Float) bool {
	switch token.Chain {
	case "eth":
		fmt.Print("Placeholder")
		return true
	default:
		return false
	}

}

func (token *ERC20) FindVerifiedPayment(address string, amount big.Float) bool {
	switch token.Chain {
	case "eth":
		fmt.Print("Placeholder")
		return true
	default:
		return false
	}
}
