package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexgaudon/budgie/config"
	"github.com/alexgaudon/budgie/storage"
	"github.com/alexgaudon/budgie/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	LoginRequest
	PasswordConfirm string `json:"password_confirm"`
}

func RefreshAccessToken(w http.ResponseWriter, r *http.Request) error {
	refresh_token, err := r.Cookie("refresh_token")

	if err != nil {
		return fmt.Errorf("refresh token not found")
	}

	tokenClaims, err := utils.GetTokenClaims(refresh_token.Value)

	if err != nil {
		return err
	}

	sub, ok := tokenClaims["sub"].(string)

	if !ok {
		return fmt.Errorf("unable to refresh token")
	}

	setAccessToken(w, sub)

	user, err := storage.DB.GetUserById(sub)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"userId":   user.ID,
		"username": user.Username,
	})
}

func Logout(w http.ResponseWriter, r *http.Request) error {
	expired := time.Now().Add(-time.Hour * 24)

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

	return WriteJSON(w, http.StatusOK, JSON{
		"message": "logged out",
	})
}

func setAccessToken(w http.ResponseWriter, userId string) error {
	config := config.GetConfig()

	accessToken, err := utils.CreateToken(userId, config.AccessTokenExpiresIn)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   config.AccessTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
	})

	return nil
}

func Login(w http.ResponseWriter, r *http.Request) error {
	loginRequest := LoginRequest{}

	err := utils.DecodeBody(r, &loginRequest)

	if err != nil {
		return err
	}

	user, err := storage.DB.GetUserByUsername(loginRequest.Username)

	if err != nil {
		return err
	}

	if !user.ValidPassword(loginRequest.Password) {
		return fmt.Errorf("incorrect username or password")
	}

	config := config.GetConfig()

	err = setAccessToken(w, user.ID)
	if err != nil {
		WriteInternalServerError(w)
		return err
	}

	refreshToken, err := utils.CreateToken(user.ID, config.RefreshTokenExpiresIn)

	if err != nil {
		WriteInternalServerError(w)
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   config.RefreshTokenMaxAge * 60,
		Secure:   false,
		HttpOnly: true,
	})

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"userId":   user.ID,
		"username": user.Username,
	})
}

func Register(w http.ResponseWriter, r *http.Request) error {
	registerRequest := RegisterRequest{}

	err := utils.DecodeBody(r, &registerRequest)

	if err != nil {
		return err
	}

	_, err = storage.DB.GetUserByUsername(registerRequest.Username)
	if err == nil {
		return fmt.Errorf("a user with this name already exists")
	}

	if registerRequest.Password != registerRequest.PasswordConfirm {
		return fmt.Errorf("password and password confirmation do not match")
	}

	user, err := storage.NewUser(registerRequest.Username, registerRequest.Password)
	if err != nil {
		return fmt.Errorf("error registering")
	}

	id, err := storage.DB.CreateUser(user)

	if err != nil {
		return fmt.Errorf("error registering user")
	}

	_ = storage.DB.InsertDefaultCategories(id)

	return WriteJSON(w, http.StatusOK, JSON{
		"data": "registered user",
		"id":   id,
	})
}
