package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SandeshNarayan/chirpy/internal/auth"
	"github.com/SandeshNarayan/chirpy/internal/database"
)


func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")

    type parameters struct {
        Email string `json:"email"`
        Password string `json:"password"`
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

	accessToken,err:= auth.MakeJWT(user.ID, cfg.jwtSecret, time.Hour)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "could not create token", err)
        return
	}

	refreshToken, err:= auth.MakeRefreshToken()
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "could not create refresh token", err)
        return
	} 

	expirationTime := time.Now().Add(60*24*time.Hour)

	
	RefreshToken, err:= cfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams {
		Token :refreshToken,
		UserID: user.ID,
		ExpiresAt: expirationTime,
		})
	if err!=nil{

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
        Token:     accessToken,
		RefreshToken: RefreshToken.Token,
	})

}