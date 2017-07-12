package main

import (
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
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

func GetRouter() *negroni.Negroni {
	r := httprouter.New()
	r.GET("/", HomeHandler)
	r.GET("/switches", SwitchesIndexHandler)
	//r.POST("/switches", SwitchesCreateHandler)
	//r.GET("/switches/:id", SwitchesShowHandler)
	//r.PUT("/switches/:id", SwitchUpdateHandler)
	//r.GET("/switches/:id/edit", SwitchEditHandler)

	negroniLogger := negroni.NewLogger()
	negroniLogger.SetDateFormat("[15:04:05 Mon 02 Jan UTC]")
	n := negroni.New(negroniLogger, negroni.NewStatic(http.Dir("public")))
	n.UseHandler(r)

	return n
}

func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var dat []byte
	dat, err := ioutil.ReadFile("public/index.html")
	if err != nil {
		errChannel <- err
	}
	rw.Write(dat)
}

func SwitchesIndexHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	text := []byte("Hello")
	rw.Write(text)
}
