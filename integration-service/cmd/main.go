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
	"github.com/yontech/ppob/integration-service/config"
	"github.com/yontech/ppob/integration-service/internal/handlers"
	"github.com/yontech/ppob/integration-service/internal/middleware"
	"github.com/yontech/ppob/integration-service/internal/models"
	"github.com/yontech/ppob/integration-service/internal/repository"
	"github.com/yontech/ppob/integration-service/internal/services"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

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
		tp, err := initTracer("integration-service", cfg.JaegerEndpoint)
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

	db.AutoMigrate(&models.IntegrationLog{}, &models.ProviderConfig{})

	redisClient := config.InitRedis(cfg)
	defer redisClient.Close()

	integrationLogRepo := repository.NewIntegrationLogRepository(db)
	providerConfigRepo := repository.NewProviderConfigRepository(db)

	integrationService := services.NewIntegrationService(integrationLogRepo, providerConfigRepo, redisClient, cfg)

	digiflazzClient := services.NewDigiflazzClient(cfg)
	compensationService := services.NewCompensationService(db, redisClient)

	integrationHandler := handlers.NewIntegrationHandlerWithClient(integrationService, digiflazzClient, compensationService)

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
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	api := r.Group("/api/v1")
	api.Use(middleware.AuthMiddleware(cfg))
	{
		api.POST("/integrations/digiflazz/transaction", integrationHandler.InitiateDigiflazzTransaction)
		api.GET("/integrations/providers", integrationHandler.ListProviders)
		api.GET("/integrations/errors", integrationHandler.GetErrorCatalog)
		api.GET("/integrations/compensation/jobs", integrationHandler.GetCompensationJobs)
		api.GET("/integrations/compensation/dead-letter", integrationHandler.GetDeadLetterQueue)
	}

	webhook := r.Group("/webhook")
	{
		webhook.POST("/digiflazz", integrationHandler.HandleWebhook)
	}

	r.GET("/debug/circuit-breaker", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"digiflazz": "circuit breaker status",
		})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	go func() {
		log.Printf("Integration Service starting on port %s", cfg.ServerPort)
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