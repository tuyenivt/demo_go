package middleware

import (
	"basic/internal/store"
	"basic/internal/tokens"
	"context"
	"net/http"
	"strings"
)

type UserMiddleware struct {
	userStore store.UserStore
}

func NewUserMiddleware(userStore store.UserStore) *UserMiddleware {
	return &UserMiddleware{userStore: userStore}
}

const UserContextKey = "userKey"

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("missing user in request")
	}
	return user
}

func (um *UserMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}
		token := headerParts[1]
		user, err := um.userStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		if user == nil {
			http.Error(w, "Token expired or invalid", http.StatusUnauthorized)
			return
		}
		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddleware) RequireUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)

		if user.IsAnonymous() {
			http.Error(w, "You must be logged in to access this route", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
