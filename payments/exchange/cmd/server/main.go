package main

import (
	"math/rand"
	"telython/payments/exchange/pkg/database"
	"telython/pkg/cfg"
	"telython/pkg/http/server"
	"telython/pkg/log"
	"time"
)

func main() {
	log.InfoLogger.Println("Starting...")
	var err error
	rand.Seed(time.Now().UnixNano())

	log.InfoLogger.Println("Loading config...")
	err = cfg.LoadConfig()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		goto Shutdown
	}

	log.InfoLogger.Println("Database start")
	err = database.InitDatabase()
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		goto Shutdown
	}

	log.InfoLogger.Println("Price listener start")
	go func() {
		err := Listener()
		if err != nil {
			log.ErrorLogger.Println(err.Error())
		}
	}()

	log.InfoLogger.Println("Fiber initializing")
	server.Init()
	registerHandlers()

	log.InfoLogger.Println("Fiber run")
	err = server.Run(":8004")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		goto Shutdown
	}

Shutdown:
	log.InfoLogger.Println("Shutdown...")
	shutdown()
	log.InfoLogger.Println("Goodbye!")
}
