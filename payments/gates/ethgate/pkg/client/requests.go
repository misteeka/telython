package client

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
	ethapi "telython/payments/gates/ethgate/pkg/ethereum/api"
	"telython/pkg/http"
	client "telython/pkg/http/client"
)

var httpclient *client.Client

func init() {
	httpclient = client.New("127.0.0.1:8003", "/")
}

func GetAddress(id uint64) (*common.Address, *http.Error, error) {
	value, err := httpclient.Get("getAddress?id=" + strconv.FormatUint(id, 10))
	if err != nil {
		return nil, nil, err
	}
	if client.GetError(value) != nil {
		return nil, client.GetError(value), nil
	}
	addressBase64 := string(value.GetStringBytes("address"))
	address, err := ethapi.Base64ToAddress(addressBase64)
	if err != nil {
		return nil, nil, err
	}
	return address, nil, nil
}

func GetWallet(id uint64) (*ethapi.Wallet, *http.Error, error) {
	private, requestError, err := GetPrivate(id)
	if err != nil {
		return nil, nil, err
	}
	if requestError == nil {
		wallet, err := ethapi.GetWallet(private)
		if err != nil {
			return nil, nil, err
		}
		return wallet, nil, nil
	} else {
		return nil, requestError, nil
	}
}

func CreateWallet(id uint64) (*ethapi.Wallet, *http.Error, error) {
	value, err := httpclient.Post("createWallet", fmt.Sprintf(`{"id":%d}`, id))
	if err != nil {
		return nil, nil, err
	}
	if client.GetError(value) != nil {
		return nil, client.GetError(value), nil
	}
	privateBase64 := string(value.GetStringBytes("data"))
	private, err := ethapi.Base64ToPrivate(privateBase64)
	if err != nil {
		return nil, nil, err
	}
	wallet, err := ethapi.GetWallet(private)
	if err != nil {
		return nil, nil, err
	}
	return wallet, nil, nil
}

func GetPrivate(id uint64) (*ecdsa.PrivateKey, *http.Error, error) {
	value, err := httpclient.Get("getPrivate?id=" + strconv.FormatUint(id, 10))
	if err != nil {
		return nil, nil, err
	}
	if client.GetError(value) != nil {
		return nil, client.GetError(value), nil
	}
	privateBase64 := string(value.GetStringBytes("data"))
	private, err := ethapi.Base64ToPrivate(privateBase64)
	if err != nil {
		return nil, nil, err
	}
	return private, nil, nil
}
