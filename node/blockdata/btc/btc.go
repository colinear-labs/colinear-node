// For all BTC-like chains

package btc

import (
	"math/big"
)

type Tx struct {
	txid   []byte
	to     []byte
	from   []byte
	amount *big.Int
}

type BtcBlock struct {
	txs []Tx
}

type BtcChain struct {
	name     string
	blocks10 []BtcBlock
}

func NewBtcChain(name string) *BtcChain {
	c := BtcChain{name: name, blocks10: []BtcBlock{}}
	return &c
}

// Set all blocks at once
func (chain *BtcChain) SetBlocks(blocks []BtcBlock) {
	if len(blocks) == 10 {
		chain.blocks10 = blocks
	}
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *BtcChain) NewBlock(block BtcBlock) {
	if len(chain.blocks10) == 10 {
		chain.blocks10 = chain.blocks10[1 : len(chain.blocks10)-1]
	}
	chain.blocks10 = append(chain.blocks10, block)
}
