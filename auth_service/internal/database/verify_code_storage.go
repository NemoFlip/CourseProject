package database

import (
	"CourseProject/auth_service/internal/entity"
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

func (vs *VerifyCodeStorage) Post(verifyCode entity.VerifyCode) error {
	expiration := time.Until(verifyCode.ExpiresAt)
	if expiration <= 0 {
		return fmt.Errorf("token has already expired")
	}
	err := vs.client.Set(vs.ctx, verifyCode.Email, verifyCode.Code, 0).Err()
	if err != nil {
		return fmt.Errorf("unable to post verify token: %s", err)
	}
	return nil
}

func (vs *VerifyCodeStorage) Get(email string) (string, error) {
	val, err := vs.client.Get(vs.ctx, email).Result()
	if err != nil {
		return "", fmt.Errorf("unable to get value by email(%s): %s", email, err)
	}
	return val, err
}

func (vs *VerifyCodeStorage) Delete(email string) error {
	err := vs.client.Del(vs.ctx, email).Err()
	if err != nil {
		return fmt.Errorf("unable to delete row by email(%s): %s", email, err)
	}
	return nil
}
