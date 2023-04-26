package server

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type JSON = map[string]any

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHandlerFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			log.Println("ERROR: ", err.Error())
			WriteJSON(w, http.StatusBadRequest, JSON{"error": err.Error()})
		}
	}
}

func WriteInternalServerError(w http.ResponseWriter) {
	WriteJSON(w, http.StatusInternalServerError, JSON{
		"message": "there was an error processing your request",
	})
}

func WriteUnauthorized(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, JSON{
		"error": "Unauthorized",
	})
}

func ErrorForbidden(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, JSON{
		"error": "forbidden",
	})
}

func WriteJSON(w http.ResponseWriter, status int, c JSON) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if status == http.StatusNoContent {
		return nil
	}
	return json.NewEncoder(w).Encode(c)
}

func SetupServer() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Content-Type", "Set-Cookie", "Cookie"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Route("/api/transactions", func(r chi.Router) {
		r.Get("/", WithUser(makeHandlerFunc(GetTransactions)))

		r.Get("/category/{id}", WithUser(makeHandlerFunc(GetTransactionsByCategory)))

		r.Get("/{id}", WithUser(makeHandlerFunc(GetTransactionById)))

		r.Post("/", WithUser(makeHandlerFunc(CreateTransaction)))

		r.Delete("/{id}", WithUser(makeHandlerFunc(DeleteTransaction)))

		r.Put("/{id}", WithUser(makeHandlerFunc(UpdateTransaction)))
	})

	r.Route("/api/categories", func(r chi.Router) {
		r.Get("/", WithUser(makeHandlerFunc(GetCategories)))

		r.Post("/", WithUser(makeHandlerFunc(CreateCategory)))

		r.Delete("/{id}", WithUser(makeHandlerFunc(DeleteCategory)))
	})

	r.Route("/api/budgets", func(r chi.Router) {
		r.Get("/", WithUser(makeHandlerFunc(GetBudgets)))
	})

	r.Route("/api/user", func(r chi.Router) {
		r.Post("/register", makeHandlerFunc(Register))
		r.Post("/login", makeHandlerFunc(Login))
		r.Get("/me", makeHandlerFunc(RefreshAccessToken))
		r.Get("/logout", makeHandlerFunc(Logout))
	})

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "client/dist")

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(filesDir + r.URL.Path); errors.Is(err, os.ErrNotExist) {
			http.ServeFile(w, r, filepath.Join(filesDir, "index.html"))
		}
		http.ServeFile(w, r, filesDir+r.URL.Path)
	})

	return r
}
