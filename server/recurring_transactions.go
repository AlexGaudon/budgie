package server

import (
	"net/http"

	"github.com/alexgaudon/budgie/models"
	"github.com/alexgaudon/budgie/utils"
	"github.com/go-chi/chi/v5"
)

type CreateRecurringTransactionRequest struct {
	Amount        int                   `json:"amount"`
	CategoryID    string                `json:"category_id"`
	Description   models.OptionalString `json:"description"`
	Vendor        string                `json:"vendor"`
	Type          string                `json:"type"`
	UnitOfMeasure string                `json:"unit_of_measure"`
	Frequency     int                   `json:"frequency_count"`
}

func (s *APIServer) registerRecurringTransactions() {
	s.Router.Route("/api/recurringtransactions", func(r chi.Router) {
		r.Get("/", s.WithUser(MakeHandler(s.getRecurringTransactions)))

		r.Post("/", s.WithUser(MakeHandler(s.createRecurringTransaction)))
	})
}

func (s *APIServer) getRecurringTransactions(w http.ResponseWriter, r *http.Request) *Response {
	user := r.Context().Value(ContextKey("user")).(*models.User)

	t, err := s.DB.RecurringTransactions.Find(user.ID)

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

func (s *APIServer) createRecurringTransaction(w http.ResponseWriter, r *http.Request) *Response {
	user := r.Context().Value(ContextKey("user")).(*models.User)

	crtr := &CreateRecurringTransactionRequest{}

	if err := utils.DecodeBody(r, &crtr); err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	t := &models.Transaction{
		UserID:      user.ID,
		CategoryID:  crtr.CategoryID,
		Description: crtr.Description,
		Type:        crtr.Type,
		Vendor:      crtr.Vendor,
		Amount:      crtr.Amount,
	}

	rt := &models.RecurringTransaction{
		Frequency:     crtr.Frequency,
		UnitOfMeasure: crtr.UnitOfMeasure,
	}
	rt.Transaction = *t

	rt, err := s.DB.RecurringTransactions.Save(rt)

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
			"data": rt,
		},
	}
}
