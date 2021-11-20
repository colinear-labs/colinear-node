// Base chain for storing & parsing recent blocks. Intended
// primarily for usage with "brittle" nodes like bitcoind,
// litecoind, etc.
package basechain

import "math/big"

type Tx struct {
	Txid   string
	To     string
	From   string
	Amount *big.Float
}

type Block struct {
	Hash string
	Txs  []Tx
}

type BaseChain struct {
	PaymentTimeout uint // in seconds
	PendingTxs     []Tx
	LatestBlock    Block
	LatestHeader   interface{}
}

func (c *BaseChain) SetLatestBlock(b Block) {
	c.LatestBlock = b
}

func (c *BaseChain) SetLatestHeader(h interface{}) {
	c.LatestHeader = h
}
