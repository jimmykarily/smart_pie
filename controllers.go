package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"html/template"
	"net/http"
	"path"
)

type TemplateData struct {
	Request  http.Request
	DataJson string
}

type SwitchUpdateData struct {
	node    string
	pin     int
	checked bool
}

var FuncMap = template.FuncMap{
	"eq": func(a, b interface{}) bool {
		return a == b
	},
}

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
	r.POST("/switches", SwitchesUpdateHandler)
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
	tmplPath := path.Join("views", "home.html")
	tmpl, err := template.ParseFiles(tmplPath)
	tmpl.Funcs(FuncMap) // Use the 'eq' function

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(rw, TemplateData{*r, "test"}); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func SwitchesUpdateHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var action string
	r.ParseForm()
	if r.Form["checked"][0] == "true" {
		action = "high"
	} else {
		action = "low"
	}
	message := fmt.Sprintf("%s/pin/%s", r.Form["node"][0], r.Form["pin"][0])
	nodeManager.SendMessage(message, action)
}

func SwitchesIndexHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	/*
		tmplPath := path.Join("views", "home.html")
		tmpl, err := template.ParseFiles(tmplPath)
		tmpl.Funcs(FuncMap) // Use the 'eq' function

		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(rw, TemplateData{*r, "test"}); err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	*/

	nodes := nodeManager.Nodes
	json.NewEncoder(rw).Encode(nodes)
}
