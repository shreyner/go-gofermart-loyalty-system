package router

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"go-gofermart-loyalty-system/internal/auth"
	"go-gofermart-loyalty-system/internal/balance"
	"go-gofermart-loyalty-system/internal/handlers"
	"go-gofermart-loyalty-system/internal/middlewares"
	"go-gofermart-loyalty-system/internal/order"
	"go-gofermart-loyalty-system/internal/pkg/jwtauth"
	"go-gofermart-loyalty-system/internal/withdrawal"
)

func New(
	log *zap.Logger,
	authService *auth.AuthService,
	balanceService *balance.BalanceService,
	withdrawalService *withdrawal.WithdrawalService,
	orderService *order.OrderService,
	asyncProcessingOrderService *order.AsyncProcessingOrderService,
) *chi.Mux {
	log.Info("Initilize REST API")

	authHandlers := handlers.NewAuthHandlers(log, authService)
	balanceHandlers := handlers.NewBalanceHandlers(log, balanceService)
	withdrawalHandlers := handlers.NewWithdrawalHandlers(log, withdrawalService)
	orderHandlers := handlers.NewOrdersHandlers(log, orderService, asyncProcessingOrderService)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(middlewares.NewStructuredLogger(log))
	r.Use(chiMiddleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.
				With(chiMiddleware.AllowContentType("application/json")).
				Group(func(r chi.Router) {
					r.Post("/register", authHandlers.Registry)
					r.Post("/login", authHandlers.Login)
				})

			r.
				With(jwtauth.Verifier([]byte(""))).
				Group(func(r chi.Router) {
					r.Use(jwtauth.Authenticator(log))

					r.Get("/orders", orderHandlers.GetAllByUser)
					r.With(chiMiddleware.AllowContentType("text/plain")).
						Post("/orders", orderHandlers.AddOrder)

					r.Get("/withdrawals", withdrawalHandlers.GetAllByUser)

					r.Get("/balance", balanceHandlers.GetUserBalance)
					r.
						With(chiMiddleware.AllowContentType("application/json")).
						Post("/balance/withdraw", withdrawalHandlers.CreateWithdrawal)
				})
		})
	})

	// // TODO: https://github.com/swaggo/http-swagger
	// r.Get("/swagger/*", httpSwagger.Handler())

	return r
}
