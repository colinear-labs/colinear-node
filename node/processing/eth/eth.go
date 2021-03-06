package eth

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/colinear-labs/colinear-node/processing"
	"github.com/colinear-labs/colinear-node/xutil"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthProcessor struct {
	Id             string
	Port           uint
	Client         *ethclient.Client
	NewBlockEvents chan *types.Header
}

func NewEthProcessor(id string, port uint) *EthProcessor {
	client, err := ethclient.Dial(fmt.Sprintf("ws://127.0.0.1:%s/wsrpc", fmt.Sprint(port)))
	if err != nil {
		panic(err)
	}

	return &EthProcessor{Id: id, Port: port, Client: client}
}

func (p *EthProcessor) CurrencyId() string {
	return p.Id
}

func (p *EthProcessor) Process(intent *processing.PaymentIntentLocal) chan xutil.PaymentStatus {
	status := make(chan xutil.PaymentStatus)

	amtInt64, _ := intent.Amount.Int64()
	amtEth := big.NewInt(amtInt64)
	toEth := common.HexToAddress(intent.To)

	go func() {

		// Pending transaction Loop

	secondsLoop:
		for {
			balance, _ := p.Client.PendingBalanceAt(context.Background(), toEth)
			comparison := balance.Cmp(amtEth)
			if comparison == 1 || comparison == 0 {
				status <- xutil.Verified
				break secondsLoop
			}

			time.Sleep(1 * time.Second)
		}

		// Verified Transaction Loop

		headers := make(chan *types.Header)
		sub, err := p.Client.SubscribeNewHead(context.Background(), headers)
		if err != nil {
			panic(err)
		}

	headerLoop:
		for {
			select {
			case err := <-sub.Err():
				panic(err)
			case header := <-headers:
				balance, _ := p.Client.BalanceAt(context.Background(), toEth, header.Number)
				comparison := balance.Cmp(amtEth)
				if comparison == 1 || comparison == 0 {
					status <- xutil.Verified
					break headerLoop
				}
			}
		}

	}()

	return status
}
