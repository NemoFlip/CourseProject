package database

import (
	customLogger "CourseProject/auth_service/pkg/log"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type VerifyCodeStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewVerifyCodeStorage(logger *customLogger.Logger) *VerifyCodeStorage {
	password := os.Getenv("REDIS_VERIFY_PASSWORD")
	if password == "" {
		logger.ErrorLogger.Error().Msg("password for redis is unable to find")
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     "redis_verify:6379",
		Password: password,
	})
	return &VerifyCodeStorage{
		client: client,
		ctx:    context.Background(),
	}
}

func (vs *VerifyCodeStorage) PostWithPrefix(prefix string, email string, value string, expTime time.Time) error {
	expiration := time.Until(expTime)
	if expiration <= 0 {
		return fmt.Errorf("token has already expired")
	}
	fullKey := fmt.Sprintf("%s:%s", prefix, email)
	err := vs.client.Set(vs.ctx, fullKey, value, 0).Err()
	if err != nil {
		return fmt.Errorf("unable to post verify token: %s", err)
	}
	return nil
}

func (vs *VerifyCodeStorage) GetWithPrefix(prefix, email string) (string, error) {
	fullKey := fmt.Sprintf("%s:%s", prefix, email)
	val, err := vs.client.Get(vs.ctx, fullKey).Result()
	if err != nil {
		return "", fmt.Errorf("unable to get value by email(%s): %s", email, err)
	}
	return val, err
}

func (vs *VerifyCodeStorage) DeleteWithPrefix(prefix, email string) error {
	fullKey := fmt.Sprintf("%s:%s", prefix, email)
	err := vs.client.Del(vs.ctx, fullKey).Err()
	if err != nil {
		return fmt.Errorf("unable to delete row by email(%s): %s", email, err)
	}
	return nil
}
