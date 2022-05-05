package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"math/big"
	"strconv"
	"telython/payments/pkg/currency"
	"telython/payments/service/pkg/payments"
	"telython/pkg/cfg"
	"telython/pkg/http"
	"telython/pkg/http/server"
	"telython/pkg/utils"
)

func registerHandlers() {
	server.Post("/payments/addPayment", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			return http.ToError(http.INVALID_REQUEST)
		}
		sender := string(data.GetStringBytes("sender"))
		receiver := string(data.GetStringBytes("receiver"))
		amount, err := utils.DecodeBigInt(string(data.GetStringBytes("amount")))
		if err != nil {
			return http.ToError(http.INVALID_REQUEST)
		}
		currencyCode := data.GetUint64("currency")
		secretKey := string(data.GetStringBytes("secretKey"))

		currency := &currency.Currency{
			Type:   currency.FromCode(currencyCode),
			Amount: amount,
		}
		if currency.Type == nil {
			return http.ToError(http.INVALID_CURRENCY_CODE)
		}

		if secretKey != cfg.GetString("secretKey") {
			return http.ToError(http.AUTHORIZATION_FAILED)
		}

		timestamp, timestampError := server.GetUniqueTimestamp("admin")
		if timestampError != nil {
			return timestampError
		}

		return addPayment(sender, receiver, currency, timestamp)
	}))
	server.Post("/payments/sendPayment", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			return http.ToError(http.INVALID_REQUEST)
		}

		sender := string(data.GetStringBytes("sender"))
		receiver := string(data.GetStringBytes("receiver"))
		amount, ok := new(big.Int).SetString(string(data.GetStringBytes("amount")), 10)
		if !ok {
			return http.ToError(http.INVALID_REQUEST)
		}
		currencyCodeFrom := data.GetUint64("currencyFrom")
		currencyTo := data.GetUint64("currencyTo")
		password := string(data.GetStringBytes("password"))

		if sender == receiver && currencyTo == currencyCodeFrom {
			return &http.Error{
				Code:    http.INVALID_REQUEST,
				Message: "You Cannot Convert The Same Currency",
			}
		}

		currencyFrom := &currency.Currency{
			Type:   currency.FromCode(currencyCodeFrom),
			Amount: amount,
		}

		if currencyFrom.Type == nil {
			return http.ToError(http.INVALID_CURRENCY_CODE)
		}

		authorizationError := server.Authorize(sender, password)
		if authorizationError != nil {
			return authorizationError
		}

		timestamp, timestampError := server.GetUniqueTimestamp(sender)
		if timestampError != nil {
			return timestampError
		}

		return sendPayment(sender, receiver, currencyFrom, currencyTo, timestamp)
	}))
	server.Get("/payments/getAccountInfo", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")

		authorizationError := server.Authorize(username, password)
		if authorizationError != nil {
			return authorizationError, nil
		}

		requestError, account := getAccountInfo(username)
		if requestError != nil {
			return requestError, nil
		}
		return requestError, account.Serialize()
	}))
	server.Get("/payments/getBalance", server.ReturnJsonHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")
		currencyCode, err := strconv.ParseUint(ctx.FormValue("c"), 10, 64)
		if err != nil {
			return http.ToError(http.INVALID_REQUEST), nil
		}

		authorizationError := server.Authorize(username, password)
		if authorizationError != nil {
			return authorizationError, nil
		}

		currencyType := currency.FromCode(currencyCode)

		getError, balance := getBalance(username, currencyType)
		if getError != nil {
			return getError, nil
		}

		return nil, fmt.Sprintf(`{"balance":"%s", "symbol":"%s", "decimals": %d}`, balance.String(), currencyType.Symbol, currencyType.Decimals)
	}))
	server.Get("/payments/getHistory", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")

		authorizationError := server.Authorize(username, password)
		if authorizationError != nil {
			return authorizationError, nil
		}

		requestError, history := getHistory(username)
		if requestError != nil {
			return requestError, nil
		}

		serialized, err := payments.SerializePayments(*history)
		if err != nil {
			return http.ToError(http.INTERNAL_SERVER_ERROR), nil
		}

		if serialized == "" {
			return nil, fmt.Sprintf(`{"payments": "%s"}`, base64.StdEncoding.EncodeToString([]byte("{}")))
		} else {
			return nil, fmt.Sprintf(`{"payments": "%s"}`, base64.StdEncoding.EncodeToString([]byte(serialized)))

		}
	}))
	server.Get("/payments/getPayment", server.ReturnJsonHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		paymentId, err := utils.ParseUint(ctx.FormValue("id"))
		if err != nil {
			return http.ToError(http.INVALID_REQUEST), nil
		}
		sender := ctx.FormValue("s")
		requesterType := ctx.FormValue("t")
		password := ctx.FormValue("p")

		if requesterType == "sender" {
			authorizationError := server.Authorize(sender, password)
			if authorizationError != nil {
				return authorizationError, nil
			}
		}

		getError, payment := getPayment(paymentId, fnv64(sender))
		if getError != nil {
			return getError, nil
		}
		if requesterType == "receiver" {
			authorizationError := server.Authorize(payment.Receiver, password)
			if authorizationError != nil {
				return authorizationError, nil
			}
		}

		return nil, *payment
	}))
	server.Post("/payments/createAccount", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			return http.ToError(http.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		password := string(data.GetStringBytes("password"))

		authorizationError := server.Authorize(username, password)
		if authorizationError != nil {
			return authorizationError
		}

		timestamp, timestampError := server.GetUniqueTimestamp(username)
		if timestampError != nil {
			return timestampError
		}

		requestError := createAccount(username, timestamp)
		return requestError
	}))
}
