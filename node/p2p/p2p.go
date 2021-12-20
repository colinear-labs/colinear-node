// Accept p2p connections & handle requests

package p2p

import (
	"context"
	"fmt"
	"time"
	"github.com/colinear-labs/colinear-node/processing"
	"github.com/colinear-labs/colinear-node/runtime"
	"github.com/colinear-labs/colinear-node/xutil"
	"github.com/colinear-labs/colinear-node/xutil/currencies"
	"github.com/colinear-labs/colinear-node/xutil/ipassign"
	"github.com/colinear-labs/colinear-node/xutil/p2pshared"

	"github.com/perlin-network/noise"
	"github.com/perlin-network/noise/kademlia"
)

var Node *noise.Node = nil

func InitP2P() {

	// logger, _ := zap.NewDevelopment()

	port := 9871
	broadcastIp := ipassign.GetIPv6Address()

	Node, err := noise.NewNode(
		noise.WithNodeBindPort((uint16)(port)),
		noise.WithNodeAddress(fmt.Sprintf("[%s]:%d", broadcastIp, port)),
		// noise.WithNodeLogger(logger),
	)
	if err != nil {
		panic(err)
	}
	defer Node.Close()

	k := kademlia.New()
	Node.Bind(k.Protocol())

	xutil.RegisterNodeMessages(Node)

	Node.Handle(func(ctx noise.HandlerContext) error {

		if ctx.IsRequest() {
			return nil
		}

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
				fmt.Println(runtime.Processors)
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
				fmt.Printf("Now processing intent %s\n", paymentIntent)
				targetProcessor.Process(localIntent)
			}()
		}

		return nil
	})

	if err := Node.Listen(); err != nil {
		panic(err)
	}

	// Ping bootstrap nodes to make the network aware of presence

	// NOTE: idk if looping this is necessary for full nodes lol
	// it could be that we just need to ping once then we're done

	go func() {

		for {

			timeoutCtx, _ := context.WithTimeout(context.Background(), time.Duration(5*time.Second))
			for _, addr := range p2pshared.BootstrapNodes {
				if _, err := Node.Ping(timeoutCtx, addr+":9871"); err != nil {
					fmt.Printf("Failed to ping bootstrap node at %s.\n", addr)
				} else {
					fmt.Printf("Pinged bootstrap node at %s.\n", addr)
				}
			}

			Peers := k.Discover()
			fmt.Printf("Peers: %s\n", fmt.Sprint(Peers))

			// wait 10 mins before contacting bootstrap nodes again

			time.Sleep(10 * time.Minute)
			// time.Sleep(5 * time.Second)

		}

	}()
}
