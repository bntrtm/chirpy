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
	"github.com/bntrtm/chirpy/internal/auth"
)

func(cfg *apiConfig) endpGetChirpByID(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("chirpID")
	id, err := uuid.Parse(idString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id", err)
		return
	}
	dbChirp, err := cfg.db.GetChirpByID(r.Context(), id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp at specified id", err)
		return
	}


	respBody := Chirp{
		ID:			dbChirp.ID,
		CreatedAt:	dbChirp.CreatedAt,
		UpdatedAt:	dbChirp.UpdatedAt,
		Body:     	dbChirp.Body,
		UserID:		dbChirp.UserID,
	}
	
	respondWithJSON(w, http.StatusOK, respBody)
	return
}

func(cfg *apiConfig) endpGetRecentChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetRecentChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get recent chirps", err)
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
		respondWithError(w, http.StatusInternalServerError, "Could not post chirp", err)
		return
    }

	cleaned, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	dbChirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:	cleaned,
		UserID:	params.UserID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create chirp", err)
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
		Password	string `json:"password"`
        Email		string `json:"email"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	hashedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failure processing request to create user", err)
		return
	}

	dbUser, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:			params.Email,
		HashedPassword:	hashedPass,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failure processing request to create user", err)
		return
	}
	respBody := User{
		ID:        		dbUser.ID,
		CreatedAt: 		dbUser.CreatedAt,
		UpdatedAt: 		dbUser.UpdatedAt,
		Email:     		dbUser.Email,
		HashedPassword:	dbUser.HashedPassword,
	}

	respondWithJSON(w, http.StatusCreated, respBody)
	return
}

