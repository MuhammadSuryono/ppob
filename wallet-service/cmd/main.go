package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yontech/ppob/wallet-service/config"
	"github.com/yontech/ppob/wallet-service/internal/handlers"
	"github.com/yontech/ppob/wallet-service/internal/middleware"
	"github.com/yontech/ppob/wallet-service/internal/models"
	"github.com/yontech/ppob/wallet-service/internal/repository"
	"github.com/yontech/ppob/wallet-service/internal/services"
)

func main() {
	cfg := config.Load()

	// Tracing disabled - OTel version mismatch
	// if cfg.JaegerEndpoint != "" {
	// 	tp, err := initTracer("wallet-service", cfg.JaegerEndpoint)
	// 	if err != nil {
	// 		log.Printf("Failed to init tracer: %v", err)
	// 	} else {
	// 		defer tp.Shutdown(context.Background())
	// 	}
	// }

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&models.Wallet{}, &models.WalletEvent{}, &models.Hold{}, &models.Commission{}, &models.DailyLimit{})

	redisClient := config.InitRedis(cfg)
	defer redisClient.Close()

	walletRepo := repository.NewWalletRepository(db)
	eventRepo := repository.NewEventRepository(db)

	walletService := services.NewWalletService(walletRepo, eventRepo, redisClient, cfg)
	walletHandler := handlers.NewWalletHandler(walletService)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.MetricsMiddleware())
	r.Use(middleware.RateLimitMiddleware(60))
	r.Use(middleware.TracingMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/health/ready", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.Ping() != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		wallets := api.Group("/wallets")
		{
			wallets.GET("/me/balance", walletHandler.GetBalance)
			wallets.GET("/me/balance-events", walletHandler.GetBalanceByEvents)
			wallets.POST("/:id/hold", walletHandler.PlaceHold)
			wallets.POST("/:id/release-hold", walletHandler.ReleaseHold)
			wallets.POST("/:id/debit", walletHandler.Debit)
			wallets.POST("/:id/credit", walletHandler.Credit)
			wallets.POST("/me/topup", walletHandler.TopUp) // Mitra self top-up
			wallets.POST("/transfer", walletHandler.Transfer)
			wallets.POST("/staff/topup", walletHandler.TopUpStaff)
			wallets.GET("/me/events", walletHandler.GetEvents)
			wallets.GET("/:id/reconcile", walletHandler.Reconcile)
		}
		txHandlers := api.Group("/wallets/transactions")
		{
			txHandlers.POST("/:transaction_id/hold", walletHandler.PlaceHoldForTransaction)
			txHandlers.DELETE("/:transaction_id/hold", walletHandler.ReleaseHoldForTransaction)
			txHandlers.POST("/:transaction_id/debit", walletHandler.DebitForTransaction)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Printf("Wallet Service starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}
