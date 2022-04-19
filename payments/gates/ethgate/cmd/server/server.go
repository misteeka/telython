package main

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"strconv"
	ethapi "telython/payments/gates/ethgate/pkg/ethereum/api"
	"telython/pkg/http"
	"telython/pkg/http/server"
	"telython/pkg/log"
)

func toBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func registerHandlers() {
	server.Post("/createWallet", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INTERNAL_SERVER_ERROR), nil
		}
		accountId := data.GetUint64("id")
		account, creatingStatus := createWallet(accountId)
		if creatingStatus == nil {
			return nil, account.GetPrivateBase64()
		} else {
			return creatingStatus, nil
		}

	}))
	server.Get("/getPrivate", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		accountId, err := strconv.ParseUint(ctx.FormValue("id"), 10, 64)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INTERNAL_SERVER_ERROR), nil
		}
		private, creatingStatus := getPrivate(accountId)
		if creatingStatus == nil {
			return nil, ethapi.PrivateToBase64(private)
		} else {
			return creatingStatus, nil
		}

	}))
	server.Get("/getAddress", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		accountId, err := strconv.ParseUint(ctx.FormValue("id"), 10, 64)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INTERNAL_SERVER_ERROR), nil
		}
		address, creatingStatus := getAddress(accountId)
		if creatingStatus == nil {
			return nil, ethapi.AddressToBase64(address)
		} else {
			return creatingStatus, nil
		}

	}))
}

func shutdown() {
	// shutdown logic
}
