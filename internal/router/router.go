package router

import (
	"go-gofermart-loyalty-system/internal/auth"
	"go-gofermart-loyalty-system/internal/handlers"
	"go-gofermart-loyalty-system/internal/middlewares"
	"go-gofermart-loyalty-system/internal/pkg/jwtauth"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func New(log *zap.Logger, authService *auth.AuthService) *chi.Mux {
	log.Info("Initilize REST API")
	userH := handlers.NewUsersHandlers(log)
	authHandlers := handlers.NewAuthHandlers(log, authService)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middlewares.NewStructuredLogger(log))
	r.Use(chiMiddleware.Recoverer)

	r.Get("/", Index)
	r.Get("/me", userH.GetMe)

	r.
		With(chiMiddleware.AllowContentType("application/json")).
		Route("/api", func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {

				r.Group(func(r chi.Router) {
					r.Post("/register", authHandlers.Registry)
					r.Post("/login", authHandlers.Login)
				})

				r.
					With(jwtauth.Verifier([]byte(""))).
					Group(func(r chi.Router) {

						r.Use(jwtauth.Authenticator(log))

						r.Get("/me", func(rw http.ResponseWriter, r *http.Request) {
							jwtData, _ := jwtauth.JwtDataFromContext(r.Context())

							rw.WriteHeader(http.StatusOK)
							rw.Write([]byte(jwtData.ID))
						})

					})

			})
		})

	// // TODO: https://github.com/swaggo/http-swagger
	// r.Get("/swagger/*", httpSwagger.Handler())

	return r
}
