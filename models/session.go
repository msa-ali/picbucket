package models

import "database/sql"

type Session struct {
	ID     int
	UserID int
	// Token is only set when creating a new session.
	// When looing up a session, this will be left empty,
	// as we only store the hash of a session token in
	// our database and we can't reverse it into a raw token.
	Token     string
	TokenHash string
}

type SessionService struct {
	DB *sql.DB
}

func (ss *SessionService) Create(userID int) (*Session, error) {
	return nil, nil
}

func (ss *SessionService) User(token string) (*User, error) {
	return nil, nil
}

func (ss *SessionService) Delete(userID int) error {
	return nil
}
