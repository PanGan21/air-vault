package pkg

import (
	"context"
	"fmt"
	"math/big"

	"github.com/PanGan21/air-vault/config"
	fudTokenContractInterface "github.com/PanGan21/air-vault/contracts/interfaces/fud-token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type FudTransfer interface {
	TransferFud(ctx context.Context, client *ethclient.Client, amount *big.Int) error
}

type fudTransfer struct {
	privateKey              string
	fudTokenContractAddress string
	address                 string
}

func NewFudTransferRunner(privateKey string, fudTokenContractAddress string, receiverAddress string) FudTransfer {
	return &fudTransfer{
		privateKey:              privateKey,
		fudTokenContractAddress: fudTokenContractAddress,
		address:                 receiverAddress,
	}
}

func (f *fudTransfer) TransferFud(ctx context.Context, client *ethclient.Client, amount *big.Int) error {
	contract, err := fudTokenContractInterface.NewContracts(common.HexToAddress(f.fudTokenContractAddress), client)
	if err != nil {
		return err
	}

	signer, err := getSigner(ctx, client, config.App.Blockchain.PrivateKey)
	if err != nil {
		return err
	}

	address := common.HexToAddress(f.address)
	tx, err := contract.Transfer(signer, address, amount)
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

	balance, err := contract.BalanceOf(&bind.CallOpts{Pending: false, Context: ctx}, address)
	if err != nil {
		return err
	}

	fmt.Printf("User with address: %s has FUD Token balance: %s\n", address, balance)

	return nil
}
