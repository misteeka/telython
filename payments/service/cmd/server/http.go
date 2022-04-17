package main

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
	"telython/payments/service/pkg/http"
	"telython/payments/service/pkg/log"
	"telython/payments/service/pkg/payments"
	"telython/payments/service/pkg/status"
)

func registerHandlers() {
	http.Post("/payments/sendPayment", http.DefaultHandler(func(ctx *fiber.Ctx) status.Status {
		data, err := http.Deserialize(ctx.Body())
		if err != nil {
			return status.INVALID_REQUEST
		}
		sender := data.GetUint64("sender")
		receiver := data.GetUint64("receiver")
		amount := data.GetUint64("amount")
		password := string(data.GetStringBytes("password"))

		username, found, err := getUsername(sender)
		if err != nil {
			return status.INTERNAL_SERVER_ERROR
		}
		if !found {
			return status.AUTHORIZATION_FAILED
		}

		authorizationStatus := http.Authorize(username, password)
		if authorizationStatus == status.SUCCESS {
			timestamp, timestampStatus := http.GetUniqueTimestamp(username)
			if timestampStatus == status.SUCCESS {
				return sendPayment(sender, receiver, amount, timestamp)
			} else {
				return timestampStatus
			}
		} else {
			return authorizationStatus
		}
	}))
	http.Get("/payments/getAccountInfo", http.ReturnDataHandler(func(ctx *fiber.Ctx) (status.Status, interface{}) {
		accountId, err := strconv.ParseUint(ctx.FormValue("a"), 10, 64)
		if err != nil {
			return status.INVALID_REQUEST, nil
		}
		password := ctx.FormValue("p")

		username, found, err := getUsername(accountId)
		if err != nil {
			return status.INTERNAL_SERVER_ERROR, nil
		}
		if !found {
			return status.AUTHORIZATION_FAILED, nil
		}

		authorizationStatus := http.Authorize(username, password)
		if authorizationStatus == status.SUCCESS {
			requestStatus, account := getAccountInfo(accountId)
			if requestStatus != status.SUCCESS {
				return requestStatus, nil
			}
			return requestStatus, http.Serialize(*account)
		} else {
			return authorizationStatus, nil
		}
	}))
	http.Get("/payments/getBalance", http.ReturnDataHandler(func(ctx *fiber.Ctx) (status.Status, interface{}) {
		accountId, err := strconv.ParseUint(ctx.FormValue("a"), 10, 64)
		if err != nil {
			return status.INVALID_REQUEST, nil
		}
		password := ctx.FormValue("p")

		username, found, err := getUsername(accountId)
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return status.INTERNAL_SERVER_ERROR, nil
		}
		if !found {
			return status.AUTHORIZATION_FAILED, nil
		}

		authorizationStatus := http.Authorize(username, password)
		if authorizationStatus == status.SUCCESS {
			return getBalance(accountId)
		} else {
			return authorizationStatus, nil
		}
	}))
	http.Get("/payments/getHistory", http.ReturnDataHandler(func(ctx *fiber.Ctx) (status.Status, interface{}) {
		accountId, err := strconv.ParseUint(ctx.FormValue("a"), 10, 64)
		if err != nil {
			return status.INVALID_REQUEST, nil
		}
		password := ctx.FormValue("p")

		username, found, err := getUsername(accountId)
		if err != nil {
			return status.INTERNAL_SERVER_ERROR, nil
		}
		if !found {
			return status.AUTHORIZATION_FAILED, nil
		}

		authorizationStatus := http.Authorize(username, password)
		if authorizationStatus == status.SUCCESS {
			requestStatus, history := getHistory(accountId)
			if requestStatus != status.SUCCESS {
				return requestStatus, nil
			}
			bytes, err := payments.SerializePayments(*history)
			if err != nil {
				return status.INTERNAL_SERVER_ERROR, nil
			}
			return status.SUCCESS, bytes
		} else {
			return authorizationStatus, nil
		}
	}))
	http.Get("/payments/getPayment", http.ReturnDataHandler(func(ctx *fiber.Ctx) (status.Status, interface{}) {
		paymentId, err := strconv.ParseUint(ctx.FormValue("id"), 10, 64)
		if err != nil {
			return status.INVALID_REQUEST, nil
		}
		senderId, err := strconv.ParseUint(ctx.FormValue("sender"), 10, 64)
		if err != nil {
			return status.INVALID_REQUEST, nil
		}
		requesterType := ctx.FormValue("t")
		password := ctx.FormValue("p")

		if requesterType == "sender" {
			username, found, err := getUsername(senderId)
			if err != nil {
				return status.INTERNAL_SERVER_ERROR, nil
			}
			if !found {
				return status.AUTHORIZATION_FAILED, nil
			}
			authorizationStatus := http.Authorize(username, password)
			if authorizationStatus != status.SUCCESS {
				return authorizationStatus, nil
			}
		}

		getStatus, payment := getPayment(paymentId, senderId)
		if getStatus != status.SUCCESS {
			return getStatus, nil
		} else {
			if requesterType == "receiver" {
				username, found, err := getUsername(payment.Receiver)
				if err != nil {
					return status.INTERNAL_SERVER_ERROR, nil
				}
				if !found {
					return status.AUTHORIZATION_FAILED, nil
				}
				authorizationStatus := http.Authorize(username, password)
				if authorizationStatus != status.SUCCESS {
					return authorizationStatus, nil
				}
			}
			return status.SUCCESS, http.Serialize(*payment)
		}
	}))
	http.Post("/payments/createAccount", http.ReturnDataHandler(func(ctx *fiber.Ctx) (status.Status, interface{}) {
		data, err := http.Deserialize(ctx.Body())
		if err != nil {
			return status.INVALID_REQUEST, nil
		}
		username := string(data.GetStringBytes("username"))
		currency := data.GetInt("currency")
		password := string(data.GetStringBytes("password"))

		authorizationStatus := http.Authorize(username, password)
		if authorizationStatus == status.SUCCESS {
			timestamp, timestampStatus := http.GetUniqueTimestamp(username)
			if timestampStatus == status.SUCCESS {
				return createAccount(username, currency, timestamp)
			} else {
				return timestampStatus, nil
			}
		} else {
			return authorizationStatus, nil
		}
	}))
}
