package main

import (
	"encoding/json"
	"net/http"

	"github.com/SandeshNarayan/chirpy/internal/auth"
	"github.com/SandeshNarayan/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request){


	accessToken, err := auth.GetBearerToken(r.Header)
	if err!=nil {
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
        return
    }
	

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
	if err!=nil{
        respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
        return
    }

	


	type parameters struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	var params parameters

	if err:= json.NewDecoder(r.Body).Decode(&params); err!=nil{
  		respondWithError(w, http.StatusBadRequest, "Invalid request parameters", err)
		return
	}

	HashPassword, err := auth.HashPassword(params.Password)
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "Couldnt hash password", err)
        return
	}

	newUser, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: params.Email,
		HashedPassword: HashPassword,
		ID: userID,
    })
	if err!=nil{
        respondWithError(w, http.StatusInternalServerError, "Couldnt update user", err)
        return
    }

	type response struct {
		User
	}
	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID: newUser.ID,
            CreatedAt: newUser.CreatedAt,
            UpdatedAt: newUser.UpdatedAt,
            Email: newUser.Email,
			IsChirpyRed: newUser.IsChirpyRed,
		},
	})

}