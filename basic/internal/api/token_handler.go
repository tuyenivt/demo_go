package api

import (
	"basic/internal/store"
	"basic/internal/tokens"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type TokenHandler struct {
	tokenStore store.TokenStore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(tokenStore store.TokenStore, userStore store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{tokenStore: tokenStore, userStore: userStore, logger: logger}
}

func (th *TokenHandler) HandleCreateNewToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		th.logger.Printf("ERROR: decoding create token user request got error: %v", fmt.Errorf("%w", err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := th.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		th.logger.Printf("ERROR: GetUserByUsername got error: %v", fmt.Errorf("%w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	checkPasswordResult, err := user.PasswordHash.CheckPassword(req.Password)
	if err != nil {
		th.logger.Printf("ERROR: PasswordHash.CheckPassword got error: %v", fmt.Errorf("%w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !checkPasswordResult {
		http.Error(w, "Invalid credentials", http.StatusInternalServerError)
		return
	}

	token, err := th.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)
	if err != nil {
		th.logger.Printf("ERROR: tokenStore.CreateNewToken got error: %v", fmt.Errorf("%w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(token)
}
