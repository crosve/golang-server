package main

import (
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable was not set")
	}

	fmt.Println("PORT is set to", portString)

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", readinessHandler)
	v1Router.Get("/err", handleErr)

	router.Mount("/v1", v1Router)

	fmt.Println("Server starting on port ", portString)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	fmt.Println("Server is up and running", portString)
	err := srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
