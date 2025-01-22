package main

import (
	"fmt"
	"net/http"
)
func main(){
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")) )
	

	server := &http.Server{
		Handler : mux,
		Addr : ":8080",
	}

	if err:= server.ListenAndServe(); err!=nil{
		fmt.Printf("Error starting server: %v\n", err)
	}
	fmt.Println("Server started on port 8080")

}