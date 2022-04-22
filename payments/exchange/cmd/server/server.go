package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"telython/pkg/cfg"
	"telython/pkg/http"
	"telython/pkg/http/server"
	"telython/pkg/utils"
)

func registerHandlers() {
	server.Get("/getPrice", server.ReturnJsonHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		symbol := ctx.FormValue("symbol")
		secretKey := ctx.FormValue("key")

		if secretKey != cfg.GetString("secretKey") {
			return &http.Error{
				Code:    http.AUTHORIZATION_FAILED,
				Message: "Secret Name Invalid",
			}, nil
		}

		price, requestError := getPrice(symbol)
		if requestError != nil {
			return requestError, nil
		}
		return nil, fmt.Sprintf(`{"price": "%s"}`, utils.EncodeBigInt(price))

	}))

	server.Put("/convert", server.ReturnJsonHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			return http.ToError(http.INVALID_REQUEST), nil
		}
		from := data.GetUint64("from")
		to := data.GetUint64("to")
		amount, err := utils.DecodeBigInt(string(data.GetStringBytes("amount")))
		if err != nil {
			return http.ToError(http.INVALID_REQUEST), nil
		}
		secretKey := string(data.GetStringBytes("key"))

		if secretKey != cfg.GetString("secretKey") {
			return &http.Error{
				Code:    http.AUTHORIZATION_FAILED,
				Message: "Secret Name Invalid",
			}, nil
		}
		requestError, fundTo := convert(from, to, amount)
		if requestError != nil {
			return requestError, nil
		}
		return nil, fmt.Sprintf(`{"fund": "%s"}`, utils.EncodeBigInt(fundTo))
	}))
}

func shutdown() {
	// shutdown logic
}
