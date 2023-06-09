package server

import (
	"fmt"
	"net/http"

	"github.com/alexgaudon/budgie/models"
	"github.com/alexgaudon/budgie/utils"
	"github.com/go-chi/chi/v5"
)

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func (s *APIServer) registerCategories() {
	s.Router.Route("/api/categories", func(r chi.Router) {
		r.Get("/", s.WithUser(MakeHandler(s.getCategories)))
		r.Get("/{id}", s.WithUser(MakeHandler(s.getCategory)))

		r.Post("/", s.WithUser(MakeHandler(s.createCategory)))

		r.Delete("/{id}", s.WithUser(MakeHandler(s.deleteCategory)))
	})
}

func (s *APIServer) getCategories(w http.ResponseWriter, r *http.Request) *Response {
	user := r.Context().Value(ContextKey("user")).(*models.User)

	categories, err := s.DB.Categories.Find(user.ID)
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
			"data": categories,
		},
	}
}

func (s *APIServer) getCategory(w http.ResponseWriter, r *http.Request) *Response {
	id := chi.URLParam(r, "id")

	user := r.Context().Value(ContextKey("user")).(*models.User)

	category, err := s.DB.Categories.FindOne(&models.Category{
		ID: id,
	})

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	if category.UserID != user.ID {
		return &Response{
			Status: http.StatusNotFound,
			Content: JSON{
				"error": fmt.Errorf("category not found"),
			},
		}
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"data": category,
		},
	}
}

func (s *APIServer) createCategory(w http.ResponseWriter, r *http.Request) *Response {
	ccr := &CreateCategoryRequest{}

	err := utils.DecodeBody(r, ccr)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	user := r.Context().Value(ContextKey("user")).(*models.User)

	newCategory := models.Category{
		Name:   ccr.Name,
		UserID: user.ID,
	}

	c, err := s.DB.Categories.Save(&newCategory)

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
			"data": c,
		},
	}
}

func (s *APIServer) deleteCategory(w http.ResponseWriter, r *http.Request) *Response {
	id := chi.URLParam(r, "id")
	user := r.Context().Value(ContextKey("user")).(*models.User)

	category, err := s.DB.Categories.FindOne(&models.Category{
		ID: id,
	})

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	if category.UserID != user.ID {
		return &Response{
			Status: http.StatusNotFound,
			Content: JSON{
				"error": fmt.Errorf("category not found"),
			},
		}
	}

	err = s.DB.Categories.Delete(category.ID)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	return &Response{
		Status:  http.StatusNoContent,
		Content: JSON{},
	}
}
