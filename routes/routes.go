package routes

import (
	"fmt"
	"net/http"

	"github.com/EnesAybeyR/chi-postgresql-jwt.git/handlers"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/mdware"
	"github.com/EnesAybeyR/chi-postgresql-jwt.git/service"
	"github.com/go-chi/chi/v5"
)

func GetRoutes() *chi.Mux {
	r := chi.NewRouter()

	r.Use(mdware.ZapLogger)
	r.Get("/users", service.GetUsers)
	r.Post("/user", service.CreateUser)
	r.Get("/user/{id}", service.GetUser)

	r.Post("/register", handlers.RegisterHandler)
	r.Post("/login", handlers.LoginHandler)
	r.Post("/refresh", handlers.RefreshHandler)
	r.Post("/logout", handlers.LogoutHandler)
	r.Group(func(r chi.Router) {
		r.Use(mdware.JwtAuth)
		r.Get("/prof", func(w http.ResponseWriter, r *http.Request) {
			uid := r.Context().Value(mdware.ContextUserId).(uint)
			fmt.Fprintf(w, "Hello user %d", uid)
		})
	})
	return r
}
