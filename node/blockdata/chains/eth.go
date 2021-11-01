// For all ETH-like chains
// NOTE: scans for erc20 transactions as well

package chains

import (
	"context"
	"fmt"
	"math/big"
	"xnode/blockdata/basechain"
	"xnode/blockdata/tokens"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type EthChain struct {
	Chain  basechain.BaseChain
	Client *ethclient.Client
}

// Ethereum header struct. SUBJECT TO CHANGE
type EthHeader struct {
	HashPrevBlock  [32]byte
	HashMerkleRoot [32]byte
	Time           uint32
	Bits           uint32
	Nonce          uint32
}

// Instantiate new ETH chain.
func NewEthChain(name string, port uint) *EthChain {
	chain := *basechain.NewChain(name, port)
	client, err := ethclient.Dial(fmt.Sprintf("ws://127.0.0.1:%s/wsrpc", fmt.Sprint(chain.Port)))
	if err != nil {
		fmt.Printf("Failed to dial full node on port %s.\n", fmt.Sprint(chain.Port))
		panic(err)
	}
	c := EthChain{Chain: chain, Client: client}
	return &c
}

// Listen for new blocks (over websocket to local node)
func (chain *EthChain) Listen() {

	headers := make(chan *types.Header)

	sub, err := chain.Client.SubscribeNewHead(context.Background(), headers)

	if err != nil {
		panic(err)
	}

	for {
		select {
		case err := <-sub.Err():
			panic(err)
		case header := <-headers:
			block, err := chain.Client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				panic(err)
			}

			newHeader := EthHeader{
				HashPrevBlock:  header.ParentHash,
				HashMerkleRoot: header.Hash(),
				Time:           uint32(header.Time),
				Bits:           uint32(header.Difficulty.Uint64()),
				Nonce:          uint32(header.Nonce.Uint64()),
			}

			txs := []basechain.Tx{}

			for _, transaction := range block.Transactions() {
				amt := new(big.Float).SetInt(transaction.Value())
				amt = amt.Mul(amt, big.NewFloat(1e-18))

				tx := basechain.Tx{
					Txid:   transaction.Hash().String(),
					To:     transaction.To().String(),
					Amount: amt,
				}
				txs = append(txs, tx)
			}

			chain.Chain.NewHeader(newHeader)
			chain.Chain.NewBlock(basechain.Block{Txs: txs})
			tokens.NewEthBlockEvents <- uint(header.Number.Uint64())
		}
	}

}
