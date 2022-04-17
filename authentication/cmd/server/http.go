package main

import (
	"crypto/tls"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fastjson"
	"golang.org/x/crypto/acme"
	"golang.org/x/crypto/acme/autocert"
	"net"
	"telython/authentication/pkg/cfg"
	"telython/authentication/pkg/log"
	"telython/authentication/pkg/status"
)

var (
	tlsConfig   *tls.Config
	app         *fiber.App
	tlsListener net.Listener
)

func initFiber() error {
	log.InfoLogger.Println("Fiber initialization")
	app = fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	requestHandler()
	tlsPortListener, err := net.Listen(cfg.GetString("networkName"), "0.0.0.0:443")
	if err != nil {
		return err
	}

	tlsListener = tls.NewListener(tlsPortListener, tlsConfig)
	return nil
}

func initTLS() {
	certManager := &autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("oppositemc.com"),
		Cache:      autocert.DirCache("./certs"),
	}
	tlsConfig = &tls.Config{
		GetCertificate: certManager.GetCertificate,
		NextProtos: []string{
			"http/1.1", acme.ALPNProto,
		},
		InsecureSkipVerify: false,
	}
}

func ListenAndServe() {
	// log("Start HTTPS redirect")
	// go runHTTPSRedirect()
	log.InfoLogger.Println("Fiber server start")

	err := app.Listen(":8001")
	// err := app.Listener(tlsListener)
	panicIfError(err)
}

func runHTTPSRedirect() {

}

func Respond(ctx *fiber.Ctx, status status.Status) error {
	_, err := ctx.WriteString(status.Serialize())
	return err
}

func Deserialize(jsonBytes []byte) (data *fastjson.Value, err error) {
	var p fastjson.Parser
	data, err = p.ParseBytes(jsonBytes)
	return
}

type Handler interface{}

type ReturnStatus func(ctx *fiber.Ctx) status.Status
type ReturnData func(ctx *fiber.Ctx) []byte
type ReturnStatusAndData func(ctx *fiber.Ctx) (status.Status, []byte)

func Post(path string, handler Handler) fiber.Router {
	switch handler := handler.(type) {
	case ReturnStatus:
		return app.Post(path, func(ctx *fiber.Ctx) error {
			return Respond(ctx, handler(ctx))
		})
	}
	log.WarnLogger.Println("Unknown handler " + path)
	return nil
}
func Put(path string, handler Handler) fiber.Router {
	switch handler := handler.(type) {
	case ReturnStatus:
		return app.Put(path, func(ctx *fiber.Ctx) error {
			return Respond(ctx, handler(ctx))
		})
	}
	log.WarnLogger.Println("Unknown handler " + path)
	return nil
}
func Get(path string, handler Handler) fiber.Router {
	switch handler := handler.(type) {
	case ReturnStatus:
		return app.Get(path, func(ctx *fiber.Ctx) error {
			return Respond(ctx, handler(ctx))
		})
	}
	log.WarnLogger.Println("Unknown handler " + path)
	return nil
}
func Delete(path string, handler Handler) fiber.Router {
	switch handler := handler.(type) {
	case ReturnStatus:
		return app.Delete(path, func(ctx *fiber.Ctx) error {
			return Respond(ctx, handler(ctx))
		})
	}
	log.WarnLogger.Println("Unknown handler " + path)
	return nil
}

func requestHandler() {
	Put("auth/signIn", ReturnStatus(func(ctx *fiber.Ctx) status.Status {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return status.INVALID_REQUEST
		}
		username := string(data.GetStringBytes("username"))
		password := string(data.GetStringBytes("password"))
		return signIn(username, password, ctx.IP())
	}))
	Get("auth/checkPassword", ReturnStatus(func(ctx *fiber.Ctx) status.Status {
		username := ctx.FormValue("u")
		password := ctx.FormValue("p")
		return checkPassword(username, password)
	}))
	Put("auth/resetPassword", ReturnStatus(func(ctx *fiber.Ctx) status.Status {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return status.INVALID_REQUEST
		}
		username := string(data.GetStringBytes("username"))
		oldPassword := string(data.GetStringBytes("oldPassword"))
		newPassword := string(data.GetStringBytes("newPassword"))
		return resetPassword(username, oldPassword, newPassword)
	}))
	Post("auth/requestSignUpCode", ReturnStatus(func(ctx *fiber.Ctx) status.Status {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return status.INVALID_REQUEST
		}
		username := string(data.GetStringBytes("username"))
		email := string(data.GetStringBytes("email"))
		return requestSignUpCode(username, email, ctx.IP())
	}))
	Put("auth/requestPasswordRecovery", ReturnStatus(func(ctx *fiber.Ctx) status.Status {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return status.INVALID_REQUEST
		}
		username := string(data.GetStringBytes("username"))
		return requestPasswordRecovery(username)
	}))
	Put("auth/recoverPassword", ReturnStatus(func(ctx *fiber.Ctx) status.Status {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return status.INVALID_REQUEST
		}
		username := string(data.GetStringBytes("username"))
		code := string(data.GetStringBytes("code"))
		newPassword := string(data.GetStringBytes("newPassword"))
		return recoverPassword(username, code, newPassword)
	}))
	Post("auth/signUp", ReturnStatus(func(ctx *fiber.Ctx) status.Status {
		data, err := Deserialize(ctx.Body())
		if err != nil {
			return status.INVALID_REQUEST
		}
		username := string(data.GetStringBytes("username"))
		password := string(data.GetStringBytes("password"))
		code := string(data.GetStringBytes("code"))
		return signUp(username, password, code, ctx.IP())
	}))
}
