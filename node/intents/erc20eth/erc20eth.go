package erc20eth

import (
	"math/big"
	"time"
	"xnode/intents"
	"xnode/intents/erc20eth/erc20abi"

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

func (p *ERC20EthProcessor) Process(intent *intents.PaymentIntent) chan intents.PaymentStatus {
	status := make(chan intents.PaymentStatus)

	amtInt64, _ := intent.Amount.Int64()
	amtEth := big.NewInt(amtInt64)
	toEth := common.HexToAddress(intent.To)

	// Pending transaction loop
	// NOTE: FIND OUT HOW TO GET PENDING BALANCES FROM ERC20ABI BINDINGS
	go func() {
	}()

	// Verified transaction loop
	go func() {
	verifiedLoop:
		for {
			balance, _ := p.Contract.BalanceOf(nil, toEth)
			comparison := balance.Cmp(amtEth)
			if comparison == 1 || comparison == 0 {
				status <- intents.Verified
				break verifiedLoop
			}
			time.Sleep(1 * time.Second)
		}
	}()

	return status
}
