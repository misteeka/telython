package main

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
	"telython/payments/gates/ethgate/pkg/ethapi"
	"telython/payments/gates/ethgate/pkg/status"
)

func GetAddress(id uint64) (*common.Address, status.Status, error) {
	value, err := Get("getPrivate?id=" + strconv.FormatUint(id, 10))
	if err != nil {
		return nil, 0, err
	}
	requestStatus := getStatus(value)
	if requestStatus == status.SUCCESS {
		publicBase64 := string(value.GetStringBytes("data"))
		address, err := ethapi.PublicBase64ToAddress(publicBase64)
		if err != nil {
			return nil, 0, err
		}
		return address, status.SUCCESS, nil
	}
	return nil, requestStatus, nil
}

func GetWallet(id uint64) (*ethapi.Wallet, status.Status, error) {
	private, getStatus, err := GetPrivate(id)
	if err != nil {
		return nil, 0, err
	}
	if getStatus == status.SUCCESS {
		wallet, err := ethapi.GetWallet(private)
		if err != nil {
			return nil, 0, err
		}
		return wallet, status.SUCCESS, nil
	} else {
		return nil, getStatus, nil
	}
}

func CreateWallet(id uint64) (*ethapi.Wallet, status.Status, error) {
	value, err := Post("createWallet", fmt.Sprintf(`{"id":%d}`, id))
	if err != nil {
		return nil, 0, err
	}
	requestStatus := getStatus(value)
	if requestStatus == status.SUCCESS {
		privateBase64 := string(value.GetStringBytes("data"))
		fmt.Println(privateBase64)
		private, err := ethapi.Base64ToPrivate(privateBase64)
		if err != nil {
			return nil, 0, err
		}
		wallet, err := ethapi.GetWallet(private)
		if err != nil {
			return nil, 0, err
		}
		return wallet, status.SUCCESS, nil
	}
	return nil, requestStatus, nil
}

func GetPrivate(id uint64) (*ecdsa.PrivateKey, status.Status, error) {
	value, err := Get("getPrivate?id=" + strconv.FormatUint(id, 10))
	if err != nil {
		return nil, 0, err
	}
	requestStatus := getStatus(value)
	if requestStatus == status.SUCCESS {
		privateBase64 := string(value.GetStringBytes("data"))
		private, err := ethapi.Base64ToPrivate(privateBase64)
		if err != nil {
			return nil, 0, err
		}
		return private, status.SUCCESS, nil
	}
	return nil, requestStatus, nil
}
