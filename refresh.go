package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/SandeshNarayan/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request){

	refreshToken := r.Header.Get("Authorization")

	parts:= strings.SplitN(refreshToken, " ", 2)
	if len(parts)!=2 || parts[0]!="Bearer"{
        respondWithError(w, http.StatusUnauthorized, "Invalid token", nil)
        return
    }

	token:=parts[1]

	user, err := cfg.dbQueries.GetUserFromRefreshToken(r.Context(), token)
	if err!=nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
        return
    }
	
	

	newAccessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err!=nil{
        respondWithError(w, http.StatusInternalServerError, "could not create token", err)
        return
    }

	type response struct{
		Token string `json:"token"`
	}


	respondWithJson(w, http.StatusOK, response{
		Token: newAccessToken,
	})


}