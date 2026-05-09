package services

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenBlacklistService struct {
	redis *redis.Client
}

func NewTokenBlacklistService(redis *redis.Client) *TokenBlacklistService {
	return &TokenBlacklistService{redis: redis}
}

func (s *TokenBlacklistService) AddToBlacklist(ctx context.Context, tokenID string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", tokenID)
	return s.redis.Set(ctx, key, "1", expiration).Err()
}

func (s *TokenBlacklistService) IsBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", tokenID)
	result, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func (s *TokenBlacklistService) RevokeAllUserTokens(ctx context.Context, userID uint) error {
	pattern := fmt.Sprintf("refresh_token:%d", userID)
	
	keys, err := s.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	for _, key := range keys {
		s.redis.Del(ctx, key)
	}

	return nil
}

func (s *TokenBlacklistService) AddAccessTokenToBlacklist(ctx context.Context, tokenID string, expiration time.Duration) error {
	key := fmt.Sprintf("access_blacklist:%s", tokenID)
	return s.redis.Set(ctx, key, "1", expiration).Err()
}

func (s *TokenBlacklistService) IsAccessTokenBlacklisted(ctx context.Context, tokenID string) (bool, error) {
	key := fmt.Sprintf("access_blacklist:%s", tokenID)
	result, err := s.redis.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}