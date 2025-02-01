package main

import (
	"net/http"

	"github.com/SandeshNarayan/chirpy/internal/auth"
	"github.com/SandeshNarayan/chirpy/internal/database"
	"github.com/google/uuid"
)


func (cfg *apiConfig) handlerChirpsDelete(w http.ResponseWriter, r *http.Request){

	token, err := auth.GetBearerToken(r.Header)
	if err!= nil {
        respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
        return
    }

	userID, err:= auth.ValidateJWT(token, cfg.jwtSecret)
	if err!=nil{
        respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
        return
    }

	id:= r.PathValue("chirpID")
	chirpID, err := uuid.Parse(id)
	if err!=nil{
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
        return
	}

	chirp, err:= cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err!=nil{
        respondWithError(w, http.StatusNotFound, "Chirp not found", err)
        return
    }

	if chirp.UserID!=userID{
		respondWithError(w, http.StatusForbidden, "You are not the owner of this chirp", nil)
        return
	}

	err = cfg.dbQueries.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{
		ID:      chirpID,
        UserID: userID,
	})
	if err!=nil{
		respondWithError(w, http.StatusNotFound, "Couldnt delete chirp", err)
        return
	}

    w.WriteHeader(http.StatusNoContent)
}