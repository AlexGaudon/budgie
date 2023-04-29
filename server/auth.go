package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexgaudon/budgie/config"
	"github.com/alexgaudon/budgie/models"
	"github.com/alexgaudon/budgie/utils"
	"github.com/go-chi/chi/v5"
)

type RegisterRequest struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ContextKey string

func (c ContextKey) String() string {
	return string(c)
}

func (s *APIServer) registerAuth() {
	s.Router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", MakeHandler(s.register))
		r.Post("/login", MakeHandler(s.login))

		r.Get("/whoami", MakeHandler(s.refresh))

		r.Get("/logout", s.WithUser(MakeHandler(s.logout)))
	})
}

func refreshInvalid() (*Response, error) {
	return &Response{
		Status: http.StatusUnauthorized,
		Content: JSON{
			"message": "refresh_token is invalid",
		},
	}, nil
}

func (s *APIServer) logout(w http.ResponseWriter, r *http.Request) (*Response, error) {
	expired := time.Now().Add(-time.Hour * 1)

	refreshCookie := http.Cookie{
		Name:    "refresh_token",
		Path:    "/",
		Value:   "",
		Expires: expired,
	}
	accessCookie := http.Cookie{
		Name:    "access_token",
		Path:    "/",
		Value:   "",
		Expires: expired,
	}

	http.SetCookie(w, &refreshCookie)
	http.SetCookie(w, &accessCookie)

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"message": "logged out successfully",
		},
	}, nil
}

func (s *APIServer) refresh(w http.ResponseWriter, r *http.Request) (*Response, error) {
	refresh_token, err := r.Cookie("refresh_token")

	if err != nil {
		return refreshInvalid()
	}

	tokenClaims, err := utils.GetTokenClaims(refresh_token.Value)

	if err != nil {
		return refreshInvalid()
	}

	sub, err := tokenClaims.GetSubject()
	if err != nil {
		return refreshInvalid()
	}

	setToken(w, sub, "access_token", config.GetConfig().AccessTokenExpiresIn)

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"user": sub,
		},
	}, nil
}

func (s *APIServer) login(w http.ResponseWriter, r *http.Request) (*Response, error) {
	loginRequest := LoginRequest{}

	err := utils.DecodeBody(r, &loginRequest)

	if err != nil {
		return nil, err
	}

	user, err := s.User.FindOne(&models.User{
		Username: loginRequest.Username,
	})

	if err != nil {
		return nil, err
	}

	if !user.IsPasswordValid(loginRequest.Password) {
		return &Response{
			Status: http.StatusUnauthorized,
			Content: JSON{
				"message": "Username or password is incorrect.",
			},
		}, nil
	}

	err = setToken(w, user.ID, "access_token", config.GetConfig().AccessTokenExpiresIn)
	if err != nil {
		return nil, err
	}

	err = setToken(w, user.ID, "refresh_token", config.GetConfig().RefreshTokenExpiresIn)
	if err != nil {
		return nil, err
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"message": "logged in successfully",
		},
	}, nil
}

func (s *APIServer) register(w http.ResponseWriter, r *http.Request) (*Response, error) {
	regReq := &RegisterRequest{}

	err := utils.DecodeBody(r, &regReq)

	if err != nil {
		return nil, err
	}

	if regReq.Password != regReq.PasswordConfirmation {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"message": "password and confirmation password must match",
			},
		}, nil
	}

	user, err := models.NewUser(regReq.Username, regReq.Password)

	if err != nil {
		return nil, err
	}

	existingUser, err := s.User.FindOne(user)

	if err != nil {
		return nil, err
	}

	if existingUser != nil && existingUser.ID != "" {
		return nil, fmt.Errorf("a user with this name already exists")
	}

	_, err = s.User.Save(user)

	if err != nil {
		return nil, err
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"message": "Registered successfully",
		},
	}, nil
}

func (s *APIServer) WithUser(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var access_token string
		authorization := r.Header.Get("Authorization")

		if strings.HasPrefix(authorization, "Bearer ") {
			access_token = strings.TrimPrefix(authorization, "Bearer ")
		} else {
			accessCookie, err := r.Cookie("access_token")

			if err != nil {
				writeResponse(w, http.StatusUnauthorized, JSON{})
				return
			}

			access_token = accessCookie.Value
		}

		if access_token == "" {
			writeResponse(w, http.StatusUnauthorized, JSON{})
			return
		}

		tokenClaims, err := utils.GetTokenClaims(access_token)

		if err != nil {
			writeResponse(w, http.StatusUnauthorized, JSON{})
			return
		}

		sub, ok := tokenClaims["sub"].(string)

		if !ok {
			writeResponse(w, http.StatusUnauthorized, JSON{})
			return
		}

		user, err := s.User.FindOne(&models.User{
			ID: sub,
		})

		if err != nil {
			writeResponse(w, http.StatusUnauthorized, JSON{})
			return
		}

		ctx := context.WithValue(r.Context(), ContextKey("user"), user)

		newReq := r.WithContext(ctx)
		handlerFunc(w, newReq)
	}
}

func setToken(w http.ResponseWriter, userId string, tokenName string, time time.Duration) error {
	token, err := utils.CreateToken(userId, time)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     tokenName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(time) * 60,
		Secure:   false,
		HttpOnly: true,
	})

	return nil
}
