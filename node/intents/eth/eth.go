package eth

import (
	"context"
	"fmt"
	"xnode/intents"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthProcessor struct {
	Id     string
	Port   uint
	Client *ethclient.Client
}

func NewEthProcessor(id string, port uint) *EthProcessor {
	return &EthProcessor{Id: id, Port: port}
}

func (p *EthProcessor) CurrencyId() string {
	return p.Id
}

func (p *EthProcessor) Process(intent *intents.PaymentIntent) chan intents.PaymentStatus {
	status := make(chan intents.PaymentStatus)

	go func() {
		headers := make(chan *types.Header)
		sub, err := p.Client.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			panic(err)
		}

		// Placeholder
		fmt.Println(sub)

	}()

	return status
}
