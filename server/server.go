package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexgaudon/budgie/models"
	"github.com/alexgaudon/budgie/storage"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type APIServer struct {
	Router *chi.Mux
	User   models.UserRepo
}

func NewAPIServer(db *storage.DBStore) *APIServer {
	return &APIServer{
		Router: chi.NewRouter(),
		User:   *db.User,
	}
}

type Response struct {
	Status  int
	Content JSON
}

type JSON = map[string]any

type apiFunc = func(http.ResponseWriter, *http.Request) (*Response, error)

func writeResponse(w http.ResponseWriter, status int, c JSON) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if status == http.StatusNoContent {
		return nil
	}
	return json.NewEncoder(w).Encode(c)
}

func MakeHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp, err := fn(w, r)

		if err != nil {
			log.Println("ERROR: ", err.Error())
			writeResponse(w, http.StatusBadRequest, JSON{
				"error": err.Error(),
			})
			return
		}

		writeResponse(w, resp.Status, resp.Content)
	}
}

func (a *APIServer) ConfigureServer() {
	a.Router.Use(middleware.Logger)
	a.Router.Use(middleware.RealIP)
	a.Router.Use(middleware.RequestID)
	a.Router.Use(middleware.Recoverer)
	a.Router.Use(middleware.Compress(5))

	a.Router.Use(middleware.Timeout(60 * time.Second))

	a.Router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Content-Type", "Set-Cookie", "Cookie"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	a.registerAuth()

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "/client/dist")

	a.Router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(filesDir + r.URL.Path); errors.Is(err, os.ErrNotExist) {
			http.ServeFile(w, r, filepath.Join(filesDir, "index.html"))
		}
		http.ServeFile(w, r, filesDir+r.URL.Path)
	})

	chi.Walk(a.Router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middlewares\n", method, route, len(middlewares))
		return nil
	})
}
