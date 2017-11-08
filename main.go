package main

import (
	"net/http"
	"os"
)

var port = os.Getenv("SMART_PIE_PORT")
var mqttBrokerUri = os.Getenv("MQTT_BROKER_URI")
var mqttUsername = os.Getenv("MQTT_USERNAME")
var mqttPassword = os.Getenv("MQTT_PASSWORD")
var logger = Logger{"Main", os.Stdout}

// If something is written in the errChannel, the program prints the error and
// exits.
var errChannel = make(chan error)

func SetupVars() error {
	if port == "" {
		port = "8080"
	}

	if mqttUsername == "" {
		logger.Log("MQTT_USERNAME not set")
		os.Exit(1)
	}

	if mqttPassword == "" {
		logger.Log("MQTT_PASSWORD not set")
		os.Exit(1)
	}

	return nil
}

func main() {
	if err := SetupVars(); err != nil {
		logger.Log(err.Error())
		os.Exit(1)
	}

	// The web server
	go StartServer(Logger{"Server", os.Stdout}, errChannel)

	nodeManager := NewNodeManager(mqttBrokerUri, mqttUsername, mqttPassword)
	nodeManager.Subscribe()

	for {
		select {
		case err := <-errChannel:
			logger.Log(err.Error())
			os.Exit(1)
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
