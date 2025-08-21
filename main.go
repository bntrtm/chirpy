package main

import (
	"log"
	"sync/atomic"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"


	config := apiConfig{}

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", config.middlewareMetricsInc(handler))
	
	// REGISTER HANDLERS
	mux.HandleFunc("GET /api/healthz", endpReadiness)
	mux.HandleFunc("GET /admin/metrics", config.endpFileserverHitCountGet)
	mux.HandleFunc("POST /admin/reset", config.endpFileserverHitCountReset)
	mux.HandleFunc("POST /api/validate_chirp", endpValidateChirp)

	server := &http.Server{
		Addr:		":" + port,
		Handler:	mux,
	}
	
	// start server (NOTE: ListenAndServe always returns non-nil error)
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}