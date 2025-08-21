package main

import (
	_ "github.com/lib/pq"

	"log"
	"os"
	"fmt"
	"time"
	"database/sql"
	"sync/atomic"
	"net/http"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	
	"github.com/bntrtm/chirpy/internal/database"
)

type apiConfig struct {
	// atomic.Int32 is a //standard-library type that allows us to 
	// safely increment and read an integer value across multiple 
	// goroutines (HTTP requests)
	fileserverHits	atomic.Int32
	db				*database.Queries
	platform		string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) middlewareMetricsReset(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Store(0)
		next.ServeHTTP(w, r)
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

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

	mux := http.NewServeMux()
	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", config.middlewareMetricsInc(handler))
	
	// REGISTER HANDLERS
	mux.HandleFunc("GET /api/healthz", endpReadiness)
	mux.HandleFunc("GET /admin/metrics", config.endpFileserverHitCountGet)
	mux.HandleFunc("POST /admin/reset", config.endpDeleteAllUsers)
	mux.HandleFunc("POST /api/validate_chirp", endpValidateChirp)
	mux.HandleFunc("POST /api/users", config.endpCreateUser)

	server := &http.Server{
		Addr:		":" + port,
		Handler:	mux,
	}
	
	// start server (NOTE: ListenAndServe always returns non-nil error)
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}