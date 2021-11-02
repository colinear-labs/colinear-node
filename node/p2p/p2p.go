// Accept p2p connections & handle requests

package p2p

import (
	"fmt"
	"xnode/intents"

	"github.com/perlin-network/noise"
)

func InitServer() *noise.Node {

	node, err := noise.NewNode()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	node.RegisterMessage(noise.NewMessage(1, &intents.PaymentIntent{}))

	node.Handle(func(ctx noise.HandlerContext) error {
		if !ctx.IsRequest() {
			return nil
		}
		fmt.Printf("Received request: %s\n", string(ctx.Data()))

		return ctx.Send([]byte("Hello from the server!"))
	})

	node.Listen()

	return node
}
