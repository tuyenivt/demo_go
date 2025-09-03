package api

import (
	"basic/internal/store"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/go-chi/chi/v5"
)

type registerUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{userStore: userStore, logger: logger}
}

func (uh *UserHandler) validateRegisterUserRequest(req *registerUserRequest) error {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return errors.New("'username' is required")
	}
	if len(req.Username) > 50 {
		return errors.New("'username' cannot be greater than 50 characters")
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		return errors.New("'email' is required")
	}
	if len(req.Email) > 255 {
		return errors.New("'email' cannot be greater than 255 characters")
	}
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errors.New("'email' is invalid")
	}
	if req.Password == "" {
		return errors.New("'password' is required")
	}
	return nil
}

func (uh *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {
	var req registerUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		uh.logger.Printf("ERROR: decoding register user request got error: %v", fmt.Errorf("%w", err))
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = uh.validateRegisterUserRequest(&req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid request data: %v", err), http.StatusBadRequest)
		return
	}

	user := &store.User{Username: req.Username, Email: req.Email, Bio: req.Bio}
	err = user.PasswordHash.SetPassword(req.Password)
	if err != nil {
		uh.logger.Printf("ERROR: hashing password got error: %v", fmt.Errorf("%w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = uh.userStore.CreateUser(user)
	if err != nil {
		uh.logger.Printf("ERROR: create user got error: %v", fmt.Errorf("%w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (uh *UserHandler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) {
	paramsUsername := strings.TrimSpace(chi.URLParam(r, "username"))
	if paramsUsername == "" {
		http.Error(w, "'username' param is required", http.StatusBadRequest)
		return
	}

	user, err := uh.userStore.GetUserByUsername(paramsUsername)
	if err != nil {
		uh.logger.Printf("ERROR: get user got error: %v", fmt.Errorf("%w", err))
		http.Error(w, "Can't fetch user data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
