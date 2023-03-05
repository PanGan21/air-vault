package pkg

import (
	"context"
	"fmt"
	"log"

	airVaultContractInterface "github.com/PanGan21/air-vault/contracts/interfaces/air-vault"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type App interface {
	Run(ctx context.Context, client *ethclient.Client) error
}

type app struct {
	servicePrivateKey       string
	airVaultContractAddress string
	minter                  WinMinter
}

func NewAppRunner(servicePrivateKey string, airVaultContractAddress string, minter WinMinter) App {
	return &app{
		servicePrivateKey:       servicePrivateKey,
		airVaultContractAddress: airVaultContractAddress,
		minter:                  minter,
	}
}

func (ap *app) Run(ctx context.Context, client *ethclient.Client) error {
	airVaultContractAddr := common.HexToAddress(ap.airVaultContractAddress)
	airVaultContract, err := airVaultContractInterface.NewContracts(airVaultContractAddr, client)
	if err != nil {
		return err
	}

	// Watch for a Deposited event
	watchOpts := &bind.WatchOpts{Context: ctx, Start: nil}
	// Setup a channel for results
	depositEventChannel := make(chan *airVaultContractInterface.ContractsDeposited)
	// Start a goroutine which watches new events
	accountRule := []common.Address{}
	subDeposit, err := airVaultContract.WatchDeposited(watchOpts, depositEventChannel, accountRule)
	if err != nil {
		log.Fatal(err)
	}
	defer subDeposit.Unsubscribe()

	// Setup a channel for results
	withdrawEventChannel := make(chan *airVaultContractInterface.ContractsWithdrawn)
	// Start a goroutine which watches new events
	subWithdraw, err := airVaultContract.WatchWithdrawn(watchOpts, withdrawEventChannel, accountRule)
	if err != nil {
		log.Fatal(err)
	}
	defer subWithdraw.Unsubscribe()

	// Loop through the event channel and print out any deposit events received
	for {
		select {
		case event := <-depositEventChannel:
			fmt.Printf("Deposited event received: account=%v, amount=%v\n", event.Account, event.Amount)
			// err := ap.minter.MintWin(ctx, client, event.Account.Hex(), event.Amount)
			if err != nil {
				log.Fatal(err)
			}
		case event := <-withdrawEventChannel:
			fmt.Printf("Withdraw event received: account=%v, amount=%v\n", event.Account, event.Amount)
			// err := ap.minter.MintWin(ctx, client, event.Account.Hex(), event.Amount)
			if err != nil {
				log.Fatal(err)
			}
		case err := <-subDeposit.Err():
			log.Fatal(err)
		}
	}
}
