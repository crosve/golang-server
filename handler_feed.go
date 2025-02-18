package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/crosve/golang/internal/database"
)

func (apicfg *apiConfig) handleCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {

	type paramaters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	params := paramaters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "Error parsing JSON")
		return
	}

	feed, err := apicfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})

	if err != nil {
		respondWithError(w, 500, "Failed to create user")
		return
	}

	respondWithJSON(w, 200, databaseFeedToFeed(feed))

}

func (apicfg *apiConfig) handleGetFeeds(w http.ResponseWriter, r *http.Request) {

	feed, err := apicfg.DB.GetFeeds(r.Context())
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to get feeds: %v", err))
	}

	respondWithJSON(w, 200, databaseFeedsToFeeds(feed))

}
