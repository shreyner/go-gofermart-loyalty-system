package router

import (
	"go-gofermart-loyalty-system/internal/handlers"
	"go-gofermart-loyalty-system/internal/middlewares"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func New(log *zap.Logger) *chi.Mux {
	log.Info("Initilize REST API")
	userH := handlers.NewUsersHandlers(log)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middlewares.NewStructuredLogger(log))
	r.Use(chiMiddleware.Recoverer)

	r.Get("/", Index)
	r.Get("/me", userH.GetMe)

	// // TODO: https://github.com/swaggo/http-swagger
	// r.Get("/swagger/*", httpSwagger.Handler())

	return r
}
