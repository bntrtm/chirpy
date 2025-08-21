package main

import (
	"log"
	"slices"
	"strings"
	"net/http"
	"encoding/json"
)

func endpValidateChirp(w http.ResponseWriter, r *http.Request){
    type parameters struct {
        Body string `json:"body"`
    }

    decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(500)
		return
    }

	if len(params.Body) > 140 {
		errMsg := "Chirp is too long (use 140 characters or less)"
		respondWithError(w, 400, errMsg)
	}

    type returnVals struct {
        CleanedBody string `json:"cleaned_body"`
    }
	respBody := returnVals{
        CleanedBody: cleanChirpBody(params.Body),
    }

	respondWithJSON(w, 200, respBody)
	return
}

func cleanChirpBody(body string) string {
	const censored = "****"
	profanities := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}

	cleanedBody := ""

	words := strings.Split(body, " ")
	for i, word := range words {
		if i > 0 { cleanedBody += " "}
		if slices.Contains(profanities, strings.ToLower(word)) {
			cleanedBody += censored
		} else {
			cleanedBody += word
		}
	}

	return cleanedBody
}