package pkg

import (
	"context"
	"math/big"

	"github.com/PanGan21/air-vault/config"
	airVaultTokenContract "github.com/PanGan21/air-vault/contracts/interfaces/air-vault"
	fudTokenContract "github.com/PanGan21/air-vault/contracts/interfaces/fud-token"
	winTokenContract "github.com/PanGan21/air-vault/contracts/interfaces/win-token"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// Deployer interface
type Deployer interface {
	Deploy(ctx context.Context, client *ethclient.Client, fudTokenSupply *big.Int) error
	AirVaultContractAddress() string
	FudTokenContractAddress() string
	WinTokenContractAddress() string
}

type deployer struct {
	airVaultContractAddress common.Address
	fudTokenContractAddress common.Address
	winTokenContractAddress common.Address
}

// NewDepoyer returns a new Deployer instance
func NewDeployer() Deployer {
	return &deployer{}
}

func (d *deployer) Deploy(ctx context.Context, client *ethclient.Client, fudTokenSupply *big.Int) error {
	fudTokenSigner, err := getSigner(ctx, client, config.App.Blockchain.PrivateKey)
	if err != nil {
		return err
	}

	fudTokenContractAddress, tx, _, err := fudTokenContract.DeployContracts(fudTokenSigner, client, fudTokenSupply)
	if err != nil {
		return err
	}
	_, err = bind.WaitDeployed(ctx, client, tx)
	if err != nil {
		return err
	}

	d.fudTokenContractAddress = fudTokenContractAddress

	winTokenSigner, err := getSigner(ctx, client, config.App.Blockchain.PrivateKey)
	if err != nil {
		return err
	}

	winTokenContractAddressAddress, tx, _, err := winTokenContract.DeployContracts(winTokenSigner, client)
	if err != nil {
		return err
	}
	_, err = bind.WaitDeployed(ctx, client, tx)
	if err != nil {
		return err
	}

	d.winTokenContractAddress = winTokenContractAddressAddress

	airVaultSigner, err := getSigner(ctx, client, config.App.Blockchain.PrivateKey)
	if err != nil {
		return err
	}

	airVaultContractAddress, tx, _, err := airVaultTokenContract.DeployContracts(airVaultSigner, client, fudTokenContractAddress, winTokenContractAddressAddress)
	if err != nil {
		return err
	}
	_, err = bind.WaitDeployed(ctx, client, tx)
	if err != nil {
		return err
	}

	d.airVaultContractAddress = airVaultContractAddress

	return nil
}

func (d *deployer) AirVaultContractAddress() string {
	return d.airVaultContractAddress.Hex()
}

func (d *deployer) FudTokenContractAddress() string {
	return d.fudTokenContractAddress.Hex()
}

func (d *deployer) WinTokenContractAddress() string {
	return d.winTokenContractAddress.Hex()
}
