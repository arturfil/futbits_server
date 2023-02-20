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
	router.Get("/api/profile/{id}", controllers.GetProfileById)
	router.Post("/api/profile/create", controllers.CreateProfile)

	// FIELD ROUTES
	router.Get("/api/fields", controllers.GetAllFields)
	router.Get("/api/fields/field/{id}", controllers.GetFieldById)
	router.Post("/api/fields/field", controllers.CreateField)
	router.Put("/api/fields/update", controllers.UpdateField)

	// GAME ROUTES
	router.Get("/api/games", controllers.GetAllGames)
	router.Get("/api/games/game/{id}", controllers.GetGameById)
	router.Post("/api/games/create", controllers.CreateGame)
	router.Put("/api/games/update", controllers.UpdateGame)

	// GROUP ROUTES
	// router.Get("/api/groups", controllers.GetAllGroups)
	router.Route("/api/groups", func(router chi.Router) {
		router.Get("/", controllers.GetAllGroups)
		router.Get("/group/{id}", controllers.GetGroupById)
		router.Get("/{user_id}", controllers.GetAllGroupsOfAUser)
		router.Post("/create", controllers.CreateGroup)
	})

	// REPORTS ROUTES
	router.Get("/api/reports/{user_id}", controllers.GetReportsOfUser)
	router.Get("/api/reports/report/{id}", controllers.GetReportById)
	router.Post("/api/reports/report", controllers.CreateReport)

	// MEMBER ROUTES
	router.Get("/api/members/{group_id}", controllers.GetAllMembersFromGroup)
	router.Post("/api/members/create", controllers.CreateMember)

	return router
}
