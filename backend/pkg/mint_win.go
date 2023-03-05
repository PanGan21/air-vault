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

// WIN tokens in airdrop = proportion * (# FUD tokens deposited) * (# blocks deposited) / ( total # blocks)
func calculateWinTokens(depositAmount *big.Int, numDepositedBlocks uint64, totalNumBlocks uint64) *big.Int {
	// calculate the average deposit amount per block
	averageDepositPerBlock := new(big.Int).Div(depositAmount, big.NewInt(int64(numDepositedBlocks)))

	// calculate the percentage of the total block interval covered by the deposit interval
	blockIntervalRatio := new(big.Float).Quo(new(big.Float).SetUint64(numDepositedBlocks), new(big.Float).SetUint64(totalNumBlocks))

	// calculate the amount of WIN tokens to mint (proportion * average deposit per block * block interval ratio)
	mintProportion := float64(config.App.Contract.MintProportion / 100)
	winTokenAmount := new(big.Float).Mul(new(big.Float).SetInt(averageDepositPerBlock), new(big.Float).SetFloat64(mintProportion))
	winTokenAmount = winTokenAmount.Mul(winTokenAmount, blockIntervalRatio)

	// convert the result to a big.Int and round down
	winTokenAmountInt := new(big.Int)
	winTokenAmount.Int(winTokenAmountInt)

	return winTokenAmountInt
}
