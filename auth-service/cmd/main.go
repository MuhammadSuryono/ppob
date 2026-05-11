package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yontech/ppob/auth-service/config"
	"github.com/yontech/ppob/auth-service/internal/handlers"
	"github.com/yontech/ppob/auth-service/internal/middleware"
	"github.com/yontech/ppob/auth-service/internal/repository"
	"github.com/yontech/ppob/auth-service/internal/services"
)

func main() {
	cfg := config.Load()

	db, err := config.InitDB(cfg)
	if err != nil {
		fmt.Printf("Failed to connect to database: %v\n", err)
		os.Exit(1)
	}

	redisClient := config.InitRedis(cfg)
	defer redisClient.Close()

	userRepo := repository.NewUserRepository(db)
	otpRepo := repository.NewOTPRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)

	authService := services.NewAuthService(userRepo, otpRepo, walletRepo, deviceRepo, redisClient, cfg)
	authHandler := handlers.NewAuthHandler(authService)

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/health/live", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()

		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"reason": "database connection failed",
			})
			return
		}
		if err := sqlDB.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"reason": "database ping failed",
			})
			return
		}

		if err := redisClient.Ping(ctx).Err(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status": "unavailable",
				"reason": "redis connection failed",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	api := r.Group("/api/v1")
	{
	auth := api.Group("/auth")
	{
		auth.POST("/initiate", authHandler.Initiate)
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/send-otp", middleware.RateLimitMiddleware(5), authHandler.SendOTP)
		auth.POST("/verify-otp", authHandler.VerifyOTP)
		auth.POST("/verify-password", authHandler.VerifyPassword)
		auth.POST("/verify-pin", authHandler.VerifyPINLogin)
		auth.POST("/verify-credential", authHandler.VerifyCredential)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/logout", middleware.AuthMiddleware(cfg), authHandler.Logout)
	}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.POST("/auth/change-password", authHandler.ChangePassword)
			protected.POST("/auth/change-pin", authHandler.ChangePIN)
		}
	}

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		fmt.Printf("Auth Service starting on port %s\n", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Failed to start server: %v\n", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}
	fmt.Println("Server exited")
}