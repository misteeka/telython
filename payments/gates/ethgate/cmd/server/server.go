package main

import (
	"encoding/base64"
	"github.com/gofiber/fiber/v2"
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
		username := string(data.GetStringBytes("username"))
		account, creatingStatus := createWallet(username)
		if creatingStatus != nil {
			return creatingStatus, nil
		}

		return nil, account.GetPrivateBase64()

	}))
	server.Get("/getPrivate", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		private, creatingStatus := getPrivate(username)
		if creatingStatus != nil {
			return creatingStatus, nil
		}

		return nil, ethapi.PrivateToBase64(private)
	}))
	server.Get("/getAddress", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		address, creatingStatus := getAddress(username)
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
