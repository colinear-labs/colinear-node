// For any ambiguous chain

package xyz

import "math/big"

type Tx struct {
	txid   []byte
	to     []byte
	from   []byte
	amount *big.Int
}

type Block struct {
	txs []Tx
}

type Chain struct {
	name     string
	blocks10 []Block
}

func NewChain(name string) *Chain {
	c := Chain{name: name, blocks10: []Block{}}
	return &c
}

// Set all blocks at once
func (chain *Chain) SetBlocks(blocks []Block) {
	if len(blocks) == 10 {
		chain.blocks10 = blocks
	}
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *Chain) NewBlock(block Block) {
	chain.blocks10 = chain.blocks10[1 : len(chain.blocks10)-1]
	chain.blocks10 = append(chain.blocks10, block)
}
