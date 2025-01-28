package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/SandeshNarayan/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type ChirpRequest struct {
	Body string `json:"body"`
}

type apiConfig struct {
	fileserverHits atomic.Int32;
	dbQueries *database.Queries;
	platform string
}

type ErrorResponse struct{
	Error string `json:"error"`
}

type SuccessResponse struct{
	Cleaned_body string `json:"cleaned_body"`
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
		
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func respondWithError(w http.ResponseWriter, code int, msg string){
	w.WriteHeader(code)
	respondWithJson(w, code, ErrorResponse{Error: msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err:=json.NewEncoder(w).Encode(payload); err!=nil {
		fmt.Println(err)
	}
}


func main(){
	mux := http.NewServeMux()
	
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err!=nil {
        log.Fatal(err)
    }

	apiCfg := &apiConfig{
		dbQueries: database.New(db),
		platform: os.Getenv("PLATFORM"),
	}


	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(".")) )))
	
	mux.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Content-Type", "text/plain")
        w.WriteHeader(http.StatusOK)
        _, err := w.Write([]byte("Welcome to Chirpy"))
		if err!= nil{
            fmt.Println(err)
            return
        }
	})

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter,r  *http.Request){
		
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		
		_, err := w.Write([]byte("OK"))
		if err!=nil{
			fmt.Println(err)
			return
		}
	} )

	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter,r  *http.Request){
		
		w.Header().Set("Content-Type", "application/json")

		var requestBody struct {
			Email string `json:"email"`
		}
		
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err!=nil{
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
            return
		}

		user, err := apiCfg.dbQueries.CreateUser(r.Context(), requestBody.Email)
		if err!=nil{
            respondWithError(w, http.StatusInternalServerError, err.Error())
            return
        }

		response:= map[string]interface{}{
			"id": user.ID,
			"email": user.Email,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		}
		

		w.WriteHeader(http.StatusCreated)
		respondWithJson(w, http.StatusCreated, response)
	} )

	mux.HandleFunc("/", func(w http.ResponseWriter,r  *http.Request){
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		
		body, err := io.ReadAll(r.Body)
		if err !=nil{
			fmt.Println(err)
            return  // Return early if there was an error reading the request body.
		}
		_, err = w.Write([]byte(body))
		if err !=nil{
			fmt.Println(err)
			return
		}
	} )

	

	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		metrics:= fmt.Sprintf(
			`<html>
			  <body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			  </body>
			</html>
			`, apiCfg.fileserverHits.Load())
		_, err := w.Write([]byte(metrics))
		if err!= nil{
            fmt.Println(err)
            return
        }
	})

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request){
		
		apiCfg.fileserverHits.Store(0)

		if apiCfg.platform!="dev"{
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		//delete all users
		if err:= apiCfg.dbQueries.DeleteAllUsers(r.Context()); err!=nil{
			http.Error(w,fmt.Sprintf("Failed to reset database: %v",err.Error()), http.StatusInternalServerError )
			return
		}
		
        w.WriteHeader(http.StatusOK)
        _, err := w.Write([]byte("Metrics reset successfully"))
        if err!=nil{
            fmt.Println(err)
            return
        }
	})

	

	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request){
		
		
		var chirp struct{
			Body string `json:"body"`
			UserID string `json:"user_id"`
		}
		err:= json.NewDecoder(r.Body).Decode(&chirp)
		if err!= nil{
			respondWithError(w, http.StatusBadRequest, "Invalid JSON body")
			return
		}

		new_body :=""
		words := strings.Split(chirp.Body, " ") 
		for i, word:= range words {
			if strings.ToLower(word) == "kerfuffle" || strings.ToLower(word) == "sharbert" || strings.ToLower(word) == "fornax"{
                words[i] = "****"			
            }
			
		}
		new_body = strings.Join(words, " ")


		body := strings.TrimSpace(new_body)

		if body==""{
			respondWithError(w, http.StatusBadRequest, "`body` field must not be empty")

			return
		}

		

		if len(body)>140{
			respondWithError(w, http.StatusBadRequest, "chirp is too long")

			return
		}

		if chirp.UserID =="" || len(chirp.UserID)!=36{
			respondWithError(w, http.StatusBadRequest, "user_id must be a valid UUID")
			return 
		}

		userId , err:=uuid.Parse(chirp.UserID)
		if err!= nil{
            respondWithError(w, http.StatusBadRequest, "user_id must be a valid UUID")
            return 
        }
		createdChirp, err := apiCfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
			Body: body,
			UserID: userId,
		})

		if err!=nil{
			respondWithError(w, http.StatusInternalServerError, err.Error())
            return
		}

		response:= map[string]interface{}{
			"id": createdChirp.ID,
            "body": createdChirp.Body,
            "user_id": createdChirp.UserID,
            "created_at": createdChirp.CreatedAt,
            "updated_at": createdChirp.UpdatedAt,
		}


		respondWithJson(w, http.StatusCreated, response)


	})

	server := &http.Server{
		Handler : mux,
		Addr : ":8080",
	}

	fmt.Println("Server started on port 8080")
	if err:= server.ListenAndServe(); err!=nil{
		fmt.Printf("Error starting server: %v\n", err)
	}
	

}