package main

import (
	"math/big"
	"telython/payments/exchange/pkg/database"
	"telython/payments/pkg/currency"
	"telython/pkg/http"
	"telython/pkg/log"
	"telython/pkg/utils"
)

func getPrice(symbol string) (*big.Int, *http.Error) {
	base64Price, found, err := database.Prices.GetString(symbol, "price")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	if !found {
		return nil, &http.Error{
			Code:    http.NOT_FOUND,
			Message: "Price For The Symbol Not Found!",
		}
	}
	price, err := utils.DecodeBigInt(base64Price)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return nil, http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return price, nil
}

func getPriceFromCode(code uint64) (*big.Int, *http.Error) {
	return getPrice(currency.FromCode(code).Symbol)
}

func convert(from uint64, to uint64, amount *big.Int) (*http.Error, *big.Int) {
	priceFrom, requestError := getPriceFromCode(from)
	if requestError != nil {
		return requestError, nil
	}
	priceTo, requestError := getPriceFromCode(to)
	if requestError != nil {
		return requestError, nil
	}

	fundFrom := new(big.Int).Mul(amount, priceFrom)
	fundTo := new(big.Int).Div(fundFrom, priceTo)

	return nil, fundTo
}
