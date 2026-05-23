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
	"github.com/yontech/ppob/transaction-service/config"
	"github.com/yontech/ppob/transaction-service/internal/clients"
	"github.com/yontech/ppob/transaction-service/internal/handlers"
	"github.com/yontech/ppob/transaction-service/internal/middleware"
	"github.com/yontech/ppob/transaction-service/internal/models"
	"github.com/yontech/ppob/transaction-service/internal/repository"
	"github.com/yontech/ppob/transaction-service/internal/services"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func startReconciliationCron(svc *services.ReconciliationService) {
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
				result, err := svc.ReconcileStalePending(ctx)
				if err != nil {
					log.Printf("Reconciliation error: %v", err)
				} else {
					log.Printf("Reconciliation completed: expired=%d, released=%d",
						result.ExpiredCount, result.ReleasedHoldCnt)
				}
				cancel()
			}
		}
	}()
}

func initTracer(serviceName, jaegerEndpoint string) (*sdktrace.TracerProvider, error) {
	var exporter sdktrace.SpanExporter
	var err error

	if jaegerEndpoint != "" {
		exporter, err = jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
		if err != nil {
			return nil, err
		}
	}

	res, err := resource.New(context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}

func main() {
	cfg := config.Load()

	if cfg.JaegerEndpoint != "" {
		tp, err := initTracer("transaction-service", cfg.JaegerEndpoint)
		if err != nil {
			log.Printf("Failed to init tracer: %v", err)
		} else {
			defer tp.Shutdown(context.Background())
		}
	}

	db, err := config.InitDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	db.AutoMigrate(&models.Transaction{}, &models.DailyLimit{}, &models.Commission{})

	redisClient := config.InitRedis(cfg)
	defer redisClient.Close()

	transactionRepo := repository.NewTransactionRepository(db)

	// Clients
	walletClient, err := clients.NewWalletClient(cfg.WalletGRPCAddr)
	if err != nil {
		log.Fatalf("failed to create wallet client: %v", err)
	}
	defer walletClient.Close()

	productClient, err := clients.NewProductClient(cfg.ProductGRPCAddr)
	if err != nil {
		log.Fatalf("failed to create product client: %v", err)
	}
	defer productClient.Close()

	marginService := services.NewMarginService(db, cfg, productClient)

	integrationClient := clients.NewIntegrationClient("http://integration-service:8080")

	transactionService := services.NewTransactionService(transactionRepo, marginService, redisClient, cfg, db, walletClient, productClient, integrationClient)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Commission System Wiring (Async)
	commissionService := services.NewCommissionService(db, marginService, walletClient)
	commissionWorker := services.NewCommissionWorker(redisClient, transactionRepo, marginService, commissionService)
	go commissionWorker.Start(context.Background())

	reconciliationService := services.NewReconciliationService(db, walletClient)
	startReconciliationCron(reconciliationService)

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

	api.POST("/transactions/initiate", middleware.AuthMiddleware(cfg), transactionHandler.InitiateTransaction)
	api.GET("/transactions", middleware.AuthMiddleware(cfg), transactionHandler.ListTransactions)
	api.GET("/transactions/history", middleware.AuthMiddleware(cfg), transactionHandler.GetTransactionHistory)
	api.GET("/transactions/:id", middleware.AuthMiddleware(cfg), transactionHandler.GetTransaction)
	api.GET("/transactions/by-id/:id", middleware.AuthMiddleware(cfg), transactionHandler.GetTransactionByID)
	api.POST("/transactions/:id/status", middleware.AuthMiddleware(cfg), transactionHandler.UpdateTransactionStatus)
	api.POST("/transactions/:id/cancel", middleware.AuthMiddleware(cfg), transactionHandler.CancelTransaction)

	// Reports endpoints
	api.GET("/reports", middleware.AuthMiddleware(cfg), transactionHandler.GetReports)

	r.POST("/webhook/digiflazz", transactionHandler.ProcessWebhook)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Printf("Transaction Service starting on port %s", cfg.ServerPort)
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