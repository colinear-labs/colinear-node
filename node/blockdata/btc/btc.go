// For all BTC-like chains

package btc

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/rpc/jsonrpc"
	"net/url"

	"github.com/pebbe/zmq4"
)

type BtcTx struct {
	txid   string
	to     string
	from   string
	amount *big.Int
}

type BtcBlock struct {
	txs []BtcTx
}

type BtcChain struct {
	name           string
	blocks10       []BtcBlock
	NewBlockEvents chan BtcBlock
}

func NewBtcChain(name string) *BtcChain {
	c := BtcChain{name: name, blocks10: []BtcBlock{}, NewBlockEvents: make(chan BtcBlock)}
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

// (ZMQ) Listen on given port for incoming blocks
//
// Cleaner than local JSON-RPC, but
// the bitcoin forks leave us no choice :c
func (chain *BtcChain) ListenZmq(port uint) {
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

// (JSON-RPC) Continuously query bitcoind JSON-RPC server over given port
// to get new blocks
//
// Auth login should simply be user:pass
func (chain *BtcChain) Listen(port uint) {

	// Make JSON-RPC connection
	client, err := jsonrpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", (string)(port)))
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// testhash := "00000000000000000009493672ee4919250e6837cd65c09157d8858043af5531"

	for {

		header := <-chain.NewBlockEvents

		reply := []byte{}

		err := client.Call("getblock", header, &reply)
		if err != nil {
			panic(err)
		} else {
			fmt.Println(reply)
		}
	}

}

func SpawnBtcBlockNotifyServer(port uint) {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("HTTP ERROR %s", err)
		}

		str := (string)(body)
		u, _ := url.ParseQuery(str)
		name := u["name"][0]
		hash := u["hash"][0]

		fmt.Printf("CHAIN NAME\t%s\nBLOCK HASH\t%s\n\n", name, hash)
	})

	http.ListenAndServe(":4999", nil)
}
