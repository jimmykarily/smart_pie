package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

// This struct implements http.Handler but also prints the requests
// in stdout.
type LoggingRouter struct {
	handler http.Handler
	logger  Logger
}

func (l LoggingRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//TODO: Log based on logging level env var
	l.logger.Log(r.Method + " " + r.URL.Path)
	l.handler.ServeHTTP(w, r)
}

func GetRouter(logger *Logger) LoggingRouter {
	r := httprouter.New()
	r.GET("/", HomeHandler)
	r.GET("/switches", SwitchesIndexHandler)
	//r.POST("/switches", SwitchesCreateHandler)
	//r.GET("/switches/:id", SwitchesShowHandler)
	//r.PUT("/switches/:id", SwitchUpdateHandler)
	//r.GET("/switches/:id/edit", SwitchEditHandler)

	return LoggingRouter{r, *logger}
}

func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Home")
}

func SwitchesIndexHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	text := []byte("Hello")
	rw.Write(text)
}
