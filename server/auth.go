package server

import (
	"context"
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

		r.Get("/me", MakeHandler(s.refresh))

		r.Get("/logout", s.WithUser(MakeHandler(s.logout)))
	})
}

func refreshInvalid() *Response {
	return &Response{
		Status: http.StatusUnauthorized,
		Content: JSON{
			"message": "refresh_token is invalid",
		},
	}
}

func (s *APIServer) logout(w http.ResponseWriter, r *http.Request) *Response {
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
	}
}

func (s *APIServer) refresh(w http.ResponseWriter, r *http.Request) *Response {
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

	err = setToken(w, sub, "access_token", config.GetConfig().AccessTokenExpiresIn)

	if err != nil {
		return &Response{
			Status:  http.StatusInternalServerError,
			Content: JSON{},
		}
	}

	user, err := s.DB.User.FindOne(&models.User{
		ID: sub,
	})

	if err != nil {
		return &Response{
			Status: http.StatusUnauthorized,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"userId":   user.ID,
			"username": user.Username,
		},
	}
}

func (s *APIServer) login(w http.ResponseWriter, r *http.Request) *Response {
	loginRequest := LoginRequest{}

	err := utils.DecodeBody(r, &loginRequest)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	user, err := s.DB.User.FindOne(&models.User{
		Username: loginRequest.Username,
	})

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	if !user.IsPasswordValid(loginRequest.Password) {
		return &Response{
			Status: http.StatusUnauthorized,
			Content: JSON{
				"message": "Username or password is incorrect.",
			},
		}
	}

	err = setToken(w, user.ID, "access_token", config.GetConfig().AccessTokenExpiresIn)
	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	err = setToken(w, user.ID, "refresh_token", config.GetConfig().RefreshTokenExpiresIn)
	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"message":  "logged in successfully",
			"userId":   user.ID,
			"username": user.Username,
		},
	}
}

func (s *APIServer) register(w http.ResponseWriter, r *http.Request) *Response {
	regReq := &RegisterRequest{}

	err := utils.DecodeBody(r, &regReq)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	if regReq.Password != regReq.PasswordConfirmation {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"message": "password and confirmation password must match",
			},
		}
	}

	user, err := models.NewUser(regReq.Username, regReq.Password)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	exists := s.DB.User.Exists(user)

	if exists {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	_, err = s.DB.User.Save(user)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"message": "Registered successfully",
		},
	}
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

		user, err := s.DB.User.FindOne(&models.User{
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
