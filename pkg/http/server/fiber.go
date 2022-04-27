package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/valyala/fastjson"
	auth "telython/authentication/pkg/client"
	"telython/pkg/http"
	"telython/pkg/log"
	"time"
)

var App *fiber.App

var ParserPool fastjson.ParserPool

type Handler interface{}

type ReturnJsonHandler func(ctx *fiber.Ctx) (*http.Error, interface{})
type ReturnDataHandler func(ctx *fiber.Ctx) (*http.Error, interface{})
type DefaultHandler func(ctx *fiber.Ctx) *http.Error

func SendJson(ctx *fiber.Ctx, error *http.Error, data interface{}) error {
	if error == nil {
		switch json := data.(type) {
		case string:
			_, err := ctx.WriteString(json)
			return err
		default:
			_, err := ctx.WriteString(string(Serialize(json)))
			return err
		}
	} else {
		_, err := ctx.WriteString(error.Serialize())
		return err
	}
}
func SendStatus(ctx *fiber.Ctx, error *http.Error) error {
	if error == nil {
		_, err := ctx.WriteString(`{"success": true}`)
		return err
	} else {
		_, err := ctx.WriteString(error.Serialize())
		return err
	}
}
func SendData(ctx *fiber.Ctx, error *http.Error, resp interface{}) error {
	if error == nil {
		if resp == nil {
			return ctx.SendStatus(200)
		}
		var data string
		switch v := resp.(type) {
		case string:
			data = fmt.Sprintf("%s", v)
		case []byte:
			data = fmt.Sprintf(`%s`, base64.StdEncoding.EncodeToString(v))
		case int64:
			data = fmt.Sprintf(`%d`, v)
		case int:
			data = fmt.Sprintf(`%d`, v)
		case uint64:
			data = fmt.Sprintf(`%d`, v)
		case bool:
			data = fmt.Sprintf(`%t`, v)
		default:
			data = fmt.Sprintf(`%v`, v)
		}
		_, err := ctx.WriteString(data)
		return err
	} else {
		_, err := ctx.WriteString(error.Serialize())
		return err
	}
}

func Post(path string, handler Handler) fiber.Router {
	switch handler := handler.(type) {
	case DefaultHandler:
		return App.Post(path, func(ctx *fiber.Ctx) error {
			return SendStatus(ctx, handler(ctx))
		})
	case ReturnDataHandler:
		return App.Post(path, func(ctx *fiber.Ctx) error {
			err, resp := handler(ctx)
			if resp == nil {
				return SendStatus(ctx, err)
			} else {
				return SendData(ctx, err, resp)
			}
		})
	case ReturnJsonHandler:
		return App.Post(path, func(ctx *fiber.Ctx) error {
			err, resp := handler(ctx)
			if resp == nil {
				return SendStatus(ctx, err)
			} else {
				return SendJson(ctx, err, resp)
			}
		})
	default:
		return nil
	}
}
func Put(path string, handler Handler) fiber.Router {
	switch handler := handler.(type) {
	case DefaultHandler:
		return App.Put(path, func(ctx *fiber.Ctx) error {
			return SendStatus(ctx, handler(ctx))
		})
	case ReturnDataHandler:
		return App.Put(path, func(ctx *fiber.Ctx) error {
			status, resp := handler(ctx)
			if resp == nil {
				return SendStatus(ctx, status)
			} else {
				return SendData(ctx, status, resp)
			}
		})
	case ReturnJsonHandler:
		return App.Put(path, func(ctx *fiber.Ctx) error {
			status, resp := handler(ctx)
			if resp == nil {
				return SendStatus(ctx, status)
			} else {
				return SendJson(ctx, status, resp)
			}
		})
	default:
		return nil
	}
}
func Delete(path string, handler Handler) fiber.Router {
	switch handler := handler.(type) {
	case DefaultHandler:
		return App.Delete(path, func(ctx *fiber.Ctx) error {
			return SendStatus(ctx, handler(ctx))
		})
	case ReturnDataHandler:
		return App.Delete(path, func(ctx *fiber.Ctx) error {
			status, resp := handler(ctx)
			if resp == nil {
				return SendStatus(ctx, status)
			} else {
				return SendData(ctx, status, resp)
			}
		})
	case ReturnJsonHandler:
		return App.Delete(path, func(ctx *fiber.Ctx) error {
			status, resp := handler(ctx)
			if resp == nil {
				return SendStatus(ctx, status)
			} else {
				return SendJson(ctx, status, resp)
			}
		})
	default:
		return nil
	}
}
func Get(path string, handler Handler) fiber.Router {
	switch handler := handler.(type) {
	case DefaultHandler:
		return App.Get(path, func(ctx *fiber.Ctx) error {
			return SendStatus(ctx, handler(ctx))
		})
	case ReturnDataHandler:
		return App.Get(path, func(ctx *fiber.Ctx) error {
			status, resp := handler(ctx)
			if resp == nil {
				return SendStatus(ctx, status)
			} else {
				return SendData(ctx, status, resp)
			}
		})
	case ReturnJsonHandler:
		return App.Get(path, func(ctx *fiber.Ctx) error {
			status, resp := handler(ctx)
			if resp == nil {
				return SendStatus(ctx, status)
			} else {
				return SendJson(ctx, status, resp)
			}
		})
	default:
		return nil
	}
}

func Init() error {
	App = fiber.New(fiber.Config{
		ReadTimeout:           time.Second * 30,
		WriteTimeout:          time.Second * 30,
		IdleTimeout:           time.Second * 60,
		DisableStartupMessage: true,
	})
	App.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		//AllowHeaders: "Origin, Content-Type, Accept",
	}))
	return nil
}

func Run(addr string) error {
	return App.ListenTLS(addr, "./certs/cert.pem", "./certs/cert.key")
}

func Deserialize(jsonBytes []byte) (data *fastjson.Value, err error) {
	parser := ParserPool.Get()
	data, err = parser.ParseBytes(jsonBytes)
	ParserPool.Put(parser)
	return
}

func GetUniqueTimestamp(username string) (uint64, *http.Error) {
	timestamp := time.Now().UnixMicro()
	return uint64(timestamp), nil
}

func Authorize(username string, password string) *http.Error {
	requestError, err := auth.CheckPassword(username, password)
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		return http.ToError(http.INTERNAL_SERVER_ERROR)
	}
	return requestError
}

func Serialize(i interface{}) []byte {
	jsonData, err := json.Marshal(i)
	if err != nil {
		return nil
	}
	return jsonData
}

/*
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
*/

/*
type Status int

const (
	YES Status = 1
	NO  Status = 0
)

const (
	SUCCESS                Status = 100
	INVALID_REQUEST        Status = 101
	INTERNAL_SERVER_ERROR  Status = 102
	AUTHORIZATION_FAILED   Status = 103
	INVALID_CURRENCY_CODE  Status = 104
	CURRENCY_CODE_MISMATCH Status = 105
	NOT_FOUND              Status = 106
	WRONG_AMOUNT           Status = 107
	INSUFFICIENT_FUNDS     Status = 108
	TOO_MANY_REQUESTS      Status = 109
	ALREADY_EXISTS         Status = 110
)

func (status Status) Equals(anotherStatus Status) bool {
	return status == anotherStatus
}
func (status Status) Serialize() string {
	return strconv.Itoa(int(status))
}

func (status Status) ToString() string {
	if status == SUCCESS {
		return "SUCCESS"
	}
	if status == INVALID_REQUEST {
		return "INVALID REQUEST"
	}
	if status == INTERNAL_SERVER_ERROR {
		return "INTERNAL SERVER ERROR"
	}
	if status == AUTHORIZATION_FAILED {
		return "AUTHORIZATION FAILED"
	}
	if status == INVALID_CURRENCY_CODE {
		return "INVALID CURRENCY CODE"
	}
	if status == CURRENCY_CODE_MISMATCH {
		return "CURRENCY CODE MISMATCH"
	}
	if status == WRONG_AMOUNT {
		return "WRONG AMOUNT"
	}
	if status == NOT_FOUND {
		return "NOT FOUND"
	}
	if status == INSUFFICIENT_FUNDS {
		return "INSUFFICIENT FUNDS"
	}
	if status == TOO_MANY_REQUESTS {
		return "TOO MANY REQUESTS"
	}
	return fmt.Sprintf("%v", status)
}
*/
