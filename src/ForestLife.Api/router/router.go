package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/itsjoetree/forest-life/controllers"
	"net/http"
)

func Routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Post("/api/v1/auth/login", controllers.SignIn)
	router.Post("/api/v1/auth/register", controllers.SignUp)
	router.Post("/api/v1/auth/refresh", controllers.Refresh)
	router.Post("/api/v1/auth/logout", controllers.Logout)

	router.Get("/api/v1/profile", controllers.GetProfile)
	router.Post("/api/v1/posts/{id}/unlike", controllers.UnlikePost)
	router.Post("/api/v1/posts/{id}/like", controllers.LikePost)
	router.Get("/api/v1/posts", controllers.GetPosts)
	router.Get("/api/v1/posts/{id}", controllers.GetPostById)
	router.Post("/api/v1/posts", controllers.CreatePost)
	router.Put("/api/v1/posts/{id}", controllers.UpdatePost)
	router.Delete("/api/v1/posts/{id}", controllers.DeletePost)

	return router
}
