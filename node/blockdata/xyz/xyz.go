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
	blocks20 []Block
}

func NewChain(name string) *Chain {
	c := Chain{name: name, blocks20: []Block{}}
	return &c
}
