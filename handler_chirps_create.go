package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/SandeshNarayan/chirpy/internal/auth"
	"github.com/SandeshNarayan/chirpy/internal/database"
	"github.com/google/uuid"
)


type Chirp struct{
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID uuid.UUID `json:"user_id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request){
	

	tokenString, err := auth.GetBearerToken(r.Header)
	if err!=nil{
		respondWithError(w, http.StatusUnauthorized, "Missing or invalid token", err)
        return
	}

	userID, err:= auth.ValidateJWT(tokenString, cfg.jwtSecret)
	if err!=nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid token", err)
        return
	}
		
	type parameters struct{
		Body string `json:"body"`
	}

	params:=parameters{}
	
	if err:= json.NewDecoder(r.Body).Decode(&params); err!= nil{
		respondWithError(w, http.StatusBadRequest, "Invalid JSON body", err)
		return
	}



	cleaned, err:= validatChirp(params.Body)
	if err!= nil{
        respondWithError(w, http.StatusBadRequest, "Invalid chirp body", err)
        return
    }

	
	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleaned,
		UserID: userID,
	})

    if err!=nil{
        respondWithError(w, http.StatusInternalServerError, "Couldnt create chirp", err)
        return
    }



	respondWithJson(w, http.StatusCreated, Chirp{
		ID: chirp.ID,
        CreatedAt: chirp.CreatedAt,
        UpdatedAt: chirp.UpdatedAt,
        UserID: userID,
        Body: cleaned,
	})


}


func validatChirp(body string)(string, error){
	const maxChirpLength =140


	if len(body)>maxChirpLength {

		return "", errors.New("Chirp is too long")
	}

	badWords:= map[string]struct{}{
		"kerfuffle":{},
		"sharbert":{},
		"fornax":{},
	}

	cleanedBody := getCleanedBody(body, badWords)
	return cleanedBody, nil
}

func getCleanedBody(body string, badWords map[string]struct{}) string{
	words := strings.Split(body, " ") 
	for i, word:= range words {
		loweredWord:= strings.ToLower(word)
		if _,ok:= badWords[loweredWord]; ok{
			words[i] = "****"			
		}
		
	}
	new_body := strings.Join(words, " ")

	return new_body
}