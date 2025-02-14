package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/crosve/golang/internal/auth"
	"github.com/crosve/golang/internal/database"
)

func (apicfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	type paramaters struct {
		Name string `json:"name"`
	}

	params := paramaters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "Error parsing JSON")
		return
	}

	user, err := apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 500, "Failed to create user")
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))

}

func (apicfg *apiConfig) handleGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)

	if err != nil {
		respondWithError(w, 401, "Invalid API Key")
		return
	}

	user, err := apicfg.DB.GetUserByAPIKey(r.Context(), apiKey)

	if err != nil {
		respondWithError(w, 404, fmt.Sprintf("User with API Key %s not found", apiKey))
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))

}
