package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/alexgaudon/budgie/storage"
	"github.com/alexgaudon/budgie/utils"
	"github.com/go-chi/chi/v5"
)

type CreateTransactionRequest struct {
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Amount      int       `json:"amount"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
}

func GetTransactionById(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	transaction, err := storage.DB.GetTransactionById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"data": transaction,
	})
}

func GetTransactions(w http.ResponseWriter, r *http.Request) error {
	user := r.Context().Value(ContextKey("user")).(*storage.User)
	transactions, err := storage.DB.GetTransactions(user.ID)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"data": transactions,
	})
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) error {
	createTransactionReq := &CreateTransactionRequest{}

	err := utils.DecodeBody(r, &createTransactionReq)

	if err != nil {
		return err
	}

	user := r.Context().Value(ContextKey("user")).(*storage.User)

	if createTransactionReq.Type == "" {
		createTransactionReq.Type = "expense"
	}

	transaction := &storage.Transaction{
		UserId:      user.ID,
		Description: createTransactionReq.Description,
		Category:    createTransactionReq.Category,
		Amount:      createTransactionReq.Amount,
		Date:        createTransactionReq.Date,
		Type:        createTransactionReq.Type,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	id, err := storage.DB.CreateTransaction(transaction)

	transaction.ID = id

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"data": transaction,
	})
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user := r.Context().Value(ContextKey("user")).(*storage.User)

	t, err := storage.DB.GetTransactionById(id)

	if err != nil {
		return err
	}

	if user.ID != t.UserId {
		return fmt.Errorf("unable to update transaction, not the owner")
	}

	ctr := &CreateTransactionRequest{}

	err = utils.DecodeBody(r, &ctr)

	if err != nil {
		return err
	}

	t.Amount = ctr.Amount
	t.Category = ctr.Category
	t.Description = ctr.Description
	t.Date = ctr.Date
	t.UpdatedAt = time.Now().UTC()
	t.Type = ctr.Type

	err = storage.DB.UpdateTransaction(t)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"data": t,
	})
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user := r.Context().Value(ContextKey("user")).(*storage.User)

	t, err := storage.DB.GetTransactionById(id)

	if err != nil {
		return err
	}

	if t.UserId != user.ID {
		return fmt.Errorf("unable to delete that transaction, not the owner")
	}

	err = storage.DB.DeleteTransaction(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"deleted": id,
	})
}
