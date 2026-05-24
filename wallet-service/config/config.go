package config

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	ServerPort      string
	GRPCPort        string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	RedisHost       string
	RedisPort       string
	JWTPublicKey    string
	PublicKey       *rsa.PublicKey
	GinMode         string
	JaegerEndpoint string
}

func Load() *Config {
	cfg := &Config{
		ServerPort:      getEnv("SERVER_PORT", "8080"),
		GRPCPort:        getEnv("GRPC_PORT", "50083"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "ppob"),
		RedisHost:       getEnv("REDIS_HOST", "localhost"),
		RedisPort:       getEnv("REDIS_PORT", "6379"),
		JWTPublicKey:    getEnv("JWT_PUBLIC_KEY", ""),
		GinMode:         getEnv("GIN_MODE", "release"),
		JaegerEndpoint: getEnv("JAEGER_ENDPOINT", ""),
	}

	if cfg.JWTPublicKey == "" {
		if b, err := os.ReadFile("certs/public.pem"); err == nil {
			cfg.JWTPublicKey = string(b)
		}
	}

	if cfg.JWTPublicKey != "" {
		pub, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cfg.JWTPublicKey))
		if err != nil {
			log.Printf("Warning: failed to parse JWT public key: %v\n", err)
		} else {
			cfg.PublicKey = pub
		}
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func InitDB(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return db, nil
}

func InitRedis(cfg *Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: "",
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return client
}