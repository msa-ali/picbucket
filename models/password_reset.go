package models

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	DefaultResetDuration = 1 * time.Hour
)

type PasswordReset struct {
	ID     int
	UserID int
	// Token is only set when a password reset is being created
	Token     string
	TokenHash string
	ExpiresAt time.Time
}

type PasswordResetService struct {
	DB            *sql.DB
	BytesPerToken int
	// amount of time that a password reset is valid for.
	// defaults to Default Reset duration
	Duration time.Duration
}

func (srv *PasswordResetService) Create(email string) (*PasswordReset, error) {
	return nil, fmt.Errorf("todo")
}

func (srv *PasswordResetService) Consume(token string) (*User, error) {
	return nil, fmt.Errorf("todo")
}
