package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type Auth struct {
	ID        string       `json:"id"`
	Username  string       `json:"username"`
	Nickname  string       `json:"nickname"`
	Email     string       `json:"email"`
	Theme     ProfileTheme `json:"theme"`
	Password  string       `json:"password,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Expiry   time.Time `json:"expiry"`
}

func (s Session) isExpired() bool {
	return s.Expiry.Before(time.Now())
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateSession(ctx context.Context, username string) (string, time.Time, error) {
	// Create new session token
	sessionToken := uuid.NewString()
	expiresAt := time.Now().Add(168 * time.Hour)

	// Attempt to delete existing session
	sessionCheckQuery := `
		SELECT id
		FROM sessions
		WHERE username = $1
    `

	var existingSessionId string
	row := db.QueryRowContext(ctx, sessionCheckQuery, username)
	row.Scan(&existingSessionId)

	if existingSessionId != "" {
		err := DeleteSession(ctx, existingSessionId)

		if err != nil {
			return "", expiresAt, err
		}
	}

	sessionQuery := `
		INSERT INTO sessions (id, username, expiry)
		VALUES ($1, $2, $3)
	`
	_, err := db.ExecContext(ctx, sessionQuery, sessionToken, username, expiresAt)

	if err != nil {
		return sessionToken, expiresAt, err
	}

	return sessionToken, expiresAt, err
}

func DeleteSession(ctx context.Context, sessionToken string) error {
	deleteQuery := `
	  DELETE FROM sessions
	  WHERE id = $1
	`

	_, err := db.ExecContext(ctx, deleteQuery, sessionToken)

	if err != nil {
		return err
	}

	return nil
}

func (a *Auth) SignUp(profile Auth) (*http.Cookie, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var cookie http.Cookie

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return &cookie, 0, err
	}

	defer tx.Rollback()

	// Insert new profile
	query := `
		INSERT INTO profiles (username, nickname, email)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	stmt, err := tx.PrepareContext(ctx, query)
	var newId string

	err = stmt.QueryRowContext(
		ctx,
		profile.Username,
		profile.Nickname,
		profile.Email,
	).Scan(&newId)

	if err != nil {
		return &cookie, http.StatusInternalServerError, errors.New("serverError")
	}

	userQuery := `
		INSERT INTO users (profile_id, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
	`

	hash, err := HashPassword(profile.Password)

	if err != nil {
		return &cookie, http.StatusInternalServerError, errors.New("serverError")
	}

	_, err = tx.ExecContext(
		ctx,
		userQuery,
		newId,
		hash,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return &cookie, 0, err
	}

	sessionToken, expiresAt, err := CreateSession(ctx, profile.Username)

	cookie = http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	}

	err = tx.Commit()

	if err != nil {
		return &cookie, http.StatusInternalServerError, errors.New("serverError")
	}

	return &cookie, http.StatusOK, nil
}

func (a *Auth) SignIn(creds Credentials) (*http.Cookie, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var cookie http.Cookie

	// Search for user
	query := `
		SELECT password
		FROM profiles
		INNER JOIN users ON profiles.id = users.profile_id
		WHERE profiles.username = $1
	`

	var storedHash string
	row := db.QueryRowContext(ctx, query, creds.Username)
	err := row.Scan(&storedHash)

	if err != nil {
		return &cookie, http.StatusBadRequest, errors.New("serverError")
	}

	if !CheckPasswordHash(creds.Password, storedHash) {
		return &cookie, http.StatusUnauthorized, errors.New("invalidPassword")
	}

	sessionToken, expiresAt, err := CreateSession(ctx, creds.Username)

	cookie = http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	}

	if err != nil {
		return &cookie, http.StatusUnauthorized, errors.New("unauthorized")
	}

	return &cookie, http.StatusOK, nil
}

func (a *Auth) Refresh(sessionToken string) (*http.Cookie, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
		SELECT id, username, expiry
		FROM sessions
		WHERE id = $1
	`

	var userSession Session
	row := db.QueryRowContext(ctx, query, sessionToken)
	err := row.Scan(&userSession.ID, &userSession.Username, &userSession.Expiry)

	var cookie http.Cookie
	if err != nil || userSession.ID == "" {
		return &cookie, http.StatusUnauthorized, errors.New("unauthorized")
	}

	err = DeleteSession(ctx, sessionToken)

	if err != nil {
		// no-op
	}

	sessionToken, expiresAt, err := CreateSession(ctx, userSession.Username)
	cookie = http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: expiresAt,
	}

	if err != nil {
		return &cookie, http.StatusInternalServerError, errors.New("serverError")
	}

	return &cookie, http.StatusOK, nil
}

func (a *Auth) Logout(sessionToken string) (*http.Cookie, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	var cookie http.Cookie

	err := DeleteSession(ctx, sessionToken)

	if err != nil {
		return &cookie, http.StatusNotFound, errors.New("notFound")
	}

	cookie = http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Now(),
	}

	return &cookie, 0, nil
}
