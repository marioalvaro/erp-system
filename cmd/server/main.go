// cmd/server/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	fmt.Println("Starting EHRP Server on :8080...")
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("EHRP System is up and running!"))
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}