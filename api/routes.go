package main

import (
	"encoding/json"
	"fmt"
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

	// users login
	router.Post("/api/users/login", app.Login)
	// users signup
	router.Post("/api/users/signup", app.Signup)
	// this returns all the users in the db
	router.Route("/api/users/all", func(router chi.Router) {
		router.Use(app.IsAuthorized)
		router.Get("/", app.GetAllUsers)

	})

	// route to test if the server is working
	router.Get("/api/test", func(w http.ResponseWriter, r *http.Request) {
		type Test struct {
			Message string `json:"msg"`
		}
		msg := &Test{
			Message: "Hello",
		}
		m, err := json.Marshal(msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		// app.writeJSON(w, http.StatusOK, "Test Updated")
		app.writeJSON(w, http.StatusOK, &m)
	})

	return router
}
