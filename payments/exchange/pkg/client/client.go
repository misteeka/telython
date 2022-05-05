package client

import (
	"errors"
	"fmt"
	"math/big"
	"telython/payments/pkg/currency"
	"telython/pkg/http"
	httpclient "telython/pkg/http/client"
)

var client *httpclient.Client

func init() {
	client = httpclient.New("127.0.0.1:8004", "/")
}

func GetPrice(symbol string, key string) (*currency.Currency, *http.Error, error) {
	json, err := client.Get(fmt.Sprintf("getPrice?symbol=%s&key=%s", symbol, key))
	if err != nil {
		return nil, nil, err
	}
	price, ok := new(big.Int).SetString(string(json.GetStringBytes("price")), 10)
	if !ok {
		return nil, httpclient.GetError(json), errors.New("Wrong pric!")
	}
	return &currency.Currency{
		Type:   currency.Types["USD"],
		Amount: price,
	}, httpclient.GetError(json), nil
}

func Convert(from *currency.Currency, to uint64, key string) (*http.Error, *currency.Currency, error) {
	json, err := client.Put("convert", fmt.Sprintf(`{"from": %d, "to": %d, "amount":"%s", "key":"%s"}`, from.Type.Id, to, from.Amount.String(), key))
	if err != nil {
		return nil, nil, err
	}
	fund, ok := new(big.Int).SetString(string(json.GetStringBytes("fund")), 10)
	if !ok {
		return httpclient.GetError(json), nil, errors.New("Wrong fund!")
	}
	return httpclient.GetError(json), &currency.Currency{
		Type:   currency.FromCode(to),
		Amount: fund,
	}, nil
}
