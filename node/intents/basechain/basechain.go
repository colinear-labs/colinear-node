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
	Txs []Tx
}

type BaseChain struct {
	HistoryLen     int
	PaymentTimeout uint // in seconds
	Blocks         []Block
	PendingTxs     []Tx
	Headers        []interface{}
}

func NewChain(historyLen int) *BaseChain {
	c := BaseChain{Blocks: []Block{}, HistoryLen: historyLen}
	return &c
}

// Set pending transactions
func (chain *BaseChain) SetPendingTxs(txs []Tx) {
	chain.PendingTxs = txs
}

// Set all blocks at once
func (chain *BaseChain) SetBlocks(blocks []Block) {
	if len(blocks) == chain.HistoryLen {
		chain.Blocks = blocks
	}
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *BaseChain) NewBlock(block Block) {
	if len(chain.Blocks) == chain.HistoryLen {
		chain.Blocks = chain.Blocks[1 : len(chain.Blocks)-1]
	}
	chain.Blocks = append(chain.Blocks, block)
}

// Set all headers at once
func (chain *BaseChain) SetHeaders(headers []interface{}) {
	if len(headers) == chain.HistoryLen {
		chain.Headers = headers
	}
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *BaseChain) NewHeader(header interface{}) {
	if len(chain.Headers) == chain.HistoryLen {
		chain.Headers = chain.Headers[1 : len(chain.Headers)-1]
	}
	chain.Headers = append(chain.Headers, header)
}

func (chain *BaseChain) Listen() {
	panic("Must be implemented by child classes.")
}
