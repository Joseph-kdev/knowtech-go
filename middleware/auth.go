package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	firebaseauth "github.com/Joseph-kdev/knowtech-go/internal/firebase"
	"github.com/Joseph-kdev/knowtech-go/response"
)

type AuthToken struct {
	UID    string
	Email  string
	Claims map[string]interface{}
}

type authTokenContextKey struct{}

func GetAuthKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")

	if val == "" {
		return "", errors.New("no authorization header found")
	}

	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("malformed auth header")
	}

	if vals[0] != "Bearer" {
		return "", errors.New("malformed first part of auth header")
	}

	return vals[1], nil
}

func MiddlewareAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := GetAuthKey(r.Header)
		if err != nil {
			response.RespondWithError(w, 401, err.Error())
			return
		}

		token, err := firebaseauth.VerifyIDToken(r.Context(), tokenString)
		if err != nil {
			response.RespondWithError(w, 401, err.Error())
			return
		}

		authToken := AuthToken{
			UID:    token.UID,
			Email:  getStringClaim(token.Claims, "email"),
			Claims: token.Claims,
		}

		ctx := context.WithValue(r.Context(), authTokenContextKey{}, authToken)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetAuthToken(r *http.Request) (AuthToken, bool) {
	authToken, ok := r.Context().Value(authTokenContextKey{}).(AuthToken)
	return authToken, ok
}

func getStringClaim(claims map[string]interface{}, key string) string {
	if claims == nil {
		return ""
	}
	if value, ok := claims[key].(string); ok {
		return value
	}
	return ""
}
