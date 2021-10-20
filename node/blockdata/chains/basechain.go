package chains

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
	Name           string
	Port           uint
	PaymentTimeout uint // in seconds
	Blocks10       []Block
	PendingTxs     []Tx
	Headers10      []interface{}
}

func NewChain(name string, port uint) *BaseChain {
	c := BaseChain{Name: name, Port: port, Blocks10: []Block{}}
	return &c
}

// Set pending transactions
func (chain *BaseChain) SetPendingTxs(txs []Tx) {
	chain.PendingTxs = txs
}

// Set all blocks at once
func (chain *BaseChain) SetBlocks(blocks []Block) {
	if len(blocks) == 10 {
		chain.Blocks10 = blocks
	}
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *BaseChain) NewBlock(block Block) {
	if len(chain.Blocks10) == 10 {
		chain.Blocks10 = chain.Blocks10[1 : len(chain.Blocks10)-1]
	}
	chain.Blocks10 = append(chain.Blocks10, block)
}

// Set all headers at once
func (chain *BaseChain) SetHeaders(headers []interface{}) {
	if len(headers) == 10 {
		chain.Headers10 = headers
	}
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *BaseChain) NewHeader(header interface{}) {
	if len(chain.Headers10) == 10 {
		chain.Headers10 = chain.Headers10[1 : len(chain.Headers10)-1]
	}
	chain.Headers10 = append(chain.Headers10, header)
}

func (chain *BaseChain) Listen() {
	panic("Must be implemented by child classes.")
}
