package pkg

import (
	"context"
	"fmt"
	"math/big"

	airVaultContractInterface "github.com/PanGan21/air-vault/contracts/interfaces/air-vault"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Withdraw interface {
	CreateWithdraw(ctx context.Context, client *ethclient.Client, amount *big.Int) error
}

type withdraw struct {
	userPrivateKey          string
	userAddress             string
	airVaultContractAddress string
	fudTokenContractAddress string
}

func NewWithdrawRunner(userPrivateKey string, userAddress string, airVaultContractAddress string, fudTokenContractAddress string) Withdraw {
	return &withdraw{
		userPrivateKey:          userPrivateKey,
		userAddress:             userAddress,
		airVaultContractAddress: airVaultContractAddress,
		fudTokenContractAddress: fudTokenContractAddress,
	}
}

func (d *withdraw) CreateWithdraw(ctx context.Context, client *ethclient.Client, amount *big.Int) error {
	airVaultContract, err := airVaultContractInterface.NewContracts(common.HexToAddress(d.airVaultContractAddress), client)
	if err != nil {
		return err
	}

	signer, err := getSigner(ctx, client, d.userPrivateKey)
	if err != nil {
		return err
	}

	tx, err := airVaultContract.Withdraw(signer, amount)
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

	userAddr := common.HexToAddress(d.userAddress)
	lockedBalance, err := airVaultContract.LockedBalanceOf(&bind.CallOpts{Pending: false, Context: ctx}, userAddr)
	if err != nil {
		return err
	}

	fmt.Printf("Locked balance of user with address: %s is: %s\n", d.userAddress, lockedBalance)

	return nil
}
