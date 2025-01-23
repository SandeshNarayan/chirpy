package main

import (
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
		
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}


func main(){
	mux := http.NewServeMux()
	apiCfg := &apiConfig{}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app",http.FileServer(http.Dir(".")) )))
	
	mux.HandleFunc("/app/", func(w http.ResponseWriter, r *http.Request){
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
        w.WriteHeader(http.StatusOK)
        _, err := w.Write([]byte("Metrics reset successfully"))
        if err!=nil{
            fmt.Println(err)
            return
        }
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