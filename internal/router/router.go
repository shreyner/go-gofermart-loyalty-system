package router

import (
	"go-gofermart-loyalty-system/internal/auth"
	"go-gofermart-loyalty-system/internal/balance"
	"go-gofermart-loyalty-system/internal/handlers"
	"go-gofermart-loyalty-system/internal/middlewares"
	"go-gofermart-loyalty-system/internal/pkg/jwtauth"
	"go-gofermart-loyalty-system/internal/withdrawal"
	"net/http"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

func New(
	log *zap.Logger,
	authService *auth.AuthService,
	balanceService *balance.BalanceService,
	withdrawalService *withdrawal.WithdrawalService,
) *chi.Mux {
	log.Info("Initilize REST API")
	userHandlers := handlers.NewUsersHandlers(log)
	authHandlers := handlers.NewAuthHandlers(log, authService)
	balanceHandlers := handlers.NewBalanceHandlers(log, balanceService)
	withdrawalHandlers := handlers.NewWithdrawalHandlers(log, withdrawalService)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middlewares.NewStructuredLogger(log))
	r.Use(chiMiddleware.Recoverer)

	r.Get("/", Index)
	r.Get("/me", userHandlers.GetMe)

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

						r.Get("/withdrawals", withdrawalHandlers.GetAllByUser)

						r.Get("/balance", balanceHandlers.GetUserBalance)
						r.Post("/balance/withdraw", withdrawalHandlers.CreateWithdrawal)
						// TODO: for example. Remove before example
						r.Get("/me", func(rw http.ResponseWriter, r *http.Request) {
							jwtData, _ := jwtauth.JwtDataFromContext(r.Context())

							rw.WriteHeader(http.StatusOK)
							_, _ = rw.Write([]byte(jwtData.ID))
						})

					})

			})
		})

	// // TODO: https://github.com/swaggo/http-swagger
	// r.Get("/swagger/*", httpSwagger.Handler())

	return r
}
