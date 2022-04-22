package main

import (
	"errors"
	"math/big"
	"telython/payments/exchange/pkg/database"
	"telython/payments/pkg/currency"
	"telython/pkg/http"
	"telython/pkg/log"
	"telython/pkg/utils"
	"time"
)

func Listener() error {
	for _, v := range currency.Types {
		_, requestError := getPrice(v.Symbol)
		if requestError != nil {
			if requestError.Code == http.NOT_FOUND {
				err := database.Prices.Put(v.Symbol, []string{"symbol", "price"}, []interface{}{v.Symbol, utils.EncodeBigInt(new(big.Int).Div(currency.GetCurrency(v.Symbol, 1, 0).Amount, v.Decimals))})
				if err != nil {
					return err
				}
			} else {
				return errors.New(requestError.Message)
			}
		}
	}
	err := database.Prices.SingleSet("ETH", "price", utils.EncodeBigInt(new(big.Int).Div(currency.GetCurrency("USD", 308562, 2).Amount, currency.Types["ETH"].Decimals)))
	if err != nil {
		return err
	}
	for {
		time.Sleep(5 * time.Minute)
		// get price
		err = database.Prices.SingleSet("ETH", "price", utils.EncodeBigInt(new(big.Int).Div(currency.GetCurrency("USD", 308562, 2).Amount, currency.Types["ETH"].Decimals)))
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			continue
		}
	}
}
