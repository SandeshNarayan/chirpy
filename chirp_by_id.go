package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerFindChirpById(w http.ResponseWriter, r *http.Request){
	
	id:= r.PathValue("chirpID")
	chirpID, err:=uuid.Parse(id)
	if err!=nil{
        respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
        return
    }

    dbchirp, err:= cfg.dbQueries.GetChirpByID(r.Context(), chirpID)
	if err!=nil{
        respondWithError(w, http.StatusNotFound, "Chirp not found", err)
        return
    }

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	respondWithJson(w, http.StatusOK, Chirp{
		ID: dbchirp.ID,
		Body: dbchirp.Body,
        UserID: dbchirp.UserID,
        CreatedAt: dbchirp.CreatedAt,
        UpdatedAt: dbchirp.UpdatedAt,
	})	
}