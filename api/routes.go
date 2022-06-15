package main

import (
	"chi_soccer/internal/data"
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

	router.Get("/api/users/login", app.Login)
	router.Post("/api/users/login", app.Login)
	router.Get("/api/users/all", func(w http.ResponseWriter, r *http.Request) {
		var users data.User
		all, err := users.GetAll()
		if err != nil {
			app.errorLog.Println()
			return
		}
		app.writeJSON(w, http.StatusOK, all)
	})
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

	// Create a user, this should be a POST method
	router.Post("/api/users/signup", func(w http.ResponseWriter, r *http.Request) {
		var u data.User
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		app.writeJSON(w, http.StatusOK, u)
		id, err := app.models.User.Signup(u)
		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err, http.StatusForbidden)
			app.infoLog.Println("Got back if of", id)
			newUser, _ := app.models.User.GetById(id)
			app.writeJSON(w, http.StatusOK, newUser)
		}
	})

	return router
}
