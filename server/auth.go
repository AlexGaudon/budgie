package server

import (
	"fmt"
	"net/http"

	"github.com/alexgaudon/budgie/models"
	"github.com/alexgaudon/budgie/utils"
	"github.com/go-chi/chi/v5"
)

type RegisterRequest struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (s *APIServer) registerAuth() {
	s.Router.Route("/api/user", func(r chi.Router) {
		r.Post("/register", MakeHandler(s.register))
	})
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
