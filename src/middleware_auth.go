package main

import (
	"github.com/tonge3199/go-RSS-project/internal/auth"
	"github.com/tonge3199/go-RSS-project/internal/database"
	"net/http"
)

type authedHandler func(w http.ResponseWriter, r *http.Request, user database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract API Key from Authorization header
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err.Error())
			return
		}
		// Get user by API Key from database
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusNotFound, "Couldn't get user")
			return
		}
		// Call the original handler with user context
		// 调用原始处理函数并传递用户上下文
		handler(w, r, user)
	}
}
