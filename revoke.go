package main

import (
	"net/http"
	"strings"
)


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request){

	refreshToken:= r.Header.Get("Authorization")

	parts:= strings.SplitN(refreshToken, " ", 2)
	if len(parts)!=2 || parts[0]!="Bearer"{
        respondWithError(w, http.StatusUnauthorized, "Invalid token", nil)
        return
    }

	token:=parts[1]

	err:= cfg.dbQueries.RevokeToken(r.Context(), token)
	if err!=nil{
        respondWithError(w, http.StatusInternalServerError, "Could not revoke token", err)
        return
    }

	w.WriteHeader(http.StatusNoContent)



}