package command

import (
	"context"
	"errors"

	"github.com/PanGan21/air-vault/config"
	"github.com/spf13/cobra"
)

func NewRootCommand(ctx context.Context) *cobra.Command {
	rootCommand := &cobra.Command{
		Use:               "airvault",
		Short:             "Run the air-vault service",
		PersistentPreRunE: config.Setup,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.New("specify a command: [deploy, run]")
		},
	}

	rootCommand.Flags().StringVarP(
		&config.Filename,
		"config",
		"f",
		"config/local.yaml",
		"Relative path to the config file",
	)

	rootCommand.PersistentFlags().StringP("blockchain.pk", "k", "", "Account private key")
	rootCommand.AddCommand(NewDeployCommand(ctx))
	rootCommand.AddCommand(NewRunnerCommand(ctx))
	rootCommand.AddCommand(NewDepositCommand(ctx))
	rootCommand.AddCommand(NewTransferFudCommand(ctx))

	return rootCommand
}
