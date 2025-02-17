package main

import (
	"net/http"

	"github.com/crosve/golang/internal/auth"
	"github.com/crosve/golang/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiConfig *apiConfig) authMiddleware(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)

		if err != nil {
			respondWithError(w, 401, "Error getting API Key: "+err.Error())
			return
		}

		user, err := apiConfig.DB.GetUserByAPIKey(r.Context(), apiKey)

		if err != nil {
			respondWithError(w, 404, "User with API Key "+apiKey+" not found")
			return
		}

		handler(w, r, user)
	}
}
