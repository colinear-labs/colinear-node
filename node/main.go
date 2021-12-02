package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
	"go.uber.org/zap"
)

func main() {

	logger, _ := zap.NewDevelopment(zap.AddStacktrace(zap.PanicLevel))
	defer logger.Sync()

	node, _ := noise.NewNode(noise.WithNodeLogger(logger), noise.WithNodeBindPort(9871))
	defer node.Close()

	overlay := kademlia.New()
	node.Bind(overlay.Protocol())

	node.Handle(func(ctx noise.HandlerContext) error {
		fmt.Printf("Got a message from %s: '%s'\n", ctx.ID().Host, string(ctx.Data()))

		return nil
	})

	if err := node.Listen(); err != nil {
		panic(err)
	}

	address := "127.0.0.1:9871"

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	_, err := node.Ping(ctx, address) // where the magic happens
	cancel()

	if err != nil {
		fmt.Printf("Failed to ping node (%s). Skipping... [error: %s]\n", address, err)
	}

	node.Send(ctx, address, []byte("hello"))
	cancel()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
