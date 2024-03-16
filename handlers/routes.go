package handlers

import (
	"chi_soccer/middlewares"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func Routes() http.Handler {
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

	// all v1 routes
	router.Route("/api/v1/", func(router chi.Router) {

		router.Post("/auth/login", Login)
		router.Post("/auth/signup", Signup)

		// USER ROUTES
		// users login
		// users signup
		// this returns all the users in the db
		router.Get("/users/search", SearchUser)
		router.Route("/users", func(router chi.Router) {
			router.Use(middlewares.IsAuthorized)
			// router.Get("/", GetAllUsers)
			router.Get("/bytoken", GetUserByToken)
			router.Get("/user/{id}", GetUserById)
		})

		// PROFILE ROUTES
		router.Get("/profile/{id}", GetProfileById)
		router.Post("/profile/create", CreateProfile)
		router.Put("/profile/update/{id}", UpdateProfile)

		// FIELD ROUTES
		router.Get("/fields", getAllFields)
		router.Get("/fields/field/{id}", getFieldById)
		router.Post("/fields/field", CreateField)
		router.Put("/fields/update", UpdateField) // TODO

		// GAME ROUTES
		router.Get("/games/{user_id}", GetAllGames)
		router.Get("/games/game/{id}", GetGameById)
		router.Post("/games/game/byDateField", GetGameByDateField)
		router.Post("/games/game", CreateGame)
		router.Put("/games/update/{id}", UpdateGame)
		router.Delete("/games/delete/{id}", DeleteGame)

		// GROUP ROUTES
		// router.Get("/groups", GetAllGroups)
		router.Route("/groups", func(router chi.Router) {
			router.Get("/", GetAllGroups)
			router.Get("/group/{id}", GetGroupById)
			router.Get("/{user_id}", GetAllGroupsOfAUser)
			router.Post("/group", CreateGroup)
		})

		// REPORTS ROUTES
		router.Get("/reports/user/{user_id}", GetReportsOfUser)
		router.Get("/reports/group/{group_id}", GetReportsOfGroup)
		router.Get("/reports/game/{game_id}", GetReportsOfGame)
		router.Get("/reports/report/{id}", GetReportById)
		router.Post("/reports/report", CreateReport)
		router.Post("/reports/upload", UploadReportCSV)

		// MEMBER ROUTES
		router.Get("/members/{group_id}", GetAllMembersFromGroup)
		router.Post("/members/member", CreateMember)
	})
	// AUTH ROUTES

	return router
}
