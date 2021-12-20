package erc20eth

import (
	"math/big"
	"time"
	"github.com/colinear-labs/colinear-node/processing"
	"github.com/colinear-labs/colinear-node/processing/erc20eth/erc20abi"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ERC20EthProcessor struct {
	Id              string
	ContractAddress string
	Client          *ethclient.Client
	Contract        *erc20abi.ERC20Caller
}

func NewERC20EthProcessor(id string, contractAddress string, client *ethclient.Client) *ERC20EthProcessor {
	contract, err := erc20abi.NewERC20Caller(common.HexToAddress(contractAddress), client)
	if err != nil {
		panic(err)
	}
	return &ERC20EthProcessor{
		Id:              id,
		ContractAddress: contractAddress,
		Client:          client,
		Contract:        contract,
	}
}

func (p *ERC20EthProcessor) CurrencyId() string {
	return p.Id
}

func (p *ERC20EthProcessor) Process(intent *processing.PaymentIntentLocal) chan processing.PaymentStatus {
	status := make(chan processing.PaymentStatus)

	amtInt64, _ := intent.Amount.Int64()
	amtEth := big.NewInt(amtInt64)
	toEth := common.HexToAddress(intent.To)

	go func() {

		// Pending transaction loop
	pendingLoop:
		for {
			balance, err := p.Contract.BalanceOf(&bind.CallOpts{
				Pending: true,
			}, toEth)
			if err != nil {
				panic(err)
			}
			comparison := balance.Cmp(amtEth)
			if comparison == 1 || comparison == 0 {
				status <- processing.Pending
				break pendingLoop
			}
			time.Sleep(1 * time.Second)
		}

		// Verified transaction loop
	verifiedLoop:
		for {
			balance, err := p.Contract.BalanceOf(&bind.CallOpts{
				Pending: false,
			}, toEth)
			if err != nil {
				panic(err)
			}
			comparison := balance.Cmp(amtEth)
			if comparison == 1 || comparison == 0 {
				status <- processing.Verified
				break verifiedLoop
			}
			time.Sleep(1 * time.Second)
		}

	}()

	return status
}
