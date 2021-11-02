package runtime

import (
	"xnode/intents"
	"xnode/intents/btc"
	"xnode/intents/erc20eth"
	"xnode/intents/eth"
)

var Processors = make(map[string]intents.Processor)

// Initialize PaymentIntent processing + define processors in dict
func InitProcessors(selectedCurrencies []string) {

	for _, currency := range selectedCurrencies {
		switch currency {
		case "btc", "ltc", "bch", "doge":
			Processors[currency] = btc.NewBtcProcessor(currency, intents.NodePorts[currency])
		case "eth":
			Processors[currency] = eth.NewEthProcessor(currency, intents.NodePorts[currency])
			for _, e := range []string{
				"dai",
				"usdt",
				"usdc",
				"ust",
			} {
				Processors["dai"] = erc20eth.NewERC20EthProcessor(e, intents.TokenAddresses[currency])
			}
		}
	}
}
