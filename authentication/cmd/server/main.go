package main

import (
	"math/rand"
	"runtime"
	"telython/authentication/cmd/server/api"
	"telython/authentication/cmd/server/mail"
	"telython/authentication/pkg/database"
	"telython/pkg/cfg"
	"telython/pkg/http/server"
	"telython/pkg/log"
	"time"
)

func main() {
	log.InfoLogger.Println("Starting...")
	runtime.GOMAXPROCS(8)
	rand.Seed(time.Now().UnixNano())

	log.InfoLogger.Println("Config loading")
	err := cfg.LoadConfig()
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

	log.InfoLogger.Println("Gomail start")
	mail.Init()
	// TestMail()

	log.InfoLogger.Println("Server initialization")
	server.Init()
	api.RegisterHandlers()

	log.InfoLogger.Println("Server run")
	err = server.Run(":8001")
	if err != nil {
		log.ErrorLogger.Println(err.Error())
		goto Shutdown
	}

Shutdown:
	log.InfoLogger.Println("Shutdown...")
	log.InfoLogger.Println("Goodbye!")
}
