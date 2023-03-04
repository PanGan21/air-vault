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

func NewDeployCommand(ctx context.Context) *cobra.Command {
	deployCommand := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy FUD Token, WIN Token, AirVault contracts to blockchain",
		RunE: func(cmd *cobra.Command, args []string) error {
			return deploy(ctx)
		},
	}
	return deployCommand
}

func deploy(ctx context.Context) error {
	fmt.Println("Deploying contract")

	ctx, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()

	client, err := ethclient.DialContext(ctx, config.App.Blockchain.Address)
	if err != nil {
		return err
	}

	deployer := pkg.NewDeployer()
	err = deployer.Deploy(ctx, client, big.NewInt(config.App.Contract.FudTokenSupply))
	if err != nil {
		return err
	}

	fmt.Printf("Contracts deployed.\nFUD Token: %s\nWIN Token: %s\nAirVault: %s\n", deployer.FudTokenContractAddress(), deployer.WinTokenContractAddress(), deployer.AirVaultContractAddress())

	return nil
}
