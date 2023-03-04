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

func NewTransferFudCommand(ctx context.Context) *cobra.Command {
	transferFudCommand := &cobra.Command{
		Use:   "transfer-fud",
		Short: "Transfer FUD Tokens to the demo user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return transferFud(ctx)
		},
	}
	return transferFudCommand
}

func transferFud(ctx context.Context) error {
	fmt.Println("Transferring")

	ctx, cancel := context.WithTimeout(ctx, config.App.Blockchain.TimeoutIn)
	defer cancel()

	client, err := ethclient.DialContext(ctx, config.App.Blockchain.Address)
	if err != nil {
		return err
	}

	transferrer := pkg.NewFudTransferRunner(config.App.Blockchain.PrivateKey, config.App.Contract.FudTokenAddress, config.App.Demo.Address)
	// TODO: Move 10000 into config; do not hardcode it
	err = transferrer.TransferFud(ctx, client, big.NewInt(10000))
	if err != nil {
		return err
	}

	return nil
}
