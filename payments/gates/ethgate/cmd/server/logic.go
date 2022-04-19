package main

import (
	"crypto/ecdsa"
	"encoding/base64"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"telython/payments/gates/ethgate/pkg/database"
	ethapi "telython/payments/gates/ethgate/pkg/ethereum/api"
	"telython/pkg/http"
	"telython/pkg/log"
)

func createWallet(username string) (*ethapi.Wallet, *http.Error) {
	wallet, err := ethapi.CreateWallet()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	log.InfoLogger.Println(wallet.GetAddressHEX())
	err = database.AccountToWallet.Put(fnv64(username),
		[]string{"id", "address", "private"},
		[]interface{}{fnv64(username), wallet.GetAddressBase64(), wallet.GetPrivateBase64()},
	)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	err = database.WalletToAccount.Put(wallet.GetAddressBase64(),
		[]string{"id", "address"},
		[]interface{}{fnv64(username), wallet.GetAddressBase64()},
	)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return wallet, nil
}

func getWallet(username string) (*ethapi.Wallet, *http.Error) {
	private, getStatus := getPrivate(username)
	if getStatus == nil {
		wallet, err := ethapi.GetWallet(private)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
		}
		return wallet, nil
	} else {
		return nil, getStatus
	}
}

func getAddress(username string) (*common.Address, *http.Error) {
	base64Address, found, err := database.AccountToWallet.GetString(fnv64(username), "address")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return nil, &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Wallet Not Found",
		}
	}
	address, err := ethapi.Base64ToAddress(base64Address)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return address, nil
}

func getPrivate(username string) (*ecdsa.PrivateKey, *http.Error) {
	base64PrivateKey, found, err := database.AccountToWallet.GetString(fnv64(username), "private")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return nil, &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Wallet Not Found",
		}
	}
	privateKeyBytes, err := base64.StdEncoding.DecodeString(base64PrivateKey)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	privateKey, err := crypto.ToECDSA(privateKeyBytes)

	return privateKey, nil
}

func fnv64(key string) uint64 {
	hash := uint64(4332272522)
	const prime64 = uint64(33555238)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime64
		hash ^= uint64(key[i])
	}
	return hash
}
