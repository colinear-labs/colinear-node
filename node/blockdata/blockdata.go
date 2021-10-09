package blockdata

import (
	"xnode/blockdata/btc"
	"xnode/blockdata/eth"
)

var BtcData = btc.NewBtcChain("btc")
var Bch = btc.NewBtcChain("bch")
var Doge = btc.NewBtcChain("doge")
var Ltc = btc.NewBtcChain("ltc")
var Eth = eth.NewEthChain("eth")
