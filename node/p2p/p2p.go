// Accept p2p connections & handle requests

package p2p

import (
	"fmt"
	"xnode/processing"
	"xnode/runtime"
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
		obj, err := ctx.DecodeMessage()
		if err != nil {
			return nil
		}
		paymentIntent, ok := obj.(xutil.PaymentIntent)
		if ok {
			targetProcessor, ok2 := runtime.Processors[paymentIntent.Currency]
			if !ok2 {
				ctx.Send([]byte(fmt.Sprintf("Currency %s is not supported.", paymentIntent.Currency)))
				return nil
			}

			localIntent := processing.NewPaymentIntentLocal(
				paymentIntent.Currency,
				paymentIntent.Amount,
				paymentIntent.To,
			)

			go func() {
				targetProcessor.Process(localIntent)
			}()
		}

		return nil
	})

	node.Listen()

	return node
}
