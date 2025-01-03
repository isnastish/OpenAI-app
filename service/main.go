package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/isnastish/openai/pkg/api"
	"github.com/isnastish/openai/pkg/log"
)

func main() {
	port := flag.Int("port", 3030, "Listening port")
	logLevel := flag.String("loglevel", "debug", "Logging level")
	flag.Parse()

	log.SetupGlobalLogLevel(*logLevel)

	// NOTE: This won't work when executed with docker compose,
	// because the .env file won't exist inside a docker container.
	// if err := godotenv.Load(".env"); err != nil {
	// 	log.Logger.Fatal("Failed to load .env file, error: %v", err)
	// }

	app, err := api.NewApp(*port)
	if err != nil {
		log.Logger.Fatal("Error occured on app creation: %v", err)
	}

	osSigChan := make(chan os.Signal, 1)
	signal.Notify(osSigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := app.Serve()
		if err != nil {
			log.Logger.Fatal("Error occured while serving: %v", err)
		}
	}()

	<-osSigChan
	if err := app.Shutdown(); err != nil {
		log.Logger.Fatal("Failed to shutdown the server: %v", err)
	}

	os.Exit(0)
}
