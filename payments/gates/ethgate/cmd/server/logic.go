package main

import (
	"crypto/ecdsa"
	"encoding/base64"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"telython/payments/gates/ethgate/pkg/database"
	"telython/payments/gates/ethgate/pkg/ethapi"
	"telython/payments/gates/ethgate/pkg/log"
	"telython/payments/gates/ethgate/pkg/status"
)

func createWallet(accountId uint64) (*ethapi.Wallet, status.Status) {
	wallet, err := ethapi.CreateWallet()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, status.INTERNAL_SERVER_ERROR
	}
	err = database.Accounts.Put(accountId,
		[]string{"id", "public", "private"},
		[]interface{}{accountId, wallet.GetAddressHEX(), wallet.GetPrivateBase64()},
	)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, status.INTERNAL_SERVER_ERROR
	}
	return wallet, status.SUCCESS
}

func getWallet(accountId uint64) (*ethapi.Wallet, status.Status) {
	private, getStatus := getPrivate(accountId)
	if getStatus == status.SUCCESS {
		wallet, err := ethapi.GetWallet(private)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return nil, status.INTERNAL_SERVER_ERROR
		}
		return wallet, status.SUCCESS
	} else {
		return nil, getStatus
	}
}

func getAddress(accountId uint64) (*common.Address, status.Status) {
	base64PublicKey, found, err := database.Accounts.GetString(accountId, "public")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return nil, status.NOT_FOUND
	}
	publicKeyBytes, err := base64.StdEncoding.DecodeString(base64PublicKey)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, status.INTERNAL_SERVER_ERROR
	}

	return ethapi.PublicKeyBytesToAddress(publicKeyBytes), status.SUCCESS
}

func getPrivate(accountId uint64) (*ecdsa.PrivateKey, status.Status) {
	base64PrivateKey, found, err := database.Accounts.GetString(accountId, "private")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, status.INTERNAL_SERVER_ERROR
	}
	if !found {
		return nil, status.NOT_FOUND
	}
	privateKeyBytes, err := base64.StdEncoding.DecodeString(base64PrivateKey)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, status.INTERNAL_SERVER_ERROR
	}
	privateKey, err := crypto.ToECDSA(privateKeyBytes)

	return privateKey, status.SUCCESS
}
