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

func NewWithdrawCommand(ctx context.Context) *cobra.Command {
	withdrawCommand := &cobra.Command{
		Use:   "withdraw",
		Short: "Withdraw FUD Token from the AirVault contract",
		RunE: func(cmd *cobra.Command, args []string) error {
			return withdraw(ctx)
		},
	}
	return withdrawCommand
}

func withdraw(ctx context.Context) error {
	fmt.Println("Withrawing")

	ctx, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()

	client, err := ethclient.DialContext(ctx, config.App.Blockchain.Address)
	if err != nil {
		return err
	}

	withrawer := pkg.NewWithdrawRunner(config.App.Demo.PrivateKey, config.App.Demo.Address, config.App.Contract.AirVaultAddress, config.App.Contract.FudTokenAddress)
	err = withrawer.CreateWithdraw(ctx, client, big.NewInt(config.App.Demo.WithdrawAmount))
	if err != nil {
		return err
	}

	return nil
}
