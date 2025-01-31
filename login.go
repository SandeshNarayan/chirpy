package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SandeshNarayan/chirpy/internal/auth"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

    type parameters struct {
        Email string `json:"email"`
        Password string `json:"password"`
		ExpiresInSeconds int `json:"expires_in_seconds"`
    }

	params := parameters{}

	if err := json.NewDecoder(r.Body).Decode(&params); err!=nil {
        respondWithError(w, http.StatusBadRequest, "Could not decode parameters", err)
        return
    }
	
	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err!= nil {
        respondWithError(w, http.StatusNotFound, "User not found", err)
        return
    }

	if err := auth.CheckPasswordHash(params.Password, user.HashedPassword); err!=nil {
        respondWithError(w, http.StatusUnauthorized, "Invalid password", nil)
        return
    }

	expirationTime := time.Hour
	if params.ExpiresInSeconds>0 && params.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(params.ExpiresInSeconds)*time.Second
    }

    tokenString, err := auth.MakeJWT(user.ID, cfg.jwtSecret, expirationTime)
	if err!=nil {
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
	}


	type response struct {
		User
		Token     string    `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	respondWithJson(w, http.StatusOK, response{
		User:User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
        Token:     tokenString,
	})

}