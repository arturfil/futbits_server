package routes

import (
	"chi_soccer/handlers"
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

	// AUTH ROUTES
	router.Post("/api/v1/auth/login", handlers.Login)
	router.Post("/api/v1/auth/signup", handlers.Signup)

	// USER ROUTES
	// users login
	// users signup
	// this returns all the users in the db
	router.Get("/api/v1/users/search", handlers.SearchUser)
	router.Route("/api/v1/users", func(router chi.Router) {
		router.Use(middlewares.IsAuthorized)
		// router.Get("/", handlers.GetAllUsers)
		router.Get("/bytoken", handlers.GetUserByToken)
		router.Get("/user/{id}", handlers.GetUserById)
	})

	// PROFILE ROUTES
	router.Get("/api/v1/profile/{id}", handlers.GetProfileById)
	router.Post("/api/v1/profile/create", handlers.CreateProfile)
	router.Put("/api/v1/profile/update/{id}", handlers.UpdateProfile)

	// FIELD ROUTES
	router.Get("/api/v1/fields", handlers.GetAllFields)
	router.Get("/api/v1/fields/field/{id}", handlers.GetFieldById)
	router.Post("/api/v1/fields/field", handlers.CreateField)
	router.Put("/api/v1/fields/update", handlers.UpdateField) // TODO

	// GAME ROUTES
	router.Get("/api/v1/games/{user_id}", handlers.GetAllGames)
	router.Get("/api/v1/games/game/{id}", handlers.GetGameById)
	router.Post("/api/v1/games/game/byDateField", handlers.GetGameByDateField)
	router.Post("/api/v1/games/game", handlers.CreateGame)
	router.Put("/api/v1/games/update/{id}", handlers.UpdateGame)
    router.Delete("/api/v1/games/delete/{id}", handlers.DeleteGame)

	// GROUP ROUTES
	// router.Get("/api/v1/groups", handlers.GetAllGroups)
	router.Route("/api/v1/groups", func(router chi.Router) {
		router.Get("/", handlers.GetAllGroups)
		router.Get("/group/{id}", handlers.GetGroupById)
		router.Get("/{user_id}", handlers.GetAllGroupsOfAUser)
		router.Post("/group", handlers.CreateGroup)
	})

	// REPORTS ROUTES
	router.Get("/api/v1/reports/user/{user_id}", handlers.GetReportsOfUser)
	router.Get("/api/v1/reports/group/{group_id}", handlers.GetReportsOfGroup)
	router.Get("/api/v1/reports/game/{game_id}", handlers.GetReportsOfGame)
	router.Get("/api/v1/reports/report/{id}", handlers.GetReportById)
	router.Post("/api/v1/reports/report", handlers.CreateReport)
	router.Post("/api/v1/reports/upload", handlers.UploadReportCSV)

	// MEMBER ROUTES
	router.Get("/api/v1/members/{group_id}", handlers.GetAllMembersFromGroup)
	router.Post("/api/v1/members/member", handlers.CreateMember)

	return router
}
