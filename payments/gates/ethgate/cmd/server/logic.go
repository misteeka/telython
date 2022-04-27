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
	"telython/pkg/utils"
)

func createWallet(username string) (*ethapi.Wallet, *http.Error) {
	wallet, err := ethapi.CreateWallet()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	log.InfoLogger.Println(wallet.GetAddressHEX())
	err = database.AccountToWallet.Put(utils.Fnv64(username),
		[]string{"id", "address", "private"},
		[]interface{}{utils.Fnv64(username), wallet.GetAddressBase64(), wallet.GetPrivateBase64()},
	)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	err = database.WalletToAccount.Put(wallet.GetAddressBase64(),
		[]string{"name", "address"},
		[]interface{}{username, wallet.GetAddressBase64()},
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
	base64Address, found, err := database.AccountToWallet.GetString(utils.Fnv64(username), "address")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		wallet, requestError := createWallet(username)
		if requestError != nil {
			log.ErrorLogger.Println(requestError.Message)
			return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
		}
		return wallet.Address, nil
	}
	address, err := ethapi.Base64ToAddress(base64Address)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return address, nil
}

func getPrivate(username string) (*ecdsa.PrivateKey, *http.Error) {
	base64PrivateKey, found, err := database.AccountToWallet.GetString(utils.Fnv64(username), "private")
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
