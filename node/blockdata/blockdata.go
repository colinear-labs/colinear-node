// CONTAINS ALL SUPPORTED PUBLIC CHAINS
package blockdata

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"xnode/blockdata/chains"
	"xnode/blockdata/tokens"
	"xnode/nodeutil"
)

// Bitcoin & forks

var Btc *chains.BtcChain = nil
var Bch *chains.BtcChain = nil
var Ltc *chains.BtcChain = nil
var Doge *chains.BtcChain = nil

// Ethereum chain(s)

var Eth *chains.EthChain = nil

// ERC20s
//
// NOTE: Start with stablecoins for now

var Dai *tokens.ERC20 = nil
var Usdt *tokens.ERC20 = nil
var Usdc *tokens.ERC20 = nil
var Ust *tokens.ERC20 = nil

var ChainDict map[string]*chains.BaseChain
var ERC20Dict map[string]*tokens.ERC20

func InitChains(selectedChains []string) {

	selectedChains = nodeutil.Unique(selectedChains)
	fmt.Println(selectedChains)

	btcNotify := false
	for _, chain := range selectedChains {
		switch chain {
		case "btc":
			Btc = chains.NewBtcChain("btc", 5000)
			ChainDict["btc"] = &Btc.Chain
			go func() { Btc.Listen() }()
			btcNotify = true // Only for BTC-based chains
		case "bch":
			Bch = chains.NewBtcChain("bch", 5002)
			ChainDict["bch"] = &Bch.Chain
			go func() { Bch.Listen() }()
			btcNotify = true
		case "ltc":
			Ltc = chains.NewBtcChain("ltc", 5003)
			ChainDict["ltc"] = &Ltc.Chain
			go func() { Ltc.Listen() }()
			btcNotify = true
		case "doge":
			Doge = chains.NewBtcChain("doge", 5003)
			ChainDict["doge"] = &Doge.Chain
			go func() { Doge.Listen() }()
			btcNotify = true
		case "eth":
			Eth = chains.NewEthChain("eth", 5001)
			ChainDict["eth"] = &Eth.Chain
			InitERC20Eth()
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

// Initialize all ERC20 tokens on Ethereum chain.
// Only run when Ethereum is selected by the user.
func InitERC20Eth() {
	Dai = tokens.NewERC20Eth("dai", "0x6B175474E89094C44Da98b954EedeAC495271d0F")
	ERC20Dict["dai"] = Dai
	Usdt = tokens.NewERC20Eth("usdt", "0xB8c77482e45F1F44dE1745F52C74426C631bDD52")
	ERC20Dict["usdt"] = Usdt
	Usdc = tokens.NewERC20Eth("usdc", "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	ERC20Dict["usdc"] = Usdc
	Ust = tokens.NewERC20Eth("ust", "0xa47c8bf37f92aBed4A126BDA807A7b7498661acD")
	ERC20Dict["ust"] = Ust
}
