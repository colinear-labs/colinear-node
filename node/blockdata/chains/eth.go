// For all ETH-like chains

package chains

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthChain struct {
	Chain BaseChain
}

// Instantiate new ETH chain.
func NewEthChain(name string, port uint) *EthChain {
	c := EthChain{Chain: *NewChain(name, port)}
	return &c
}

// Listen for new blocks (over websocket to local node)
func (chain *EthChain) Listen() {
	client, err := ethclient.Dial(fmt.Sprintf("ws://127.0.0.1:%s/wsrpc", (string)(chain.Chain.Port)))
	if err != nil {
		panic(err)
	}

	headers := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), headers)

	if err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-sub.Err():
			panic(err)
		case header := <-headers:
			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				panic(err)
			}

			txs := []Tx{}

			for _, transaction := range block.Transactions() {
				tx := Tx{
					Txid:   transaction.Hash().String(),
					To:     transaction.To().String(),
					Amount: transaction.Value(),
				}
				txs = append(txs, tx)
			}

			chain.Chain.NewBlock(Block{Txs: txs})
		}
	}

}
