package main

import (
	"log"
	"fmt"
	"slices"
	"strings"
	"net/http"
	"encoding/json"
)

func endpReadiness(w http.ResponseWriter, r *http.Request) {
	respondWithText(w, http.StatusOK, "OK")
}

func(cfg *apiConfig) endpFileserverHitCountGet(w http.ResponseWriter, r *http.Request) {
	output := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())

	respondWithHTML(w, http.StatusOK, output)
}

func(cfg *apiConfig) endpFileserverHitCountReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(fmt.Sprintf("Reset hits to: %d", cfg.fileserverHits.Load()))); err != nil {
		log.Print(err)
	}
}

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

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type returnVals struct {
        Error string `json:"error"`
    }
    respBody := returnVals{
        Error: msg,
    }
    dat, err := json.Marshal(respBody)
	if err != nil {
			log.Printf("Error marshalling JSON for Error response: %s", err)
			w.WriteHeader(500)
			return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
	return
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    dat, err := json.Marshal(payload)
	if err != nil {
			log.Printf("Error marshalling JSON for response: %s", err)
			w.WriteHeader(500)
			return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
	return
}

func respondWithHTML(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Print(err)
	}
	return
}

func respondWithText(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	if _, err := w.Write([]byte(msg)); err != nil {
		log.Print(err)
	}
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