package database

import (
	"CourseProject/auth_service/internal/entity"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type RefreshStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRefreshStorage() *RefreshStorage {
	password := os.Getenv("REDIS_REFRESH_PASSWORD")
	if password == "" {
		log.Println("password for redis is unable to find")
		return nil
	}
	client := redis.NewClient(&redis.Options{
		Addr:     "redis_refresh:6379",
		Password: password,
	})
	return &RefreshStorage{
		client: client,
		ctx:    context.Background(),
	}
}

func (rs *RefreshStorage) Post(refreshToken entity.RefreshToken) error {
	expiration := time.Until(refreshToken.ExpiresAt)
	if expiration <= 0 {
		return fmt.Errorf("token has already expired")
	}
	err := rs.client.Set(rs.ctx, refreshToken.UserID, refreshToken.RefreshToken, expiration).Err()
	if err != nil {
		return fmt.Errorf("unable to post refresh token: %s", err)
	}
	return nil
}

func (rs *RefreshStorage) Get(userID string) (string, error) {
	val, err := rs.client.Get(rs.ctx, userID).Result()
	if err != nil {
		return "", fmt.Errorf("unable to get value by user_id(%s): %s", userID, err)
	}
	return val, err
}

func (rs *RefreshStorage) Delete(userID string) error {
	err := rs.client.Del(rs.ctx, userID).Err()
	if err != nil {
		return fmt.Errorf("unable to delete row by user_id(%s): %s", userID, err)
	}
	return nil
}
