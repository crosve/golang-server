package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/crosve/golang/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handleCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	FeedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to create feed follow: %v", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(FeedFollow))

}

func (apiCfg *apiConfig) handleGetFeedFOllows(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollow, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)

	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Failed to get feed follows: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseFeedFollowToFeedFollowList(feedFollow))

}

func (apiConfig *apiConfig) handlerDeleteFeedFollower(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowIdStr := chi.URLParam(r, "feedFollowID")

	feedFollowID, err := uuid.Parse(feedFollowIdStr)

	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Invalid feed follow ID: %v", err))
		return
	}

	err = apiConfig.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Couldn't delete feed error %v", err))
		return
	}

	respondWithJSON(w, 200, map[string]string{"message": "Feed follow deleted"})

}
