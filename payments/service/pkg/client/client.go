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

var client *httpclient.Client

func init() {
	client = httpclient.New("127.0.0.1:8002", "/payments/")
}

func SendPayment(sender string, receiver uint64, currency *currency.Currency, currencyCodeTo uint64, password string) (*http.Error, error) {
	if currency.Type == nil {
		return http.ToError(http.INVALID_CURRENCY_CODE), nil
	}
	if currency.Amount == nil {
		return http.ToError(http.WRONG_AMOUNT), nil
	}
	if currency == nil {
		return http.ToError(http.WRONG_AMOUNT), nil
	}
	fmt.Println(fmt.Sprintf(`{"sender":"%s","receiver":%d,"amount":"%s","currencyFrom":%d,"currencyTo":%d, "password":"%s"}`, sender, receiver, base64.StdEncoding.EncodeToString(currency.Amount.Bytes()), currency.Type.Id, currencyCodeTo, password))
	json, err := client.Post("sendPayment", fmt.Sprintf(`{"sender":"%s","receiver":%d,"amount":"%s","currencyFrom":%d,"currencyTo":%d, "password":"%s"}`, sender, receiver, base64.StdEncoding.EncodeToString(currency.Amount.Bytes()), currency.Type.Id, currencyCodeTo, password))
	return httpclient.GetError(json), err
}
func AddPayment(receiver uint64, currency *currency.Currency, secretKey string) (*http.Error, error) {
	if currency.Type == nil {
		return http.ToError(http.INVALID_CURRENCY_CODE), nil
	}
	if currency.Amount == nil {
		return http.ToError(http.WRONG_AMOUNT), nil
	}
	if currency == nil {
		return http.ToError(http.WRONG_AMOUNT), nil
	}
	json, err := client.Post("addPayment", fmt.Sprintf(`{"receiver":%d,"amount":"%s","currency":%d,"secretKey":"%s"}`, receiver, base64.StdEncoding.EncodeToString(currency.Amount.Bytes()), currency.Type.Id, secretKey))
	return httpclient.GetError(json), err
}
func GetBalance(username string, password string, currencyCode uint64) (*big.Int, *http.Error, error) {
	json, err := client.Get(fmt.Sprintf("getBalance?u=%s&p=%s&c=%d", username, password, currencyCode))
	if err != nil {
		return nil, nil, err
	}
	if httpclient.GetError(json) == nil {
		balance, err := base64.StdEncoding.DecodeString(string(json.GetStringBytes("balance")))
		if err != nil {
			return nil, nil, err
		}
		return new(big.Int).SetBytes(balance), nil, nil
	} else {
		return nil, httpclient.GetError(json), nil
	}
}
func CreateAccount(username string, password string) (uint64, *http.Error, error) {
	json, err := client.Post("createAccount", fmt.Sprintf(`{"username":"%s","password":"%s"}`, username, password))
	return json.GetUint64("id"), httpclient.GetError(json), err
}

func GetHistory(username string, password string) ([]payments.Payment, *http.Error, error) {
	json, err := client.Get("getHistory?u=" + username + "&p=" + password)
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
