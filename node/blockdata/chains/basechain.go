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
	Name       string
	Port       uint
	Blocks10   []Block
	PendingTxs []Tx
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
	chain.Blocks10 = chain.Blocks10[1 : len(chain.Blocks10)-1]
	chain.Blocks10 = append(chain.Blocks10, block)
}

func (chain *BaseChain) Listen() {
	panic("Must be implemented by child classes.")
}
