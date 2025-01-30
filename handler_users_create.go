package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/SandeshNarayan/chirpy/internal/auth"
	"github.com/SandeshNarayan/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string `json:"email"`
	Password string `json:"-"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter,r  *http.Request){
		
	w.Header().Set("Content-Type", "application/json")

	type parameter struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	type response struct{
		User
	}

	params:= parameter{}
	
	if err := json.NewDecoder(r.Body).Decode(&params); err!=nil{
		respondWithError(w, http.StatusBadRequest, "Couldnt decode parameters", err)
		return
	}



	hashedPassword, err := auth.HashPassword(params.Password)
	if err!=nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt hash password", err)
        return
	}



	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		HashedPassword: hashedPassword,
		Email: params.Email,
	})
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, "Couldnt create users", err)
		return
	}

	

	
	respondWithJson(w, http.StatusCreated, response{
		User:User{
			ID: user.ID,
			CreatedAt: user.CreatedAt,
            UpdatedAt: user.UpdatedAt,
            Email: params.Email,
		},
	})
} 