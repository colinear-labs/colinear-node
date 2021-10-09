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
