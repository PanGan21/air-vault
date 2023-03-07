package pkg

import (
	"context"
	"fmt"
	"math/big"

	"github.com/PanGan21/air-vault/config"
	winTokenContractInterface "github.com/PanGan21/air-vault/contracts/interfaces/win-token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type WinMinter interface {
	MintWin(ctx context.Context, client *ethclient.Client, toAddress string, amount *big.Int) error
}

type winMinter struct {
	privateKey              string
	winTokenContractAddress string
}

func NewWinMinterRunner(privateKey string, winTokenContractAddress string) WinMinter {
	return &winMinter{
		privateKey:              privateKey,
		winTokenContractAddress: winTokenContractAddress,
	}
}

func (w *winMinter) MintWin(ctx context.Context, client *ethclient.Client, toAddress string, amount *big.Int) error {
	contract, err := winTokenContractInterface.NewContracts(common.HexToAddress(w.winTokenContractAddress), client)
	if err != nil {
		return err
	}

	signer, err := getSigner(ctx, client, config.App.Blockchain.PrivateKey)
	if err != nil {
		return err
	}

	fmt.Printf("Minting %d WIN for account with address %s\n", amount, toAddress)

	toAddr := common.HexToAddress(toAddress)
	tx, err := contract.Mint(signer, toAddr, amount)
	if err != nil {
		return err
	}

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return err
	}

	if receipt.Status != types.ReceiptStatusSuccessful || err != nil {
		return err
	}

	balance, err := contract.BalanceOf(&bind.CallOpts{Pending: false, Context: ctx}, toAddr)
	if err != nil {
		return err
	}

	fmt.Printf("User with address: %s has WIN Token balance: %s\n", toAddr, balance)

	return nil
}
