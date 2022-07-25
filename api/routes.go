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
	// users signup
	router.Post("/api/users/login", app.Login)
	router.Post("/api/users/signup", app.Signup)
	router.Get("/", app.GetAllUsers)
	router.Get("/api/users/search", app.SearchUser)

	// this returns all the users in the db
	router.Route("/api/users", func(router chi.Router) {
		router.Use(app.IsAuthorized)
		router.Get("/bytoken", app.GetUserByToken)
	})

	// PROFILE ROUTES
	router.Post("/api/profile/create", app.CreateProfile)
	router.Get("/api/profile/{id}", app.GetProfileById)

	// FIELD ROUTES
	router.Get("/api/fields", app.GetAllFields)
	router.Get("/api/fields/field/{id}", app.GetFieldById)
	router.Post("/api/fields/field", app.CreateField)
	router.Put("/api/fields/update", app.UpdateField)

	// GAME ROUTES
	router.Get("/api/games", app.GetAllGames)
	router.Get("/api/games/game/{id}", app.GetGameById)
	router.Post("/api/games/create", app.CreateGame)
	router.Put("/api/games/update", app.UpdateGame)

	// GROUP ROUTES
	router.Route("/api/groups", func(router chi.Router) {
		router.Use(app.IsAuthorized)
		router.Get("/", app.GetAllGroups)
		router.Get("/group/{id}", app.GetGroupById)
		router.Post("/create", app.CreateGroup)
	})

	// MEMBER ROUTES
	router.Get("/api/members/{group_id}", app.GetAllMembersFromGroup)
	router.Post("/api/members/create", app.CreateMember)

	return router
}
