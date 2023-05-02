package server

import (
	"net/http"
	"time"

	"github.com/alexgaudon/budgie/models"
	"github.com/alexgaudon/budgie/utils"
	"github.com/go-chi/chi/v5"
)

type CreateTransactionRequest struct {
	Amount      int                   `json:"amount"`
	CategoryID  string                `json:"category_id"`
	Description models.OptionalString `json:"description"`
	Vendor      string                `json:"vendor"`
	Date        time.Time             `json:"date"`
	Type        string                `json:"type"`
}

func (s *APIServer) registerTransactions() {
	s.Router.Route("/api/transactions", func(r chi.Router) {
		r.Get("/", s.WithUser(MakeHandler(s.getTransactions)))
		r.Get("/{id}", s.WithUser(MakeHandler(s.getTransction)))

		r.Post("/", s.WithUser(MakeHandler(s.createTransaction)))

		r.Put("/{id}", s.WithUser(MakeHandler(s.updateTransaction)))

		r.Delete("/{id}", s.WithUser(MakeHandler(s.deleteTransaction)))
	})
}

func (s *APIServer) getTransactions(w http.ResponseWriter, r *http.Request) *Response {
	user := r.Context().Value(ContextKey("user")).(*models.User)

	transactions, err := s.DB.Transactions.Find(user.ID)

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
			"data": transactions,
		},
	}
}

func (s *APIServer) getTransction(w http.ResponseWriter, r *http.Request) *Response {
	id := chi.URLParam(r, "id")
	user := r.Context().Value(ContextKey("user")).(*models.User)

	transaction, err := s.DB.Transactions.FindOne(&models.Transaction{
		ID: id,
	})

	if transaction.UserID != user.ID {
		return &Response{
			Status: http.StatusNotFound,
			Content: JSON{
				"error": "transaction not found",
			},
		}
	}

	if err != nil {
		return &Response{
			Status: http.StatusInternalServerError,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"data": transaction,
		},
	}
}

func (s *APIServer) createTransaction(w http.ResponseWriter, r *http.Request) *Response {
	ctr := &CreateTransactionRequest{}

	err := utils.DecodeBody(r, &ctr)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	user := r.Context().Value(ContextKey("user")).(*models.User)

	t := &models.Transaction{
		UserID:      user.ID,
		CategoryID:  ctr.CategoryID,
		Date:        ctr.Date,
		Description: ctr.Description,
		Type:        ctr.Type,
		Vendor:      ctr.Vendor,
		Amount:      ctr.Amount,
	}

	t, err = s.DB.Transactions.Save(t)

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
			"data": t,
		},
	}
}

func (s *APIServer) updateTransaction(w http.ResponseWriter, r *http.Request) *Response {
	ctr := &CreateTransactionRequest{}

	err := utils.DecodeBody(r, &ctr)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	id := chi.URLParam(r, "id")

	user := r.Context().Value(ContextKey("user")).(*models.User)

	t, err := s.DB.Transactions.FindOne(&models.Transaction{
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

	if user.ID != t.UserID {
		return &Response{
			Status: http.StatusNotFound,
			Content: JSON{
				"error": "transaction not found",
			},
		}
	}

	tempTransaction := models.Transaction{
		ID:          t.ID,
		Amount:      ctr.Amount,
		Date:        ctr.Date,
		UserID:      t.UserID,
		CategoryID:  ctr.CategoryID,
		Vendor:      ctr.Vendor,
		Type:        ctr.Type,
		Description: ctr.Description,
	}

	updatedTransaction, err := s.DB.Transactions.Save(&tempTransaction)

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
			"data": updatedTransaction,
		},
	}
}

func (s *APIServer) deleteTransaction(w http.ResponseWriter, r *http.Request) *Response {
	id := chi.URLParam(r, "id")

	user := r.Context().Value(ContextKey("user")).(*models.User)

	t, err := s.DB.Transactions.FindOne(&models.Transaction{
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

	if user.ID != t.UserID {
		return &Response{
			Status: http.StatusNotFound,
			Content: JSON{
				"error": "transaction not found",
			},
		}
	}

	err = s.DB.Transactions.Delete(id)

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
			"deleted": id,
		},
	}
}
