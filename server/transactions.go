package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alexgaudon/budgie/storage"
	"github.com/go-chi/chi/v5"
)

type CreateTransactionRequest struct {
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Amount      int       `json:"amount"`
	Date        time.Time `json:"date"`
}

func DecodeBody(r *http.Request, rv any) error {
	if err := json.NewDecoder(r.Body).Decode(&rv); err != nil {
		return err
	}

	return nil
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

	err := DecodeBody(r, &createTransactionReq)

	if err != nil {
		return err
	}

	user := r.Context().Value(ContextKey("user")).(*storage.User)

	transaction := &storage.Transaction{
		UserId:      user.ID,
		Description: createTransactionReq.Description,
		Category:    createTransactionReq.Category,
		Amount:      createTransactionReq.Amount,
		Date:        createTransactionReq.Date,
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

func DeleteTransaction(w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	user := r.Context().Value(ContextKey("user")).(*storage.User)

	err := storage.DB.DeleteTransaction(id, user.ID)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, JSON{
		"deleted": id,
	})
}
