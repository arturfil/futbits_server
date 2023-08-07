package routes

import (
	"chi_soccer/controllers"
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
	router.Post("/api/v1/auth/login", controllers.Login)
	router.Post("/api/v1/auth/signup", controllers.Signup)

	// USER ROUTES
	// users login
	// users signup
	// this returns all the users in the db
	router.Get("/api/v1/users/search", controllers.SearchUser)
	router.Route("/api/v1/users", func(router chi.Router) {
		router.Use(middlewares.IsAuthorized)
		// router.Get("/", controllers.GetAllUsers)
		router.Get("/bytoken", controllers.GetUserByToken)
		router.Get("/user/{id}", controllers.GetUserById)
	})

	// PROFILE ROUTES
	router.Get("/api/v1/profile/{id}", controllers.GetProfileById)
	router.Post("/api/v1/profile/create", controllers.CreateProfile)
	router.Put("/api/v1/profile/update/{id}", controllers.UpdateProfile)

	// FIELD ROUTES
	router.Get("/api/v1/fields", controllers.GetAllFields)
	router.Get("/api/v1/fields/field/{id}", controllers.GetFieldById)
	router.Post("/api/v1/fields/field", controllers.CreateField)
	router.Put("/api/v1/fields/update", controllers.UpdateField) // TODO

	// GAME ROUTES
	router.Get("/api/v1/games", controllers.GetAllGames)
	router.Get("/api/v1/games/game/{id}", controllers.GetGameById)
	router.Post("/api/v1/games/game/byDateField", controllers.GetGameByDateField)
	router.Post("/api/v1/games/game", controllers.CreateGame)
	router.Put("/api/v1/games/update", controllers.UpdateGame)

	// GROUP ROUTES
	// router.Get("/api/v1/groups", controllers.GetAllGroups)
	router.Route("/api/v1/groups", func(router chi.Router) {
		router.Get("/", controllers.GetAllGroups)
		router.Get("/group/{id}", controllers.GetGroupById)
		router.Get("/{user_id}", controllers.GetAllGroupsOfAUser)
		router.Post("/group", controllers.CreateGroup)
	})

	// REPORTS ROUTES
	router.Get("/api/v1/reports/{user_id}", controllers.GetReportsOfUser)
	router.Get("/api/v1/reports/report/{id}", controllers.GetReportById)
	router.Post("/api/v1/reports/report", controllers.CreateReport)
	router.Post("/api/v1/reports/upload", controllers.UploadReportCSV)

	// MEMBER ROUTES
	router.Get("/api/v1/members/{group_id}", controllers.GetAllMembersFromGroup)
	router.Post("/api/v1/members/member", controllers.CreateMember)

	return router
}
