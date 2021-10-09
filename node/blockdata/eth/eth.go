// For all ETH-like chains

package eth

import "math/big"

type EthTx struct {
	txid   []byte
	to     []byte
	from   []byte
	amount *big.Int
}

type EthBlock struct {
	txs []EthTx
}

type EthChain struct {
	name     string
	blocks20 []EthBlock
}

func NewEthChain(name string) *EthChain {
	c := EthChain{name: name, blocks20: []EthBlock{}}
	return &c
}

// Set all blocks at once
func (chain *EthChain) SetBlocks(blocks []EthBlock) {
	chain.blocks20 = blocks
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *EthChain) NewBlock(block EthBlock) {
	chain.blocks20 = chain.blocks20[1 : len(chain.blocks20)-1]
	chain.blocks20 = append(chain.blocks20, block)
}
