// Accept p2p connections & handle requests

package p2p

import (
	"fmt"

	"xnode/xutil"

	"github.com/perlin-network/noise"
)

func InitServer() *noise.Node {

	node, err := noise.NewNode()
	if err != nil {
		panic(err)
	}
	defer node.Close()

	node.RegisterMessage(xutil.PaymentIntent{}, xutil.UnmarshalPaymentIntent)
	node.RegisterMessage(xutil.PaymentResponse{}, xutil.UnmarshalPaymentResponse)

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
