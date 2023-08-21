package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func transferETH(client *ethclient.Client, fromPrivKey *ecdsa.PrivateKey, to common.Address, amount *big.Int) error {
	fmt.Println("transfer method called")
	ctx := context.Background()
	publicKey := fromPrivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting pubic key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	balance, err := client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("balance: %d\n", new(big.Int).Div(balance, big.NewInt(1000000000000000000)))
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return err
	}

	gasLimit := uint64(21000) // in units
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)

	chainID := big.NewInt(1337)

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromPrivKey)
	if err != nil {
		return err
	}

	balance, err = client.BalanceAt(ctx, fromAddress, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("balance after transfer: %d\n", new(big.Int).Div(balance, big.NewInt(1000000000000000000)))
	return client.SendTransaction(ctx, signedTx)
}

func (ex *Exchange) getReserve() (*big.Int, error) {
	exchangePubKey := ex.PrivateKey.Public()
	exchangePublicKeyECDSA, ok := exchangePubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("error casting pubic key to ECDSA")
	}
	exchangeAddress := crypto.PubkeyToAddress(*exchangePublicKeyECDSA)	
	balance, err := ex.Client.BalanceAt(context.Background(), exchangeAddress, nil)
	balanceInEth := new(big.Int).Div(balance, big.NewInt(1000000000000000000))
	if err != nil {
		return nil, fmt.Errorf("error getting balance: %v", err)	
	}	
	return balanceInEth, nil
}
