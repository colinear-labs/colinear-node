package runtime

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/colinear-labs/colinear-node/processing"
	"github.com/colinear-labs/colinear-node/processing/btc"
	"github.com/colinear-labs/colinear-node/processing/erc20eth"
	"github.com/colinear-labs/colinear-node/processing/eth"
	"github.com/colinear-labs/colinear-node/xutil/currencies"
)

var Processors = make(map[string]processing.Processor)

// Initialize PaymentIntent processing + define processors in dict
func InitProcessors(selectedChains []string) {

	for _, currency := range selectedChains {
		switch currency {
		case "btc", "ltc", "bch", "doge":
			// later, consider changing historyLen per-chain-- but leave it at 10 for now.
			Processors[currency] = btc.NewBtcProcessor(currency, processing.NodePorts[currency], 10)
			if !btcBlockNotifyServerRunning {
				go func() {
					SpawnBtcBlockNotifyServer()
				}()
				btcBlockNotifyServerRunning = true
			}
		case "eth":
			ethProcessor := eth.NewEthProcessor(currency, processing.NodePorts[currency])
			currencies.SupportAllEthTokens()
			Processors[currency] = ethProcessor
			for _, e := range currencies.EthTokens {
				Processors[e] = erc20eth.NewERC20EthProcessor(
					e,
					// processing.TokenAddresses[currency],
					currencies.CurrencyData[currency].TokenAddress,
					ethProcessor.Client,
				)
			}
		}
	}
}

var btcBlockNotifyServerRunning bool = false

// Spawn a local HTTP server to listen for new block events from
// BTC blockchains. Triggered by cURL requests from node containers,
// & pushes new events to the NewBlockEvents channel
func SpawnBtcBlockNotifyServer() {

	mux := http.NewServeMux() // see: https://varunksaini.com/posts/go-http-multiple-registration-error/
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Printf("HTTP ERROR %s", err)
		}

		str := (string)(body)
		u, _ := url.ParseQuery(str)
		name := u["name"][0]
		hash := u["hash"][0]

		go func() {
			// Only works for BtcProcessors ((((theoretically))))
			Processors[name].(*btc.BtcProcessor).NewBlockRpcEvents <- hash
			fmt.Printf("Pushed new %s block hash %s to processor.\n", name, hash)
		}()

	})

	http.ListenAndServe("127.0.0.1:4999", nil)

}
