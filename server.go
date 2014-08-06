package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	router := httprouter.New()

	router.PUT("/upload/:id", handleFilePut)
	router.POST("/upload/:id", handleCommit)
}

func handleFilePut(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
	if id := p.ByName("id"); id == "0" {
		// handle chunk put and return session id
	} else {
		// identify upload session and put chunk accordinly.
		// if session not found - send HTTP 404
	}
}

func handleCommit(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

}
