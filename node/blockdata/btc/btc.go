// For all BTC-like chains

package btc

import (
	"fmt"
	"math/big"

	"github.com/pebbe/zmq4"
)

type Tx struct {
	txid   string
	to     string
	from   string
	amount *big.Int
}

type BtcBlock struct {
	txs []Tx
}

type BtcChain struct {
	name     string
	blocks10 []BtcBlock
}

func NewBtcChain(name string) *BtcChain {
	c := BtcChain{name: name, blocks10: []BtcBlock{}}
	return &c
}

// Set all blocks at once
func (chain *BtcChain) SetBlocks(blocks []BtcBlock) {
	if len(blocks) == 10 {
		chain.blocks10 = blocks
	}
}

// Cycles old old block `[0]` out and appends the latest block
func (chain *BtcChain) NewBlock(block BtcBlock) {
	if len(chain.blocks10) == 10 {
		chain.blocks10 = chain.blocks10[1 : len(chain.blocks10)-1]
	}
	chain.blocks10 = append(chain.blocks10, block)
}
func (chain *BtcChain) Listen(port uint) {
	socket, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		panic(err)
	}
	defer socket.Close()
	// if we need filtering, check out http://api.zeromq.org/4-1:zmq-setsockopt#toc41

	socket.Bind(fmt.Sprintf("tcp://127.0.0.1:%s", (string)(port)))
	for {
		msg, _ := socket.Recv(0)
		fmt.Printf("Received %s\n", msg)
	}
}
