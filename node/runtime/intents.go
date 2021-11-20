package runtime

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"xnode/processing"
	"xnode/processing/btc"
	"xnode/processing/erc20eth"
	"xnode/processing/eth"
)

var Processors = make(map[string]processing.Processor)

// Initialize PaymentIntent processing + define processors in dict
func InitProcessors(selectedCurrencies []string) {

	for _, currency := range selectedCurrencies {
		switch currency {
		case "btc", "ltc", "bch", "doge":
			// later, consider changing historyLen per-chain-- but leave it at 10 for now.
			Processors[currency] = btc.NewBtcProcessor(currency, processing.NodePorts[currency], 10)
			if !btcBlockNotifyServerRunning {
				go func() {
					SpawnBtcBlockNotifyServer()
					btcBlockNotifyServerRunning = true
				}()
			}
		case "eth":
			ethProcessor := eth.NewEthProcessor(currency, processing.NodePorts[currency])
			Processors[currency] = ethProcessor
			for _, e := range []string{
				"dai",
				"usdt",
				"usdc",
				"ust",
				// "ampl",
			} {
				Processors[currency] = erc20eth.NewERC20EthProcessor(
					e,
					processing.TokenAddresses[currency],
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		}()

	})

	http.ListenAndServe("127.0.0.1:4999", nil)

}
