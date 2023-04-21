package server

import (
	"net/http"

	"github.com/alexgaudon/budgie/storage"
	"github.com/alexgaudon/budgie/utils"
	"github.com/go-chi/chi/v5"
)

type CreateBudgetRequest struct {
	Category  string `json:"category"`
	Amount    int    `json:"amount"`
	Recurring bool   `json:"recurring"`
}

func GetBudgetById(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	budget, err := storage.DB.GetBudgetById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"data": budget,
	})
}

func GetBudgets(w http.ResponseWriter, r *http.Request) error {
	user := r.Context().Value(ContextKey("user")).(*storage.User)

	budgets, err := storage.DB.GetBudgets(user.ID)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"data": budgets,
	})
}

func CreateBudget(w http.ResponseWriter, r *http.Request) error {
	createBudgetReq := &CreateBudgetRequest{}

	err := utils.DecodeBody(r, &createBudgetReq)

	if err != nil {
		return err
	}

	user := r.Context().Value(ContextKey("user")).(*storage.User)

	budget := &storage.Budget{
		UserId:    user.ID,
		Category:  createBudgetReq.Category,
		Amount:    createBudgetReq.Amount,
		Recurring: createBudgetReq.Recurring,
	}

	id, err := storage.DB.CreateBudget(budget)

	budget.ID = id

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"data": budget,
	})

}

func DeleteBudget(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user := r.Context().Value(ContextKey("user")).(*storage.User)

	err := storage.DB.DeleteBudget(id, user.ID)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"deleted": id,
	})
}
