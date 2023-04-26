package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexgaudon/budgie/storage"
	"github.com/alexgaudon/budgie/utils"
	"github.com/go-chi/chi/v5"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func GetCategories(w http.ResponseWriter, r *http.Request) error {
	user := r.Context().Value(ContextKey("user")).(*storage.User)

	categories, err := storage.DB.GetCategories(user.ID)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"data": categories,
	})
}

func CreateCategory(w http.ResponseWriter, r *http.Request) error {
	ccr := &CreateCategoryRequest{}

	err := utils.DecodeBody(r, &ccr)

	if err != nil {
		return nil
	}

	user := r.Context().Value(ContextKey("user")).(*storage.User)

	category := &storage.Category{
		Name:   ccr.Name,
		UserID: user.ID,
	}

	id, err := storage.DB.CreateCategory(category)

	category.ID = id
	category.CreatedAt = time.Now().UTC()
	category.UpdatedAt = time.Now().UTC()

	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, JSON{
		"data": category,
	})
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user := r.Context().Value(ContextKey("user")).(*storage.User)

	category, err := storage.DB.GetCategoryById(user.ID)

	if err != nil {
		return err
	}

	if category.UserID != user.ID {
		return fmt.Errorf("not authorized")
	}

	err = storage.DB.DeleteCategory(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"deleted": id,
	})
}
