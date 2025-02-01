package main

import (
	"encoding/json"
	"net/http"

	"github.com/SandeshNarayan/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerWebhooksUpdate(w http.ResponseWriter, r *http.Request){

	apiKey, err:= auth.GetAPIKey(r.Header)
	if err!=nil || apiKey!=cfg.apiKey{
        respondWithError(w, http.StatusUnauthorized, "Invalid API key", err)
        return
    }

	type Parameters struct{
		Event string `json:"event"`
		Data struct{
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	var params Parameters

	if err:= json.NewDecoder(r.Body).Decode(&params); err!=nil{
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}

	if params.Event!="user.upgraded"{
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err= cfg.dbQueries.GetUserByID(r.Context(), params.Data.UserID)
	if err!=nil{
		respondWithError(w, http.StatusNotFound, "user not found", err)
		return
	}

	err = cfg.dbQueries.RichGuy(r.Context(), params.Data.UserID)
	if err!=nil{
        respondWithError(w, http.StatusInternalServerError, "Couldnt upgrade user", err)
        return
    }

	w.WriteHeader(http.StatusNoContent)
}
