package main

import (
	"encoding/json"
	"net/http"

	//"github.com/bntrtm/chirpy/internal/database"
	"github.com/bntrtm/chirpy/internal/auth"
)

func(cfg *apiConfig) endpLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password	string	`json:"password"`
		Email		string	`json:"email`
	}

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not post chirp", err)
		return
    }

	dbUser, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}
	
	respBody := User{
		ID:        		dbUser.ID,
		CreatedAt: 		dbUser.CreatedAt,
		UpdatedAt: 		dbUser.UpdatedAt,
		Email:     		dbUser.Email,
	}

	respondWithJSON(w, http.StatusOK, respBody)
	return
}