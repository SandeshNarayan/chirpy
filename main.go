package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	_ "github.com/lib/pq"

	"github.com/SandeshNarayan/chirpy/internal/database"
	"github.com/joho/godotenv"
)





type apiConfig struct {
	fileserverHits atomic.Int32;
	dbQueries *database.Queries;
	platform string
	jwtSecret string
}






func main(){
	const filepathroot =  "."
	const port = "8080"
	
	
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == ""{
        log.Fatal("DB_URL is not set")
    }

	platform := os.Getenv("PLATFORM")
	if platform==""{
		log.Fatal("PLATFORM must be set")
	}


	db, err := sql.Open("postgres", dbURL)
	if err!=nil {
        log.Fatalf("Error opening database: %v", err)
    }
	
	dbQueries := database.New(db)

	secret:= os.Getenv("SERVER_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}

	apiCfg := &apiConfig{
		dbQueries: dbQueries,
		fileserverHits: atomic.Int32{},
		platform: platform,
		jwtSecret: secret,
	}

	mux := http.NewServeMux()
	fsHandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(filepathroot)) ))
	mux.Handle("/app/", fsHandler)
	


	mux.HandleFunc("GET /api/healthz", handlerReadiness)

	mux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	

	mux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)

	mux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieve)

	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerFindChirpById)

	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)

	mux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	mux.HandleFunc("POST /api/refresh", apiCfg.handlerRefresh)

	mux.HandleFunc("POST /api/revoke", apiCfg.handlerRevoke)

	

	


	server := &http.Server{
		Handler : mux,
		Addr : ":" + port,
	}

	log.Printf("Server started on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
	
	

}