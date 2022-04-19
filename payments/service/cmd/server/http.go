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
)

func DecodeBigInt(base64Int string) (*big.Int, error) {
	bytes, err := base64.StdEncoding.DecodeString(base64Int)
	if err != nil {
		return nil, err
	}
	bigint := new(big.Int).SetBytes(bytes)
	return bigint, nil
}
func EncodeBigInt(int big.Int) string {
	return base64.StdEncoding.EncodeToString(int.Bytes())
}

func registerHandlers() {
	server.Post("/payments/addPayment", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			return http.ToError(http.INVALID_REQUEST)
		}
		receiver := string(data.GetStringBytes("receiver"))
		amount, err := DecodeBigInt(string(data.GetStringBytes("amount")))
		if err != nil {
			return http.ToError(http.INVALID_REQUEST)
		}
		currencyCode := data.GetUint64("currencyCode")
		password := string(data.GetStringBytes("password"))

		if password == cfg.GetString("secretKey") {
			timestamp, timestampError := server.GetUniqueTimestamp("admin")
			if timestampError == nil {
				return addPayment(receiver, &currency.Currency{
					Type:   currency.FromCode(currencyCode),
					Amount: amount,
				}, timestamp)
			} else {
				return timestampError
			}
		} else {
			return http.ToError(http.AUTHORIZATION_FAILED)
		}
	}))
	server.Post("/payments/sendPayment", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			return http.ToError(http.INVALID_REQUEST)
		}
		sender := string(data.GetStringBytes("sender"))
		receiver := string(data.GetStringBytes("receiver"))
		amount, err := DecodeBigInt(string(data.GetStringBytes("amount")))
		if err != nil {
			return http.ToError(http.INVALID_REQUEST)
		}
		currencyCode := data.GetUint64("currency")
		password := string(data.GetStringBytes("password"))

		authorizationError := server.Authorize(sender, password)
		if authorizationError == nil {
			timestamp, timestampError := server.GetUniqueTimestamp(sender)
			if timestampError == nil {
				return sendPayment(sender, receiver, &currency.Currency{
					Type:   currency.FromCode(currencyCode),
					Amount: amount,
				}, timestamp)
			} else {
				return timestampError
			}
		} else {
			return authorizationError
		}
	}))
	server.Get("/payments/getAccountInfo", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")

		authorizationError := server.Authorize(username, password)
		if authorizationError == nil {
			requestError, account := getAccountInfo(username)
			if requestError != nil {
				return requestError, nil
			}
			return requestError, server.Serialize(*account)
		} else {
			return authorizationError, nil
		}
	}))
	server.Get("/payments/getBalance", server.ReturnJsonHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")
		currencyCode, err := strconv.ParseUint(ctx.FormValue("c"), 10, 64)
		if err != nil {
			return http.ToError(http.INVALID_REQUEST), nil
		}
		authorizationError := server.Authorize(username, password)
		if authorizationError == nil {
			getError, balance := getBalance(username, currencyCode)
			if getError == nil {
				return nil, fmt.Sprintf(`{"balance":"%s"}`, base64.StdEncoding.EncodeToString(balance.Bytes()))
			} else {
				return getError, nil
			}
		} else {
			return authorizationError, nil
		}
	}))
	server.Get("/payments/getHistory", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")

		authorizationError := server.Authorize(username, password)
		if authorizationError == nil {
			requestError, history := getHistory(username)
			if requestError != nil {
				return requestError, nil
			}
			bytes, err := payments.SerializePayments(*history)
			if err != nil {
				return http.ToError(http.INTERNAL_SERVER_ERROR), nil
			}
			return nil, bytes
		} else {
			return authorizationError, nil
		}
	}))
	server.Get("/payments/getPayment", server.ReturnJsonHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		paymentId, err := strconv.ParseUint(ctx.FormValue("id"), 10, 64)
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

		getError, payment := getPayment(paymentId, sender)
		if getError != nil {
			return getError, nil
		} else {
			if requesterType == "receiver" {
				username, found, err := getUsername(payment.Receiver)
				if err != nil {
					return http.ToError(http.INTERNAL_SERVER_ERROR), nil
				}
				if !found {
					return &http.Error{
						Code:    http.AUTHORIZATION_FAILED,
						Message: "Account Not Found",
					}, nil
				}
				authorizationError := server.Authorize(username, password)
				if authorizationError != nil {
					return authorizationError, nil
				}
			}
			return nil, *payment
		}
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

		requestError := createAccount(username)
		return requestError
	}))
}
