package middlewares

import (
	"context"
	"encoding/json"
	"main/helpers"
	"net/http"
)

type status map[string]interface{}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientToken := r.Header.Get("token")
		if clientToken == "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": "no authorization header provided"})
			return
		}

		claims, msg := helpers.ValidateToken(clientToken)
		if msg != "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(status{"error": msg})
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "first_name", claims.FirstName)
		ctx = context.WithValue(ctx, "last_name", claims.LastName)
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "user_type", claims.UserType)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
