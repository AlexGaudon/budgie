package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexgaudon/budgie/models"
	"github.com/alexgaudon/budgie/utils"
	"github.com/go-chi/chi/v5"
)

type CreateBudgetRequest struct {
	Category string    `json:"category"`
	Amount   int       `json:"amount"`
	Period   time.Time `json:"period"`
}

func (s *APIServer) registerBudgets() {
	s.Router.Route("/api/budgets", func(r chi.Router) {
		r.Get("/", s.WithUser(MakeHandler(s.getBudgets)))
		r.Get("/{id}", s.WithUser(MakeHandler(s.getBudget)))

		r.Post("/", s.WithUser(MakeHandler(s.createBudget)))

		r.Delete("/{id}", s.WithUser(MakeHandler(s.deleteBudget)))
	})
}

func (s *APIServer) getBudgets(w http.ResponseWriter, r *http.Request) *Response {
	user := r.Context().Value(ContextKey("user")).(*models.User)

	budgets, err := s.DB.Budgets.Find(user.ID)

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
			"data": budgets,
		},
	}
}

func (s *APIServer) getBudget(w http.ResponseWriter, r *http.Request) *Response {
	id := chi.URLParam(r, "id")
	user := r.Context().Value(ContextKey("user")).(*models.User)

	budget, err := s.DB.Budgets.FindOne(&models.Budget{
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

	if budget.UserID != user.ID {
		return &Response{
			Status: http.StatusNotFound,
			Content: JSON{
				"error": fmt.Errorf("budget not found"),
			},
		}
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"data": budget,
		},
	}
}

func (s *APIServer) createBudget(w http.ResponseWriter, r *http.Request) *Response {
	cbr := &CreateBudgetRequest{}

	err := utils.DecodeBody(r, cbr)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	user := r.Context().Value(ContextKey("user")).(*models.User)

	newBudget := models.Budget{
		UserID:   user.ID,
		Category: cbr.Category,
		Amount:   cbr.Amount,
		Period:   cbr.Period,
	}

	b, err := s.DB.Budgets.Save(&newBudget)

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
			"data": b,
		},
	}
}

func (s *APIServer) deleteBudget(w http.ResponseWriter, r *http.Request) *Response {
	id := chi.URLParam(r, "id")
	user := r.Context().Value(ContextKey("user")).(*models.User)

	budget, err := s.DB.Budgets.FindOne(&models.Budget{
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

	if budget.UserID != user.ID {
		return &Response{
			Status: http.StatusNotFound,
			Content: JSON{
				"error": fmt.Errorf("budget not found"),
			},
		}
	}

	err = s.DB.Budgets.Delete(budget.ID)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	return &Response{
		Status:  http.StatusOK,
		Content: JSON{},
	}

}
