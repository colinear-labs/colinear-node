// For all BTC-like chains

package chains

import (
	"fmt"
	"math/big"

	"github.com/imroc/req"
	"github.com/pebbe/zmq4"
)

// Base BTC chain type
type BtcChain struct {
	NewBlockEvents chan string
	Chain          BaseChain
}

// Data schema for working with JSON-RPC
type jsonRequest struct {
	Jsonrpc string       `json:"jsonrpc"`
	Id      string       `json:"id"`
	Method  string       `json:"method"`
	Params  getBlockArgs `json:"params"`
}

type getBlockArgs struct {
	Blockhash string `json:"blockhash"`
	Verbosity int    `json:"verbosity"`
}

type response struct {
	Result result `json:"result"`
}

type result struct {
	Version           int     `json:"version"`
	Height            int     `json:"height"`
	Confirmations     int     `json:"confirmations"`
	Merkleroot        string  `json:"merkleroot"`
	Tx                []tx    `json:"tx"`
	Time              int     `json:"time"`
	Nonce             int     `json:"nonce"`
	Difficulty        float32 `json:"difficulty"`
	Previousblockhash string  `json:"previousblockhash"`
}

type tx struct {
	Txid string `json:"txid"`
	Hash string `json:"hash"`
	Vout []vout `json:"vout"`
}

type vout struct {
	Value        float32      `json:"value"`
	ScriptPubKey scriptPubKey `json:"scriptPubKey"`
}

type scriptPubKey struct {
	Hex       string   `json:"hex"`
	Addresses []string `json:"addresses"`
}

// Helpers

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

	for {

		result := response{}

		blockHash := <-chain.NewBlockEvents

		opts := getBlockArgs{Blockhash: blockHash, Verbosity: 2}
		payload := jsonRequest{Jsonrpc: "1.0", Id: "xyz", Method: "getblock", Params: opts}
		resp, err := req.Post("http://user:pass@127.0.0.1:5003", req.BodyJSON(&payload))

		if err != nil {
			panic(err)
		}

		resp.ToJSON(&result)

		fmt.Println(result.Result.Version)

		newBlock := Block{}
		for _, tx := range result.Result.Tx {

			addresses := tx.Vout[0].ScriptPubKey.Addresses
			addr := ""
			if len(addresses) != 0 {
				addr = addresses[0]
			}

			newTx := Tx{
				Txid:   tx.Txid,
				Amount: big.NewFloat((float64)(tx.Vout[0].Value)),
				To:     addr,
			}
			newBlock.Txs = append(newBlock.Txs, newTx)
		}
	}

}
