package main

import (
	"net/http"
	"os"
	"time"
)

var port = os.Getenv("SMART_PIE_PORT")

func SetupVars() error {
	if port == "" {
		port = "8080"
	}

	return nil
}

func main() {
	logger := Logger{"Main", os.Stdout}
	if err := SetupVars(); err != nil {
		logger.Log(err.Error())
		os.Exit(1)
	}

	// If something is written in the errChannel, the program prints the error
	// and exits.
	errChannel := make(chan error)

	// The web server
	go StartServer(Logger{"Server", os.Stdout}, errChannel)

	// Just a dummy timer. Should be replaced with something more clever.
	tickChan := time.Tick(1 * time.Second)

	for {
		select {
		case err := <-errChannel:
			logger.Log(err.Error())
			os.Exit(1)
		case <-tickChan:
			logger.Log("Ticking")
		}
	}
}

func StartServer(logger Logger, errChannel chan error) {
	router := GetRouter()
	logger.Log("Listening on port " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		errChannel <- err
	}
}
