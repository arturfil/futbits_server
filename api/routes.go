package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// USER ROUTES
	// users login
	router.Post("/api/users/login", app.Login)
	// users signup
	router.Post("/api/users/signup", app.Signup)
	// this returns all the users in the db
	router.Route("/api/users/all", func(router chi.Router) {
		router.Use(app.IsAuthorized)
		router.Get("/", app.GetAllUsers)
	})

	// FIELD ROUTES
	// GET/all
	router.Get("/api/fields", app.GetAllFields)
	router.Post("/api/fields/create", app.CreateField)
	router.Put("/api/fields/update", app.UpdateField)

	return router
}
