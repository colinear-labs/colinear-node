// For all BTC-like chains

package btc

import "math/big"

type BtcTx struct {
	txid   []byte
	to     []byte
	from   []byte
	amount *big.Int
}

type BtcBlock struct {
	txs []BtcTx
}

type BtcChain struct {
	name     string
	blocks20 []BtcBlock
}

func NewBtcChain(name string) *BtcChain {
	c := BtcChain{name: name, blocks20: []BtcBlock{}}
	return &c
}
