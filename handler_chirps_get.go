package main

import (
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request ){
	

	authorIDString := r.URL.Query().Get("author_id")
	sortDirection := r.URL.Query().Get("sort")

	if sortDirection!="" && (sortDirection!="asc" && sortDirection!="desc"){
        respondWithError(w, http.StatusBadRequest, "Invalid sort order", nil)
        return
    }
	


	dbchirps, err:= cfg.dbQueries.GetAllChirps(r.Context())
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
	}
	
	


	if authorIDString!=""{

		authorID, err := uuid.Parse(authorIDString)
		if err!=nil{
			respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
			return
		}

		dbchirps, err = cfg.dbQueries.GetChirpByUserID(r.Context(), authorID)
		if err!=nil{
            respondWithError(w, http.StatusInternalServerError, err.Error(), err)
            return
        }
	}
	if sortDirection == "desc" {
		sort.Slice(dbchirps, func(i, j int) bool {
			return dbchirps[i].CreatedAt.After(dbchirps[j].CreatedAt)
		})
	}
	
	
	chirps :=[]Chirp{}

	for _, chirp:= range dbchirps{
		chirps = append(chirps, Chirp{
			ID: chirp.ID,
			Body: chirp.Body,
			UserID: chirp.UserID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
		})
	}
	respondWithJson(w, http.StatusOK, chirps)

	
}