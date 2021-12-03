// Accept p2p connections & handle requests

package p2p

import (
	"fmt"
	"xnode/processing"
	"xnode/runtime"
	"xnode/xutil"
	"xnode/xutil/currencies"

	"github.com/perlin-network/noise"
)

var Node *noise.Node = nil

func InitServer() {

	Node, err := noise.NewNode(noise.WithNodeBindPort(9871))
	if err != nil {
		panic(err)
	}
	defer Node.Close()

	Node.RegisterMessage(xutil.PaymentIntent{}, xutil.UnmarshalPaymentIntent)
	Node.RegisterMessage(xutil.PaymentResponse{}, xutil.UnmarshalPaymentResponse)

	Node.Handle(func(ctx noise.HandlerContext) error {
		obj, err := ctx.DecodeMessage()

		if err != nil {
			// try []byte encoding instead

			if (string)(ctx.Data()) == "peerinfo" {
				ctx.SendMessage(xutil.PeerInfo{
					Type:       xutil.Node,
					Currencies: currencies.Chains,
				})
			}
			return nil
		}

		paymentIntent, ok := obj.(xutil.PaymentIntent)
		if ok {
			targetProcessor, ok2 := runtime.Processors[paymentIntent.Currency]
			if !ok2 {
				fmt.Printf("Request failed: Currency %s is not supported.\n", paymentIntent.Currency)
				ctx.SendMessage(xutil.PaymentResponse{
					To:     paymentIntent.To,
					Status: xutil.IntentError,
				})
				return nil
			}

			localIntent := processing.NewPaymentIntentLocal(
				paymentIntent.Currency,
				paymentIntent.Amount,
				paymentIntent.To,
			)

			go func() {
				fmt.Printf("Now processing intent %s", paymentIntent)
				targetProcessor.Process(localIntent)
			}()
		}

		return nil
	})

	Node.Listen()
}
