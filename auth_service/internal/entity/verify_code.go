package entity

import "time"

type VerifyCode struct {
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}
