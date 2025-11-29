package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ArtisGulbis/pokemon-companion-go-backend/middleware"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /", helloHandler)

	stack := middleware.CreateStack(
		middleware.Logging,
		middleware.Authorization,
	)

	server := http.Server{
		Addr:    ":3000",
		Handler: stack(router),
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server listening on port 3000")

	server.ListenAndServe()

}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Send "Hello Pokemon!" back to whoever made the request
	fmt.Fprintf(w, "Hello Pokemon!")
}
