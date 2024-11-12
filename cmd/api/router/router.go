package router

import (
	"github.com/Nel-sokchhunly/go-rest-api/cmd/api/resource/book"
	"github.com/Nel-sokchhunly/go-rest-api/cmd/api/resource/health"
	validatorUtil "github.com/Nel-sokchhunly/go-rest-api/util/validator"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"gorm.io/gorm"
)

func New(db *gorm.DB) *chi.Mux {
	router := chi.NewRouter()
	validator := validatorUtil.New()

	router.Get("/health", health.Read)

	bookRoutes := BookRoutes(db, validator)

	router.Route("/api/v1", func(r chi.Router) {
		r.Mount("/books", bookRoutes)
	})

	return router
}

func BookRoutes(db *gorm.DB, validator *validator.Validate) *chi.Mux {
	r := chi.NewRouter()
	bookApi := book.New(db, validator)

	r.Get("/", bookApi.List)
	r.Post("/", bookApi.Create)
	r.Get("/{id}", bookApi.Read)
	r.Put("/{id}", bookApi.Update)
	r.Delete("/{id}", bookApi.Delete)

	return r
}
