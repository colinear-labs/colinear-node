package erc20eth

import (
	"xnode/intents"

	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20EthProcessor struct {
	Id              string
	Client          *ethclient.Client
	ContractAddress string
}

func NewERC20EthProcessor(id string, contractAddress string) *ERC20EthProcessor {
	return &ERC20EthProcessor{
		Id:              id,
		ContractAddress: contractAddress,
	}
}

func (p *ERC20EthProcessor) CurrencyId() string {
	return p.Id
}

func (p *ERC20EthProcessor) Process(intent *intents.PaymentIntent) chan intents.PaymentStatus {
	status := make(chan intents.PaymentStatus)

	go func() {

	}()

	return status
}
