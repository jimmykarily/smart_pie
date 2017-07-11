package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

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

func HomeHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Fprintln(rw, "Home")
}

func SwitchesIndexHandler(rw http.ResponseWriter, r *http.Request, p httprouter.Params) {
	text := []byte("Hello")
	rw.Write(text)
}
