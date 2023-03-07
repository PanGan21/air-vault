package pkg

import (
	"context"
	"fmt"
	"log"

	"github.com/PanGan21/air-vault/config"
	airVaultContractInterface "github.com/PanGan21/air-vault/contracts/interfaces/air-vault"
	airdropRepo "github.com/PanGan21/air-vault/repository/airdrop"
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
	airdropRepo             airdropRepo.AirdropDataRepository
}

func NewAppRunner(servicePrivateKey string, airVaultContractAddress string, minter WinMinter, airdropRepo airdropRepo.AirdropDataRepository) App {
	return &app{
		servicePrivateKey:       servicePrivateKey,
		airVaultContractAddress: airVaultContractAddress,
		minter:                  minter,
		airdropRepo:             airdropRepo,
	}
}

func (ap *app) Run(ctx context.Context, client *ethclient.Client) error {
	airVaultContractAddr := common.HexToAddress(ap.airVaultContractAddress)
	airVaultContract, err := airVaultContractInterface.NewContracts(airVaultContractAddr, client)
	if err != nil {
		return err
	}

	// Watch for events
	watchOpts := &bind.WatchOpts{Context: ctx, Start: nil}

	// Setup a channel for deposit results
	depositEventChannel := make(chan *airVaultContractInterface.ContractsDeposited)
	accountRule := []common.Address{}
	subDeposit, err := airVaultContract.WatchDeposited(watchOpts, depositEventChannel, accountRule)
	if err != nil {
		log.Fatal(err)
	}
	defer subDeposit.Unsubscribe()

	// Setup a channel for withdraw results
	withdrawEventChannel := make(chan *airVaultContractInterface.ContractsWithdrawn)
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

			err := ap.handleMessage(airVaultContract, client, int64(event.Raw.BlockNumber), event.Account.Hex())
			if err != nil {
				log.Println(err)
			}
		case event := <-withdrawEventChannel:
			fmt.Printf("Withdraw event received: account=%v, amount=%v\n", event.Account, event.Amount)
			err := ap.handleMessage(airVaultContract, client, int64(event.Raw.BlockNumber), event.Account.Hex())
			if err != nil {
				log.Println(err)
			}
		case err := <-subDeposit.Err():
			log.Fatal(err)
		case err := <-subWithdraw.Err():
			log.Fatal(err)
		}
	}
}

func (ap *app) handleMessage(airVaultContract *airVaultContractInterface.Contracts, client *ethclient.Client, blockNumber int64, account string) error {
	lastAirdropBlockNumber := ap.airdropRepo.GetLastAirdropBlockNumber()
	currentBlockNumber := int64(blockNumber)

	userAddress := common.HexToAddress(account)
	balance, err := airVaultContract.LockedBalanceOf(&bind.CallOpts{Pending: false, Context: context.Background()}, userAddress)
	if err != nil {
		return err
	}

	if currentBlockNumber-lastAirdropBlockNumber >= config.App.Contract.BlocksInterval {
		err := ap.triggerAirdrop(client, currentBlockNumber)
		if err != nil {
			return err
		}
	}
	ap.airdropRepo.UpdateCurrentBlockNumber(currentBlockNumber)

	airdropData, err := ap.airdropRepo.CreateChunk(account, currentBlockNumber, balance)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("New airdrop total FUD balance for user: %s is: %v\n", account, airdropData.TokenToBlocks)
	return nil
}

func (ap *app) triggerAirdrop(client *ethclient.Client, currentBlockNumber int64) error {
	fmt.Println("Triggering airdrop")
	depositors := ap.airdropRepo.GetDepositors()
	for _, depositor := range depositors {
		amount, err := ap.airdropRepo.CalculateWinTokens(depositor, config.App.Contract.MintProportion, config.App.Contract.BlocksInterval, currentBlockNumber)
		if err != nil {
			log.Println("err", err)
			continue
		}
		err = ap.minter.MintWin(context.Background(), client, depositor, amount)
		if err != nil {
			log.Println("err", err)
			continue
		}

		err = ap.airdropRepo.CleanAirdropDataByUser(depositor, currentBlockNumber, config.App.Contract.BlocksInterval)
		if err != nil {
			log.Println("err", err)
			continue
		}
	}
	err := ap.airdropRepo.UpdateLastAirdropNumber(currentBlockNumber)
	if err != nil {
		return err
	}
	return nil
}
