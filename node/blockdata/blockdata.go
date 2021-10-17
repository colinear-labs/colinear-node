package blockdata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"xnode/blockdata/chains"
	"xnode/nodeutil"
)

var Btc *chains.BtcChain = nil
var Bch *chains.BtcChain = nil
var Ltc *chains.BtcChain = nil
var Doge *chains.BtcChain = nil
var Eth *chains.EthChain = nil

func InitChains(selectedChains []string) {

	selectedChains = nodeutil.Unique(selectedChains)
	fmt.Println(selectedChains)

	btcNotify := false
	for _, chain := range selectedChains {
		switch chain {
		case "btc":
			Btc = chains.NewBtcChain("btc", 5000)
			go func() { Btc.Listen() }()
			btcNotify = true
		case "bch":
			Bch = chains.NewBtcChain("bch", 5002)
			go func() { Bch.Listen() }()
			btcNotify = true
		case "ltc":
			Ltc = chains.NewBtcChain("ltc", 5003)
			go func() { Ltc.Listen() }()
			btcNotify = true
		case "doge":
			Doge = chains.NewBtcChain("doge", 5003)
			go func() { Doge.Listen() }()
			btcNotify = true
		case "eth":
			Eth = chains.NewEthChain("eth", 5001)
			go func() { Eth.Listen() }()
		default:
			fmt.Printf("WARN: Unrecognized chain %s.\n", chain)
		}
	}
	if btcNotify {
		go func() { SpawnBtcBlockNotifyServer() }()
	}
}

// Spawn a local HTTP server to listen for new block events from
// BTC blockchains. Triggered by cURL requests from node containers,
// & pushes new events to the NewBlockEvents channel
func SpawnBtcBlockNotifyServer() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("HTTP ERROR %s", err)
		}

		str := (string)(body)
		u, _ := url.ParseQuery(str)
		name := u["name"][0]
		hash := u["hash"][0]

		var ch *chains.BtcChain = nil

		switch name {
		case "btc":
			ch = Btc
		case "bch":
			ch = Bch
		case "ltc":
			ch = Ltc
		case "doge":
			ch = Doge
		}

		go func() {
			ch.NewBlockEvents <- hash
		}()

	})

	http.ListenAndServe("127.0.0.1:4999", nil)
}
