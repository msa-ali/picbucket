package models

import (
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/msa-ali/picbucket/utils"
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
	// verify if email is valid and get userid
	email = strings.ToLower(email)
	var userID int
	row := srv.DB.QueryRow(`
		SELECT id FROM users WHERE email = $1;
	`, email)
	err := row.Scan(&userID)
	if err != nil {
		// @TODO: return user doesn't exist error
		return nil, fmt.Errorf("create: %w", err)
	}
	// build password reset
	bytesPerToken := srv.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	token, err := utils.UUID(bytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	duration := srv.Duration
	if duration == 0 {
		duration = DefaultResetDuration
	}
	pwReset := PasswordReset{
		UserID:    userID,
		Token:     token,
		TokenHash: srv.hash(token),
		ExpiresAt: time.Now().Add(duration),
	}
	// insert into db
	row = srv.DB.QueryRow(`
			INSERT INTO password_resets (user_id, token_hash, expires_at)
			VALUES ($1, $2, $3)
			ON CONFLICT (user_id) DO
				UPDATE
				SET token_hash = $2, expires_at = $3
			RETURNING id;`, pwReset.UserID, pwReset.TokenHash, pwReset.ExpiresAt,
	)
	err = row.Scan(&pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("create: %w", err)
	}
	return &pwReset, nil
}

func (srv *PasswordResetService) Consume(token string) (*User, error) {
	tokenHash := srv.hash(token)
	var user User
	var pwReset PasswordReset
	row := srv.DB.QueryRow(`
		SELECT password_resets.id,
			   password_resets.expires_at,
			   users.id,
			   users.email,
			   users.password_hash
		FROM password_resets
		JOIN users ON users.id = password_resets.user_id
		WHERE password_resets.token_hash = $1
	`, tokenHash)
	err := row.Scan(
		&pwReset.ID,
		&pwReset.ExpiresAt,
		&user.ID,
		&user.Email,
		&user.PasswordHash,
	)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	if time.Now().After(pwReset.ExpiresAt) {
		return nil, fmt.Errorf("token expired: %v", token)
	}
	err = srv.delete(pwReset.ID)
	if err != nil {
		return nil, fmt.Errorf("consume: %w", err)
	}
	return &user, nil
}

func (srv *PasswordResetService) hash(token string) string {
	tokenHash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(tokenHash[:])
}

func (srv *PasswordResetService) delete(id int) error {
	_, err := srv.DB.Exec(`
	DELETE FROM password_resets
	WHERE id = $1;
	`, id)
	if err != nil {
		return fmt.Errorf("delete : %w", err)
	}
	return nil
}
