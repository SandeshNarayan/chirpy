package main

import (
	"encoding/json"
	"net/http"

	"github.com/SandeshNarayan/chirpy/internal/auth"
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

	type response struct {
		User
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID: user.ID,
			Email: user.Email,
		},
	})

}