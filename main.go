package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	server := http.Server{}
	server.Handler = mux
	server.Addr = ":8080"
	
	// start server (NOTE: ListenAndServe always returns non-nil error)
	server.ListenAndServe()
}