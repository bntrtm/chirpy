package main

import (
	"log"
	"errors"
	"slices"
	"strings"
	"net/http"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/bntrtm/chirpy/internal/database"
)

func(cfg *apiConfig) endpGetRecentChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetRecentChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get recent chirps")
		return
	}

	var respBody []Chirp
	for _, chirp := range chirps {
		addChirp := Chirp{
			ID:			chirp.ID,
			CreatedAt:	chirp.CreatedAt,
			UpdatedAt:	chirp.UpdatedAt,
			Body:     	chirp.Body,
			UserID:		chirp.UserID,
		}
		respBody = append(respBody, addChirp)
	}
	
	respondWithJSON(w, http.StatusOK, respBody)
	return
}

func(cfg *apiConfig) endpCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:	cleaned,
		UserID:	params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp")
		return
	}
	respBody := Chirp{
		ID:			dbChirp.ID,
		CreatedAt:	dbChirp.CreatedAt,
		UpdatedAt:	dbChirp.UpdatedAt,
		Body:     	dbChirp.Body,
		UserID:		dbChirp.UserID,
	}

	respondWithJSON(w, http.StatusCreated, respBody)
	return
}

func validateChirp(body string) (string, error) {
    const maxChirpLength = 140
	const censor = "****"
	profanities := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	if len(body) > 140 {
		return "", errors.New("Chirp is too long (use 140 characters or less)")
	}
	
	return cleanChirpBody(body, censor, profanities), nil
}

func cleanChirpBody(body, censor string, profanities []string) string {

	cleanedBody := ""

	words := strings.Split(body, " ")
	for i, word := range words {
		if i > 0 { cleanedBody += " "}
		if slices.Contains(profanities, strings.ToLower(word)) {
			cleanedBody += censor
		} else {
			cleanedBody += word
		}
	}

	return cleanedBody
}

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
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user")
		return
	}
	respBody := User{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
	}

	respondWithJSON(w, http.StatusCreated, respBody)
	return
}

