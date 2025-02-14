package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/crosve/golang/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	godotenv.Load()

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT environment variable was not set")
	}

	//get connection to the database url
	dbURL := os.Getenv("DB_URL")

	if dbURL == "" {
		log.Fatal("DB_URL environment variable was not set")
	}

	conn, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatal(err)
	}

	apiConfig := apiConfig{
		DB: database.New(conn),
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
	v1Router.Get("/users", apiConfig.handleGetUser)
	v1Router.Post("/users", apiConfig.handleCreateUser)

	router.Mount("/v1", v1Router)

	fmt.Println("Server starting on port ", portString)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	fmt.Println("Server is up and running", portString)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
