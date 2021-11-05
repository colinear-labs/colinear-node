package erc20eth

import (
	"xnode/intents"

	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20EthProcessor struct {
	Id              string
	ContractAddress string
	Client          *ethclient.Client
}

func NewERC20EthProcessor(id string, contractAddress string, client *ethclient.Client) *ERC20EthProcessor {
	return &ERC20EthProcessor{
		Id:              id,
		ContractAddress: contractAddress,
		Client:          client,
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
