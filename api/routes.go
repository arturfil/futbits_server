package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)

	router.Get("/api/users/login", app.Login)

	return router
}
