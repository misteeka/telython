package api

import (
	"github.com/gofiber/fiber/v2"
	"telython/pkg/http"
	"telython/pkg/http/server"
	"telython/pkg/log"
)

func RegisterHandlers() {
	server.Put("auth/signIn", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		password := string(data.GetStringBytes("password"))
		return signIn(username, password, ctx.IP())
	}))
	server.Get("auth/checkPassword", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")
		return checkPassword(username, password)
	}))
	server.Put("auth/resetPassword", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		oldPassword := string(data.GetStringBytes("oldPassword"))
		newPassword := string(data.GetStringBytes("newPassword"))
		return resetPassword(username, oldPassword, newPassword)
	}))
	server.Post("auth/requestSignUpCode", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		email := string(data.GetStringBytes("email"))
		return requestSignUpCode(username, email, ctx.IP())
	}))
	server.Put("auth/requestPasswordRecovery", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		return requestPasswordRecovery(username)
	}))
	server.Put("auth/recoverPassword", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		code := string(data.GetStringBytes("code"))
		newPassword := string(data.GetStringBytes("newPassword"))
		return recoverPassword(username, code, newPassword)
	}))
	server.Post("auth/signUp", server.DefaultHandler(func(ctx *fiber.Ctx) *http.Error {
		data, err := server.Deserialize(ctx.Body())
		if err != nil {
			log.ErrorLogger.Println(err.Error())
			return http.ToError(http.INVALID_REQUEST)
		}
		username := string(data.GetStringBytes("username"))
		password := string(data.GetStringBytes("password"))
		code := string(data.GetStringBytes("code"))
		return signUp(username, password, code, ctx.IP())
	}))
}
