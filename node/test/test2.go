package main

import (
	"xnode/blockdata/chains"
)

func main() {
	chains.SpawnBtcBlockNotifyServer(4999)
	// btc := btc.NewBtcChain("btc")
	// btc.Listen(5003)
}
