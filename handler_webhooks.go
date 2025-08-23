package main

import (
	"net/http"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/bntrtm/chirpy/internal/database"
)

func(cfg *apiConfig) endpUpgradeUserPlan(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event	string `json:"event"`
		Data	struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
    params := parameters{}
    err := decoder.Decode(&params)
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

