package main

import (
	"github.com/perlin-network/noise"
	"go.uber.org/zap"
)

func main() {

	logger, err := zap.NewDevelopment(zap.AddStacktrace(zap.PanicLevel))
	if err != nil {
		panic(err)
	}

	node, err := noise.NewNode(noise.WithNodeLogger(logger), noise.WithNodeBindPort(9000))
	if err != nil {
		panic(err)
	}

	if err := node.Listen(); err != nil {
		panic(err)
	}

	defer node.Close()
}
