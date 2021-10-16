package blockdata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"xnode/blockdata/chains"
)

var Btc *chains.BtcChain = nil
var Bch *chains.BtcChain = nil
var Ltc *chains.BtcChain = nil
var Doge *chains.BtcChain = nil
var Eth *chains.EthChain = nil

func InitChains(selectedChains []string) {
	for _, chain := range selectedChains {
		switch chain {
		case "btc":
			// Btc = chains.NewBtcChain("btc", 5000)
			Btc = chains.NewBtcChain("btc", 5000)
		case "bch":
			Btc = chains.NewBtcChain("bch", 5002)
		case "ltc":
			Btc = chains.NewBtcChain("ltc", 5003)
		case "doge":
			Btc = chains.NewBtcChain("doge", 5003)
		case "eth":
			Eth = chains.NewEthChain("eth", 5002)
		}
	}
}

// Spawn a local HTTP server to listen for new block events from
// BTC blockchains. Triggered by cURL requests from node containers,
// & pushes new events to the NewBlockEvents channel
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

		bmatch := make(map[string]*chains.BtcChain)
		bmatch["btc"] = Btc
		bmatch["bch"] = Bch
		bmatch["ltc"] = Ltc
		bmatch["doge"] = Doge

		go func() {
			bmatch[name].NewBlockEvents <- hash
		}()

	})

	http.ListenAndServe(":4999", nil)
}
