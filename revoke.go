package main

import (
	"net/http"

	"github.com/SandeshNarayan/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request){

	accessToken, err := auth.GetBearerToken(r.Header)
	if err!=nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
        return
    }

	err= cfg.dbQueries.RevokeToken(r.Context(), accessToken)
	if err!=nil{
        respondWithError(w, http.StatusInternalServerError, "Could not revoke token", err)
        return
    }

	w.WriteHeader(http.StatusNoContent)



}