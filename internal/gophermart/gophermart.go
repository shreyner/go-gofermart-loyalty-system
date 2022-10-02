package gophermart

import (
	"context"
	"database/sql"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"go-gofermart-loyalty-system/internal/auth"
	"go-gofermart-loyalty-system/internal/balance"
	"go-gofermart-loyalty-system/internal/config"
	"go-gofermart-loyalty-system/internal/order"
	client_loyalty_points "go-gofermart-loyalty-system/internal/pkg/client-loyalty-points"
	"go-gofermart-loyalty-system/internal/pkg/database"
	"go-gofermart-loyalty-system/internal/pkg/httpserver"
	"go-gofermart-loyalty-system/internal/router"
	"go-gofermart-loyalty-system/internal/user"
	"go-gofermart-loyalty-system/internal/withdrawal"
)

func Run(log *zap.Logger, cfg *config.Config) {
	log.Info("Connection to database...")
	db, err := database.New(cfg.DBURI)

	if err != nil {
		log.Fatal("Can't connection to db", zap.Error(err))
		os.Exit(1)
	}
	defer func(db *sql.DB) {
		log.Info("Close database connection")
		err := db.Close()
		if err != nil {
			log.Error("Error to close connection db", zap.Error(err))
		}
	}(db)

	log.Info("Staring the application...")

	// Repositories
	userRepository := user.NewUserRepository(db)
	balanceRepository := balance.NewBalanceRepository(db)
	orderRepository := order.NewBalanceRepository(db)
	withdrawalRepository := withdrawal.NewWithdrawalRepository(log, db)

	// Initialize DataBase schemas
	log.Info("Start initialize database schemas ...")

	err = database.InitSchemas(
		context.Background(),
		db,
		userRepository,
		balanceRepository,
		orderRepository,
		withdrawalRepository,
	)

	log.Info("Finish initialize database schemas")

	if err != nil {
		log.Fatal("Can't initialize db schema", zap.Error(err))

		os.Exit(1)
	}

	//log.Info("Request test ...")
	//
	//log.Info(cfg.AccrualSystemAddress.String())
	//
	client := client_loyalty_points.NewClientLoyaltyPoints(log, cfg.AccrualSystemAddress.String())
	//
	//response, err := client.GetOrder(context.Background(), "12345678903")
	//log.Error("client err", zap.Error(err))
	//fmt.Println(response)
	//
	//log.Info("Request test end ...")

	// Services
	userService := user.NewUserService(userRepository)
	balanceService := balance.NewBalanceService(balanceRepository)
	orderService := order.NewOrderService(orderRepository)
	authService := auth.NewAuthService(userService, balanceService)
	withdrawalService := withdrawal.NewWithdrawalService(log, withdrawalRepository, balanceService)

	orderWorkerPool := order.NewWorkerPool(
		log,
		orderService,
		balanceService,
		client,
		5,
	)
	asyncProcessingOrderService := order.NewAsyncProcessingOrderService(orderService, orderWorkerPool)

	apiMux := router.New(
		log,
		authService,
		balanceService,
		withdrawalService,
		orderService,
		asyncProcessingOrderService,
	)

	apiServer := httpserver.NewHttpServer(log, apiMux, cfg.Address)
	log.Info("Staring rest api server...")
	apiServer.Start()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case x := <-interrupt:
		log.Info("Received a signal.", zap.String("signal", x.String()))
	case err := <-apiServer.Notify():
		log.Error("Received an error from the start rest api server", zap.Error(err))
	}

	log.Info("Stopping server...")

	if err := apiServer.Stop(context.Background()); err != nil {
		log.Error("Got an error while stopping th rest api server", zap.Error(err))
	}

	orderWorkerPool.Stop()

	log.Info("The gophermart is calling the last defers and will be stopped.")
}
