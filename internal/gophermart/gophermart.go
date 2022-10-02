package gophermart

import (
	"context"
	"database/sql"
	"go-gofermart-loyalty-system/internal/auth"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"go-gofermart-loyalty-system/internal/balance"
	"go-gofermart-loyalty-system/internal/config"
	"go-gofermart-loyalty-system/internal/order"
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

	//log.Info("Start insert user")
	//userEntity := user.UserEntity{Login: "alex4"}
	//_ = userEntity.SetPassword("1235")
	//err = userRepository.Create(context.Background(), &userEntity)
	//if err != nil {
	//	if errors.Is(err, user.ErrLoginAlreadyExist) {
	//		log.Info("this error is user already exists", zap.Error(err))
	//	} else {
	//		log.Error("can't create user", zap.Error(err))
	//	}
	//} else {
	//	log.Info("Success created user")
	//	fmt.Println(userEntity)
	//}
	//log.Info("Finished insert user")
	//
	//log.Info("Find user")
	//userEntity2, err := userRepository.FindByLogin(context.Background(), "alex5")
	//
	//if err != nil {
	//	if errors.Is(err, user.ErrUserNotFound) {
	//		log.Info("Current user not found")
	//	} else {
	//		log.Error("unknown error FindByLogin", zap.Error(err))
	//	}
	//
	//} else {
	//	log.Info("User fined ")
	//	fmt.Println(userEntity2)
	//}
	//
	//if err := balanceRepository.Create(context.Background(), "51a7be8f-985e-4312-9b48-452e31c4efc9"); err != nil {
	//	log.Error("Error created", zap.Error(err))
	//} else {
	//	log.Info("Success Created")
	//}
	//
	//balanceEntity, err := balanceRepository.FindByUser(context.Background(), "51a7be8f-985e-4312-9b48-452e31c4efc9")
	//if err != nil {
	//	if errors.Is(err, balance.ErrBalanceNotFound) {
	//		log.Warn("Error not found")
	//	} else {
	//		log.Error("unknown error", zap.Error(err))
	//	}
	//} else {
	//	log.Info("finded blance")
	//	fmt.Println(balanceEntity)
	//}

	// Services
	userService := user.NewUserService(userRepository)
	balanceService := balance.NewBalanceService(balanceRepository)
	//orderService := order.NewOrderService(orderRepository)
	authService := auth.NewAuthService(userService, balanceService)
	withdrawalService := withdrawal.NewWithdrawalService(log, withdrawalRepository, balanceService)

	//orderWorkerPool := order.NewWorkerPool(log, orderService, &client_loyalty_points.ClientLoyaltyPoints{}, 5)
	//// TODO: определить порядок defer. На случай завершения connection к базе раньше чем очистится очередь
	//defer orderWorkerPool.Stop()
	//asyncProcessingOrder := order.NewAsyncProcessingOrder(orderService, orderWorkerPool)

	apiMux := router.New(
		log,
		authService,
		balanceService,
		withdrawalService,
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

	log.Info("The gophermart is calling the last defers and will be stopped.")
}
