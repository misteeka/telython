package main

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"telython/payments/gates/ethgate/pkg/ethapi"
	"telython/payments/gates/ethgate/pkg/http"
	"telython/payments/gates/ethgate/pkg/log"
	"telython/payments/gates/ethgate/pkg/status"
)

func toBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func registerHandlers() {
	http.Post("/createWallet", http.ReturnDataHandler(func(ctx *fiber.Ctx) (status.Status, interface{}) {
		data, err := http.Deserialize(ctx.Body())
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return status.INTERNAL_SERVER_ERROR, nil
		}
		accountId := data.GetUint64("id")
		account, creatingStatus := createWallet(accountId)
		if creatingStatus == status.SUCCESS {
			return status.SUCCESS, account.GetPrivateBase64()
		} else {
			return creatingStatus, nil
		}

	}))
	http.Get("/getPrivate", http.ReturnDataHandler(func(ctx *fiber.Ctx) (status.Status, interface{}) {
		accountId, err := strconv.ParseUint(ctx.FormValue("id"), 10, 64)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return status.INTERNAL_SERVER_ERROR, nil
		}
		private, creatingStatus := getPrivate(accountId)
		if creatingStatus == status.SUCCESS {
			return status.SUCCESS, ethapi.PrivateToBase64(private)
		} else {
			return creatingStatus, nil
		}

	}))
	http.Get("/getAddress", http.ReturnDataHandler(func(ctx *fiber.Ctx) (status.Status, interface{}) {
		accountId, err := strconv.ParseUint(ctx.FormValue("id"), 10, 64)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return status.INTERNAL_SERVER_ERROR, nil
		}
		address, creatingStatus := getAddress(accountId)
		if creatingStatus == status.SUCCESS {
			return status.SUCCESS, toBase64(address.Bytes())
		} else {
			return creatingStatus, nil
		}

	}))
}

func shutdown() {
	// shutdown logic
}
