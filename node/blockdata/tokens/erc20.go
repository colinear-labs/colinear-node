package tokens

import (
	"fmt"
	"math/big"
	"xnode/blockdata/basechain"
	"xnode/blockdata/tokens/erc20abi"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20 struct {
	Chain           string
	Id              string
	ContractAddress string
	Txs             []basechain.Tx
}

var NewEthBlockEvents = make(chan uint)

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

func ListenForERC20EthPayment(
	ethClient *ethclient.Client,
	contractAddress string,
	walletAddress string,
	amount *big.Float,
	/* timeout uint, */
) bool {
	for {
		select {
		case <-NewEthBlockEvents:
			contract, err := erc20abi.NewERC20(common.HexToAddress(contractAddress), ethClient)
			if err != nil {
				panic(err)
			}
			bal, _ := contract.BalanceOf(nil, common.HexToAddress(walletAddress))

			amt := new(big.Float).SetInt(bal)
			amt = amt.Mul(amt, big.NewFloat(1e-18))
			if amt == amount {
				return true
			}
			// already implemented in intents.go ;)
			// case <-time.After(time.Second * time.Duration(timeout)):
			// 	return false
		}
	}
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
