package main

import (
	"net/http"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/bntrtm/chirpy/internal/auth"
	"github.com/bntrtm/chirpy/internal/database"
)

func(cfg *apiConfig) endpUpgradeUserPlan(w http.ResponseWriter, r *http.Request) {
	
	apiKey, ok := (*cfg.apiKeys)["polka"]
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "Internal server error", nil)
		return
	}
	headerKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error(), err)
	}
	if apiKey != headerKey {
		respondWithText(w, http.StatusUnauthorized, "401 Unauthorized")
		return
	}
	type parameters struct {
		Event	string `json:"event"`
		Data	struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
    }
	upgradedUserID, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid id", err)
		return
	}
	if params.Event != "user.upgraded" {
		respondWithText(w, http.StatusNoContent, "204 No Content")
		return
	} else {
		err = cfg.db.UpdateUserPlan(r.Context(), database.UpdateUserPlanParams{
			ID: 		 upgradedUserID,
			IsChirpyRed: true,
		})
		if err != nil {
			respondWithError(w, http.StatusNotFound, err.Error(), err)
			return
		}

		respondWithText(w, http.StatusNoContent, "User upgraded to Chirpy Red plan!")
		return
	}
}

