package main

import (
	"math/rand"
	"telython/payments/gates/ethgate/pkg/cfg"
	"telython/payments/gates/ethgate/pkg/database"
	"telython/payments/gates/ethgate/pkg/ethapi"
	"telython/payments/gates/ethgate/pkg/http"
	"telython/payments/gates/ethgate/pkg/log"
	"time"
)

func main() {
	log.InfoLogger.Println("Starting...")
	var err error
	rand.Seed(time.Now().UnixNano())

	log.InfoLogger.Println("Loading config...")
	err = cfg.LoadConfig()
	if err != nil {
		goto Shutdown
	}

	log.InfoLogger.Println("Database start")
	err = database.InitDatabase()
	if err != nil {
		goto Shutdown
	}

	log.InfoLogger.Println("Ethapi start")
	err = ethapi.Init()
	if err != nil {
		goto Shutdown
	}

	log.InfoLogger.Println("Fiber initializing")
	http.Init()
	registerHandlers()

	log.InfoLogger.Println("Fiber run")
	err = http.Run()
	if err != nil {
		goto Shutdown
	}

Shutdown:
	log.InfoLogger.Println("Shutdown...")
	shutdown()
	log.InfoLogger.Println("Goodbye!")
}
