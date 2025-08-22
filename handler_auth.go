package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/bntrtm/chirpy/internal/auth"
)

func(cfg *apiConfig) endpLoginUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password			string	`json:"password"`
		Email				string	`json:"email`
		ExpiresInSeconds	int		`json: expires_in_seconds`
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

	var expirationTime time.Duration
	if params.ExpiresInSeconds > 0 && params.ExpiresInSeconds <= 3600 {
		expirationTime = time.Second * time.Duration(params.ExpiresInSeconds)
	} else {
		expirationTime = time.Hour * 1
	}

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Trouble logging in", err)
		return
	}
	
	respBody := User{
		ID:        		dbUser.ID,
		CreatedAt: 		dbUser.CreatedAt,
		UpdatedAt: 		dbUser.UpdatedAt,
		Email:     		dbUser.Email,
		Token:			token,
	}

	respondWithJSON(w, http.StatusOK, respBody)
	return
}