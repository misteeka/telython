package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"telython/pkg/http"
	"telython/pkg/http/server"
)

func registerHandlers() {
	server.Get("/getPrivate", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		secret_key := ctx.FormValue("key")

		if secret_key != "qazwsx" {
			return http.ToError(http.AUTHORIZATION_FAILED), nil
		}

		private, requestError := getPrivate(username)
		if requestError != nil {
			return requestError, nil
		}
		return nil, private

	}))
	server.Get("/getAddress", server.ReturnDataHandler(func(ctx *fiber.Ctx) (*http.Error, interface{}) {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")

		authorizationError := server.Authorize(username, password)
		if authorizationError != nil {
			return authorizationError, nil
		}

		address, requestError := getAddress(username)
		if requestError != nil {
			return requestError, nil
		}
		return nil, fmt.Sprintf(`{"address": "%s"}`, address)
	}))
}

func shutdown() {
	// shutdown logic
}
