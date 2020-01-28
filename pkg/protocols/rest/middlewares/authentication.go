package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/bungysheep/catalogue-api/pkg/commons/contextkey"
	"github.com/bungysheep/catalogue-api/pkg/configs"
	"github.com/bungysheep/catalogue-api/pkg/models/v1/signinclaimresource"
	"github.com/dgrijalva/jwt-go"
)

// AuthenticationMiddleware - Authentication middleware
func AuthenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Applying authentication middleware.\n")

		authToken := r.Header.Get("Authorization")
		authToken = strings.TrimSpace(authToken)
		if authToken == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Missing auth token.",
			})
			return
		}

		splittedToken := strings.Split(authToken, " ")
		if len(splittedToken) != 2 {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Invalid auth token.",
			})
			return
		}

		tokenClaim := signinclaimresource.NewSignInClaimResource()
		token, err := jwt.ParseWithClaims(splittedToken[1], tokenClaim, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.TOKENSIGNKEY), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"message": "Invalid auth token.",
			})
			return
		}

		ctx := context.WithValue(r.Context(), contextkey.ClaimToken, *tokenClaim) //nolint

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
