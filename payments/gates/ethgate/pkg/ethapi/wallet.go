package ethapi

import (
	"context"
	"crypto/ecdsa"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/status-im/keycard-go/hexutils"
	"math/big"
)

type Wallet struct {
	Address *common.Address
	Private *ecdsa.PrivateKey
}

func PublicBase64ToAddress(publicBase64 string) (*common.Address, error) {
	publicBytes, err := base64.StdEncoding.DecodeString(publicBase64)
	if err != nil {
		return nil, err
	}
	return PublicKeyBytesToAddress(publicBytes), nil
}

func Base64ToPrivate(privateBase64 string) (*ecdsa.PrivateKey, error) {
	privateBytes, err := base64.StdEncoding.DecodeString(privateBase64)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func GetWallet(private *ecdsa.PrivateKey) (*Wallet, error) {
	address, ok := PrivateToAddress(private)
	if !ok {
		return nil, errors.New("PrivateToAddress error! ")
	}
	return &Wallet{address, private}, nil
}

func (account *Wallet) GetPrivateHEX() string {
	return hexutils.BytesToHex(crypto.FromECDSA(account.Private))
}

func (account *Wallet) GetPrivateBase64() string {
	return base64.StdEncoding.EncodeToString(crypto.FromECDSA(account.Private))
}

func (account *Wallet) GetAddressHEX() string {
	return account.Address.Hex()
}

func (account *Wallet) GetAddressBase64() string {
	return base64.StdEncoding.EncodeToString(account.Address.Bytes())
}

func PrivateToBase64(privateKey *ecdsa.PrivateKey) string {
	return base64.StdEncoding.EncodeToString(crypto.FromECDSA(privateKey))
}

func (account *Wallet) TransferEth(client ethclient.Client, to string, amount int64) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	// Function requires the public address of the account we're sending from -- which we can derive from the private key.

	fromAddress := *account.Address

	// Now we can read the nonce that we should use for the account's transaction.
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", err
	}

	value := big.NewInt(amount) // in wei (1 eth)
	gasLimit := uint64(21000)   // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", err
	}

	// We figure out who we're sending the ETH to.
	toAddress := common.HexToAddress(to)
	var data []byte

	// We create the transaction payload
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return "", err
	}

	// We sign the transaction using the sender's private key
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), account.Private)
	if err != nil {
		return "", err
	}

	// Now we are finally ready to broadcast the transaction to the entire network
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", err
	}

	// We return the transaction hash
	return signedTx.Hash().String(), nil
}
