// For all BTC-like chains

package chains

import (
	"fmt"
	"net/rpc/jsonrpc"

	"github.com/pebbe/zmq4"
)

type BtcChain struct {
	NewBlockEvents chan string
	Chain          BaseChain
}

func NewBtcChain(name string, port uint) *BtcChain {
	c := BtcChain{NewBlockEvents: make(chan string), Chain: *NewChain(name, port)}
	return &c
}

// (ZMQ) Listen on given port for incoming blocks
//
// Cleaner than local JSON-RPC, but
// the bitcoin forks leave us no choice :c
func (chain *BtcChain) ListenZmq() {
	socket, err := zmq4.NewSocket(zmq4.SUB)
	if err != nil {
		panic(err)
	}
	defer socket.Close()
	// if we need filtering, check out http://api.zeromq.org/4-1:zmq-setsockopt#toc41

	socket.Bind(fmt.Sprintf("tcp://127.0.0.1:%s", (string)(chain.Chain.Port)))
	for {
		msg, _ := socket.Recv(0)
		fmt.Printf("Received %s\n", msg)
	}
}

// (JSON-RPC) Query bitcoind JSON-RPC server over given port
// to get new blocks whenever new block hashes are passed
// over the NewBlockEvents channel
//
// Auth login should simply be user:pass
func (chain *BtcChain) Listen() {

	// Make JSON-RPC connection
	client, err := jsonrpc.Dial("tcp", fmt.Sprintf("127.0.0.1:%s", (string)(chain.Chain.Port)))
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
