package main

import (
	"net/http"
	"os"
)

var port = os.Getenv("SMART_PIE_PORT")
var logger = Logger{"Main", os.Stdout}

func SetupVars() error {
	if port == "" {
		port = "8080"
	}

	return nil
}

func main() {
	if err := SetupVars(); err != nil {
		logger.Log(err.Error())
		os.Exit(1)
	}

	// If something is written in the errChannel, the program prints the error
	// and exits.
	errChannel := make(chan error)

	// The web server
	go StartServer(Logger{"Server", os.Stdout}, errChannel)

	for {
		select {
		case err := <-errChannel:
			logger.Log(err.Error())
			os.Exit(1)
		}
	}
}

func StartServer(logger Logger, errChannel chan error) {
	router := GetRouter(&logger)
	logger.Log("Listening on port " + port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		errChannel <- err
	}
}
