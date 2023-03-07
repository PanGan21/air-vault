package command

import (
	"context"
	"fmt"

	"github.com/PanGan21/air-vault/config"
	"github.com/PanGan21/air-vault/pkg"
	airdropRepository "github.com/PanGan21/air-vault/repository/airdrop"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
)

func NewRunnerCommand(ctx context.Context) *cobra.Command {
	runCommand := &cobra.Command{
		Use:   "run",
		Short: "Run the backend",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(ctx)
		},
	}

	return runCommand
}

func run(ctx context.Context) error {
	fmt.Println("Start backend")

	ctx, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()

	client, err := ethclient.DialContext(ctx, config.App.Blockchain.Ws)
	if err != nil {
		return err
	}

	// Perform dependency injection so implementations can be independent
	airdropRepo := airdropRepository.NewAirdropRepository()
	minter := pkg.NewWinMinterRunner(config.App.Blockchain.PrivateKey, config.App.Contract.WinTokenAddress)
	appRunner := pkg.NewAppRunner(config.App.Blockchain.PrivateKey, config.App.Contract.AirVaultAddress, minter, airdropRepo)
	err = appRunner.Run(ctx, client)
	if err != nil {
		return err
	}

	return nil
}
