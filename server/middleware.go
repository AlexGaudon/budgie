package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/alexgaudon/budgie/storage"
	"github.com/alexgaudon/budgie/utils"
)

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

func writeForbidden(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, JSON{
		"message": "token is invalid or session has expired",
	})
}

func WithUser(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var access_token string

		authorization := r.Header.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			access_token = strings.TrimPrefix(authorization, "Bearer ")
		} else {
			accessCookie, err := r.Cookie("access_token")

			if err != nil {
				WriteUnauthorized(w)
				return
			}

			access_token = accessCookie.Value
		}

		if access_token == "" {
			WriteUnauthorized(w)
			return
		}

		tokenClaims, err := utils.GetTokenClaims(access_token)

		if err != nil {
			writeForbidden(w)
			return
		}

		sub, ok := tokenClaims["sub"].(string)

		if !ok {
			writeForbidden(w)
			return
		}

		user, err := storage.DB.GetUserById(sub)

		if err != nil {
			writeForbidden(w)
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey("user"), user)

		newReq := r.WithContext(ctx)
		handlerFunc(w, newReq)
	}
}
