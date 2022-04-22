package client

import (
	"fmt"
	"telython/payments/pkg/currency"
	"telython/pkg/http"
	httpclient "telython/pkg/http/client"
	"telython/pkg/utils"
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
	price, err := utils.DecodeBigInt(string(json.GetStringBytes("price")))
	if err != nil {
		return nil, httpclient.GetError(json), nil
	}
	return &currency.Currency{
		Type:   currency.Types["USD"],
		Amount: price,
	}, httpclient.GetError(json), nil
}

func Convert(from *currency.Currency, to uint64, key string) (*http.Error, *currency.Currency, error) {
	json, err := client.Put("convert", fmt.Sprintf(`{"from": %d, "to": %d, "amount":"%s", "key":"%s"}`, from.Type.Id, to, utils.EncodeBigInt(from.Amount), key))
	if err != nil {
		return nil, nil, err
	}
	fund, err := utils.DecodeBigInt(string(json.GetStringBytes("fund")))
	if err != nil {
		return httpclient.GetError(json), nil, err
	}
	return httpclient.GetError(json), &currency.Currency{
		Type:   currency.FromCode(to),
		Amount: fund,
	}, nil
}
