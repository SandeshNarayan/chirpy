package main

import (
	"net/http"
	"time"

	"github.com/SandeshNarayan/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request){

	accessToken, err := auth.GetBearerToken(r.Header)
	if err!=nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
        return
    }

	user, err := cfg.dbQueries.GetUserFromToken(r.Context(), accessToken)
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