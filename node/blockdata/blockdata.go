package blockdata

import (
	"xnode/blockdata/btc"
	"xnode/blockdata/eth"
	"xnode/blockdata/xyz"
)

var Btc = btc.NewBtcChain("btc")
var Bch = btc.NewBtcChain("bch")
var Doge = btc.NewBtcChain("doge")
var Ltc = btc.NewBtcChain("ltc")

var Eth = eth.NewEthChain("eth")

var Xmr = xyz.NewChain("xmr")
var Zec = xyz.NewChain("zec")
