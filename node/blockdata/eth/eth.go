// For all ETH-like chains

package eth

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthTx struct {
	txid   string
	to     string
	from   string
	amount *big.Int
}

type EthBlock struct {
	txs []EthTx
}

type EthChain struct {
	name     string
	blocks10 []EthBlock
}

func NewEthChain(name string) EthChain {
	c := EthChain{name: name, blocks10: []EthBlock{}}
	return c
}

// Set all blocks at once
func (chain *EthChain) setBlocks(blocks []EthBlock) {
	if len(blocks) == 10 {
		chain.blocks10 = blocks
	}
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *EthChain) newBlock(block EthBlock) {
	chain.blocks10 = chain.blocks10[1 : len(chain.blocks10)-1]
	chain.blocks10 = append(chain.blocks10, block)
}

func (chain *EthChain) Listen(port uint) {
	client, err := ethclient.Dial("ws://127.0.0.1:5001/wsrpc")
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

			txs := []EthTx{}

			for _, transaction := range block.Transactions() {
				tx := EthTx{
					txid:   transaction.Hash().String(),
					to:     transaction.To().String(),
					amount: transaction.Value(),
				}
				txs = append(txs, tx)
			}

			chain.newBlock(EthBlock{txs: txs})
		}
	}

}
