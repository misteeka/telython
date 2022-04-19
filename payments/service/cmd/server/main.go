package main

import (
	"math/rand"
	"runtime"
	"telython/payments/service/pkg/database"
	"telython/pkg/cfg"
	"telython/pkg/http/server"
	"telython/pkg/log"
	"time"
)

func panicIfError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	runtime.GOMAXPROCS(16)
	rand.Seed(time.Now().UnixNano())

	log.InfoLogger.Println("Starting...")
	log.InfoLogger.Println("Config loading")
	panicIfError(cfg.LoadConfig())
	log.InfoLogger.Println("Database start")
	panicIfError(database.InitDatabase())

	log.InfoLogger.Println("Fiber initializing")
	server.Init()
	registerHandlers()

	log.InfoLogger.Println("Fiber run")
	panicIfError(server.Run(":8002"))

	log.InfoLogger.Println("Shutdown...")
	log.InfoLogger.Println("Goodbye!")

}
