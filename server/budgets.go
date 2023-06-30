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

type BudgetWithUtilization struct {
	*models.Budget
	Utilization int `json:"utilization"`
}

func (s *APIServer) registerBudgets() {
	s.Router.Route("/api/budgets", func(r chi.Router) {
		r.Get("/", s.WithUser(MakeHandler(s.getBudgets)))
		r.Get("/{id}", s.WithUser(MakeHandler(s.getBudget)))
		r.Get("/utilization/{period}", s.WithUser(MakeHandler(s.getBudgetsWithUtilization)))

		r.Post("/", s.WithUser(MakeHandler(s.createBudget)))

		r.Delete("/{id}", s.WithUser(MakeHandler(s.deleteBudget)))

		r.Get("/copy-last-period-budgets", s.WithUser(MakeHandler(s.copyLastPeriodsBudgets)))
	})
}

func (s *APIServer) copyLastPeriodsBudgets(w http.ResponseWriter, r *http.Request) *Response {
	lastPeriod, err := time.Parse("2006-01", time.Now().AddDate(0, -1, 0).Format("2006-01"))
	if err != nil {
		return &Response{Status: http.StatusBadRequest, Content: JSON{
			"error": err.Error(),
		}}
	}

	user := r.Context().Value(ContextKey("user")).(*models.User)

	err = s.copyBudgetsFromPeriod(user.ID, lastPeriod)

	if err != nil {
		return &Response{Status: http.StatusBadRequest, Content: JSON{
			"error": err.Error(),
		}}
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"data": "ok",
		},
	}
}

func (s *APIServer) copyBudgetsFromPeriod(userId string, period time.Time) error {
	budgets, err := s.getBudgetsForPeriod(userId, period)

	if err != nil {
		return err
	}
	if len(budgets) > 0 {
		for _, budget := range budgets {
			period, err := time.Parse("2006-01", time.Now().Format("2006-01"))

			if err != nil {
				return err
			}

			budget.Period = period
			budget.Category = budget.CategoryID
			budget.ID = ""

			_, err = s.DB.Budgets.Save(budget)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *APIServer) getBudgetsForPeriod(userId string, period time.Time) ([]*models.Budget, error) {
	budgets, err := s.DB.Budgets.Find(userId)

	if err != nil {
		return nil, err
	}

	filteredBudgets := s.DB.Budgets.Filter(budgets, func(b *models.Budget) bool {
		return b.Period.Equal(period)
	})

	return filteredBudgets, nil
}

func (s *APIServer) getBudgetsWithUtilization(w http.ResponseWriter, r *http.Request) *Response {
	periodString := chi.URLParam(r, "period")

	period, err := time.Parse("2006-01", periodString)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	user := r.Context().Value(ContextKey("user")).(*models.User)

	budgets, err := s.getBudgetsForPeriod(user.ID, period)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	allTransactions, err := s.DB.Transactions.Find(user.ID)

	if err != nil {
		return &Response{
			Status: http.StatusBadRequest,
			Content: JSON{
				"error": err.Error(),
			},
		}
	}

	budgetsWithUtil := []*BudgetWithUtilization{}

	for _, budget := range budgets {
		transactions := s.DB.Transactions.Filter(allTransactions, func(t *models.Transaction) bool {
			return t.Date.Month() == period.Month() && t.Date.Year() == period.Year() && t.CategoryID == budget.CategoryID
		})

		tranSum := sumTransactions(transactions)

		budgetWithUtil := &BudgetWithUtilization{
			Budget:      budget,
			Utilization: tranSum,
		}

		budgetsWithUtil = append(budgetsWithUtil, budgetWithUtil)
	}

	return &Response{
		Status: http.StatusOK,
		Content: JSON{
			"data": budgetsWithUtil,
		},
	}
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

func sumTransactions(transactions []*models.Transaction) int {
	sum := 0
	for _, t := range transactions {
		sum += t.Amount
	}
	return sum
}
