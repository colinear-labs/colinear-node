// For all BTC-like chains

package btc

import "math/big"

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
	blocks20 []BtcBlock
}

func NewBtcChain(name string) *BtcChain {
	c := BtcChain{name: name, blocks20: []BtcBlock{}}
	return &c
}

// Set all blocks at once
func (chain *BtcChain) SetBlocks(blocks []BtcBlock) {
	if len(blocks) == 20 {
		chain.blocks20 = blocks
	}
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *BtcChain) NewBlock(block BtcBlock) {
	chain.blocks20 = chain.blocks20[1 : len(chain.blocks20)-1]
	chain.blocks20 = append(chain.blocks20, block)
}
