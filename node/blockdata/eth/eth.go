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
	blocks10 []EthBlock
}

func NewEthChain(name string) *EthChain {
	c := EthChain{name: name, blocks10: []EthBlock{}}
	return &c
}

// Set all blocks at once
func (chain *EthChain) SetBlocks(blocks []EthBlock) {
	chain.blocks10 = blocks
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *EthChain) NewBlock(port uint) {

	block := EthBlock{} // change later
	chain.blocks10 = chain.blocks10[1 : len(chain.blocks10)-1]
	chain.blocks10 = append(chain.blocks10, block)
}
