package main

import (
	"log"
	"net/http"
	"encoding/json"
)

func(cfg *apiConfig) endpCreateUser(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Email string `json:"email"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	dbUser, err := cfg.db.CreateUser(r.Context(), params.Email)
	respBody := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(201)

	respondWithJSON(w, 201, respBody)
	return
}

