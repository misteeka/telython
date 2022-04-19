package client

import (
	"encoding/base64"
	"fmt"
	"math/big"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/payments"
	http "telython/pkg/http"
	httpclient "telython/pkg/http/client"
)

func init() {
	httpclient.Init("127.0.0.1:8002", "/payments/")
}

func SendPayment(sender string, receiver string, currency *currency.Currency, password string) (*http.Error, error) {
	json, err := httpclient.Post("sendPayment", fmt.Sprintf(`{"sender":"%s","receiver":"%s",amount":"%s","currency":%d,"password":"%s"}`, sender, receiver, base64.StdEncoding.EncodeToString(currency.Amount.Bytes()), currency.Type.Id, password))
	return httpclient.GetError(json), err
}
func AddPayment(receiver string, currency *currency.Currency, password string) (*http.Error, error) {
	json, err := httpclient.Post("addPayment", fmt.Sprintf(`{"receiver":"%s","amount":"%s","currency":%d,"password":"%s"}`, receiver, base64.StdEncoding.EncodeToString(currency.Amount.Bytes()), currency.Type.Id, password))
	return httpclient.GetError(json), err
}
func GetBalance(username string, password string) (*currency.Currency, *http.Error, error) {
	json, err := httpclient.Get("getBalance?u=" + username + "&p=" + password)
	if err != nil {
		return nil, nil, err
	}
	if httpclient.GetError(json) == nil {
		balance, err := base64.StdEncoding.DecodeString(string(json.GetStringBytes("balance")))
		if err != nil {
			return nil, nil, err
		}
		return &currency.Currency{
			Type:   currency.FromCode(json.GetUint64("currency")),
			Amount: new(big.Int).SetBytes(balance),
		}, nil, nil
	} else {
		return nil, httpclient.GetError(json), nil
	}
}
func CreateAccount(username string, password string) (uint64, *http.Error, error) {
	json, err := httpclient.Post("createAccount", fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	return json.GetUint64("id"), httpclient.GetError(json), err
}

func GetHistory(username string, password string) ([]payments.Payment, *http.Error, error) {
	json, err := httpclient.Get("getHistory?u=" + username + "&p=" + password)
	if err != nil {
		return nil, httpclient.GetError(json), err
	}
	paymentsBytes, err := base64.StdEncoding.DecodeString(string(json.GetStringBytes("data")))
	if err != nil {
		return nil, httpclient.GetError(json), err
	}
	payments, err := payments.DeserializePayments(paymentsBytes)
	if err != nil {
		return nil, httpclient.GetError(json), err
	}
	return *payments, httpclient.GetError(json), nil
}
