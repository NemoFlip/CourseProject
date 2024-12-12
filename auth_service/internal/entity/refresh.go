package entity

import "time"

type RefreshToken struct {
	UserID       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}
