package main

import (
	"avito-merch/internal/config"
	"avito-merch/internal/handler"
	"avito-merch/internal/repository"
	"avito-merch/internal/service"
	"context"
	"database/sql"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.MustLoad()

	db, err := sql.Open("postgres", cfg.DBConnString())
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		logger.Fatal("Database connection failed", zap.Error(err))
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserPostgres(db)
	merchRepo := repository.NewMerchPostgres(db)
	transRepo := repository.NewTransactionPostgres(db)
	purchaseRepo := repository.NewPurchasePostgres(db)

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret)
	buyService := service.NewBuyService(userRepo, merchRepo, purchaseRepo, transRepo)
	infoService := service.NewInfoService(userRepo, purchaseRepo, transRepo)
	sendService := service.NewSendService(userRepo, transRepo)

	// Создание обработчика
	h := handler.NewHandler(
		authService,
		buyService,
		infoService,
		sendService,
		logger,
	)

	router := gin.New()
	router.Use(gin.Recovery(), h.LoggingMiddleware())
	h.SetupRoutes(router)

	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, shutdownCancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer shutdownCancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", zap.Error(err))
	}
}
