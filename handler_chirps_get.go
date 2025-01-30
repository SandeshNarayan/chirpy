package main

import "net/http"

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request ){
		
	dbchirps, err:= cfg.dbQueries.GetAllChirps(r.Context())
	if err!=nil{
		respondWithError(w, http.StatusInternalServerError, err.Error(), err)
		return
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