package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"main/controller"
	"main/middlewares"
	"net/http"
)

type status map[string]interface{}

func main() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Group(func(r chi.Router) {
		r.Post("/users/signup", controller.Signup())
		r.Post("/users/login", controller.Login())
	})

	router.Group(func(r chi.Router) {
		r.Use(middlewares.Authenticator)
		r.Get("/users", controller.GetUsers())
		router.Get("/users/{user_id}", controller.GetUserByID())
	})

	http.ListenAndServe(":5000", router)
}
