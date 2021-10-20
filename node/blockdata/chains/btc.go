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

type BtcHeader struct {
	Hash              string
	Version           int
	Merkleroot        string
	Time              int
	Nonce             int
	Bits              string
	Difficulty        float32
	Previousblockhash string
}

// Get block headers via JSON-RPC

type getBlockHeaderJsonRequest struct {
	Jsonrpc string             `json:"jsonrpc"`
	Id      string             `json:"id"`
	Method  string             `json:"method"`
	Params  getBlockHeaderArgs `json:"params"`
}

type getBlockHeaderArgs struct {
	Hash    string `json:"hash"`
	Verbose bool   `json:"verbose"`
}

type getBlockHeaderResponse struct {
	Result getBlockHeaderResult `json:"result"`
	Error  interface{}          `json:"error"`
	Id     interface{}          `json:"id"`
}

type getBlockHeaderResult struct {
	Hash              string  `json:"hash"`
	Confirmations     int     `json:"confirmations"`
	Height            int     `json:"height"`
	Version           int     `json:"version"`
	VersionHex        string  `json:"versionHex"`
	Merkleroot        string  `json:"merkleRoot"`
	Time              int     `json:"time"`
	MedianTime        int     `json:"medianTime"`
	Nonce             int     `json:"nonce"`
	Bits              string  `json:"bits"`
	Difficulty        float32 `json:"difficulty"`
	Chainwork         string  `json:"chainwork"`
	Ntx               int     `json:"nTx"`
	Previousblockhash string  `json:"previousblockhash"`
	Nextblockhash     string  `json:"nextblock"`
}

// Get full blocks JSON-RPC

type getBlockJsonRequest struct {
	Jsonrpc string       `json:"jsonrpc"`
	Id      string       `json:"id"`
	Method  string       `json:"method"`
	Params  getBlockArgs `json:"params"`
}

type getBlockArgs struct {
	Blockhash string `json:"blockhash"`
	Verbosity int    `json:"verbosity"`
}

type getBlockResponse struct {
	Result getBlockResult `json:"result"`
}

type getBlockResult struct {
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

// [NOT FULLY IMPLEMENTED] (ZMQ) Listen on given port
// for incoming blocks
//
// Cleaner than blocknotify + local JSON-RPC, but
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

		blockHash := <-chain.NewBlockEvents

		resultGb := getBlockResponse{}

		optsGb := getBlockArgs{Blockhash: blockHash, Verbosity: 2}
		payloadGb := getBlockJsonRequest{Jsonrpc: "1.0", Id: "xyz", Method: "getblock", Params: optsGb}
		respGb, err := req.Post("http://user:pass@127.0.0.1:5003", req.BodyJSON(&payloadGb))

		if err != nil {
			panic(err)
		}

		respGb.ToJSON(&resultGb)

		resultGbh := getBlockHeaderResponse{}
		optsGbh := getBlockHeaderArgs{Hash: blockHash, Verbose: true}
		payloadGbh := getBlockHeaderJsonRequest{Jsonrpc: "1.0", Id: "xyz", Method: "getblockheader", Params: optsGbh}
		respGbh, err := req.Post("http://user:pass@127.0.0.1:5003", req.BodyJSON(&payloadGbh))

		if err != nil {
			panic(err)
		}

		respGbh.ToJSON(&resultGbh)

		newBlock := Block{}
		for _, tx := range resultGb.Result.Tx {

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

		rgbh := resultGbh.Result
		newHeader := BtcHeader{
			Hash:              rgbh.Hash,
			Version:           rgbh.Version,
			Merkleroot:        rgbh.Merkleroot,
			Time:              rgbh.Time,
			Nonce:             rgbh.Time,
			Bits:              rgbh.Bits,
			Difficulty:        rgbh.Difficulty,
			Previousblockhash: rgbh.Previousblockhash,
		}

		chain.Chain.NewBlock(newBlock)
		chain.Chain.NewHeader(newHeader)
	}

}
