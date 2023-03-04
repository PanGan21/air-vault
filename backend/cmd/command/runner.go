package command

import (
	"context"
	"fmt"

	"github.com/PanGan21/air-vault/config"
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

	_, err := ethclient.DialContext(ctx, config.App.Blockchain.Address)
	if err != nil {
		return err
	}

	return nil
}
