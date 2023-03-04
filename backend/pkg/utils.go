package pkg

import (
	"context"
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/PanGan21/air-vault/config"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	ErrInvalidKey = errors.New("invalid key")
)

func getSigner(ctx context.Context, client *ethclient.Client, privKey string) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(privKey)
	if err != nil {
		return nil, err
	}
	publicKey, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return nil, ErrInvalidKey
	}
	address := crypto.PubkeyToAddress(*publicKey)
	nonce, err := client.PendingNonceAt(ctx, address)
	if err != nil {
		return nil, err
	}
	chainId, err := client.ChainID(ctx)
	if err != nil {
		return nil, err
	}
	signer, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		return nil, err
	}

	signer.Nonce = big.NewInt(int64(nonce))
	signer.Value = big.NewInt(config.App.Contract.WeiFunds)
	signer.GasLimit = uint64(config.App.Contract.GasLimit)
	signer.GasPrice = big.NewInt(config.App.Contract.GasPrice)

	return signer, nil
}
