package main

import (
	_ "github.com/lib/pq"

	"log"
	"os"
	"fmt"
	"database/sql"
	"net/http"

	"github.com/joho/godotenv"
	
	"github.com/bntrtm/chirpy/internal/database"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	config := apiConfig{}
	dbQueries := database.New(db)
	config.db = dbQueries
	config.platform = os.Getenv("PLATFORM")
	config.secret = os.Getenv("SECRET")

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", config.middlewareMetricsInc(handler))
	
	// REGISTER HANDLERS
	mux.HandleFunc("GET /api/healthz", endpReadiness)
	mux.HandleFunc("GET /admin/metrics", config.endpFileserverHitCountGet)
	mux.HandleFunc("POST /admin/reset", config.endpDeleteAllUsers)
	mux.HandleFunc("POST /api/users", config.endpCreateUser)
	mux.HandleFunc("POST /api/chirps", config.endpCreateChirp)
	mux.HandleFunc("GET /api/chirps", config.endpGetRecentChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", config.endpGetChirpByID)
	mux.HandleFunc("POST /api/login", config.endpLoginUser)
	mux.HandleFunc("POST /api/refresh", config.endpCheckRefreshToken)
	mux.HandleFunc("POST /api/revoke", config.endpRevokeRefreshToken)
	mux.HandleFunc("PUT /api/users", config.endpUpdateUserCredentials)

	server := &http.Server{
		Addr:		":" + port,
		Handler:	mux,
	}
	
	// start server (NOTE: ListenAndServe always returns non-nil error)
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}