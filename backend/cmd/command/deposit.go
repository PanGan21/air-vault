package command

import (
	"context"
	"fmt"
	"math/big"

	"github.com/PanGan21/air-vault/config"
	"github.com/PanGan21/air-vault/pkg"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

func NewDepositCommand(ctx context.Context) *cobra.Command {
	depositCommand := &cobra.Command{
		Use:   "deposit",
		Short: "Deposit FUD Token to AirVault contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deposit(ctx)
		},
	}
	return depositCommand
}

func deposit(ctx context.Context) error {
	fmt.Println("Depositing")

	ctx, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()

	client, err := ethclient.DialContext(ctx, config.App.Blockchain.Address)
	if err != nil {
		return err
	}

	depositer := pkg.NewDepositRunner(config.App.Demo.PrivateKey, config.App.Demo.Address, config.App.Contract.AirVaultAddress, config.App.Contract.FudTokenAddress)
	err = depositer.CreateDeposit(ctx, client, big.NewInt(config.App.Demo.DepositAmount))
	if err != nil {
		return err
	}

	return nil
}
