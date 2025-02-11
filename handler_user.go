package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/crosve/golang/internal/database"
)

func (apicfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	type paramaters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)

	params := paramaters{}

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "Invalid request payload")
		return
	}

	user, err := apicfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		Name:      params.Name,
	})

	if err != nil {
		respondWithError(w, 500, "Failed to create user")
		return
	}

	respondWithJSON(w, 200, user)

}
