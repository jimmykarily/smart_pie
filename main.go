package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
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

func GetRouter() *httprouter.Router {
	r := httprouter.New()
	r.GET("/", HomeHandler)
	r.GET("/switches", SwitchesIndexHandler)
	//r.POST("/switches", SwitchesCreateHandler)
	//r.GET("/switches/:id", SwitchesShowHandler)
	//r.PUT("/switches/:id", SwitchUpdateHandler)
	//r.GET("/switches/:id/edit", SwitchEditHandler)

	return r
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
	if err := http.ListenAndServe(":8080", router); err != nil {
		errChannel <- err
	}
}

func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Home")
}

func SwitchesIndexHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	text := []byte("Hello")
	rw.Write(text)
}
