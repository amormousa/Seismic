package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// User represents a row in the users table.
type User struct {
	ID             string     `json:"id"`
	Username       string     `json:"username"`
	Email          string     `json:"email"`
	APIKey         string     `json:"apiKey"`
	Country        *string    `json:"country"`
	Bio            *string    `json:"bio"`
	Website        *string    `json:"website"`
	AvatarURL      *string    `json:"avatarUrl"`
	AvatarPublicID *string    `json:"-"`
	CreatedAt      time.Time  `json:"createdAt"`
	DeletedAt      *time.Time `json:"-"`
}

// FindUserByEmail looks up a user by their email address.
// Returns nil (no error) if no user was found — this is
// normal, not a failure, since a new email means we should
// create a new account instead. (We use magic links ma dude)
func FindUserByEmail(ctx context.Context, pool *pgxpool.Pool, email string) (*User, error) {
	var u User

	err := pool.QueryRow(ctx, `
		SELECT id, username, email, api_key, country, bio, website,
		       avatar_url, avatar_public_id, created_at, deleted_at
		FROM users
		WHERE email = $1 AND deleted_at IS NULL
	`, email).Scan(
		&u.ID, &u.Username, &u.Email, &u.APIKey, &u.Country, &u.Bio,
		&u.Website, &u.AvatarURL, &u.AvatarPublicID, &u.CreatedAt, &u.DeletedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil // no user found, not an error
	}

	if err != nil {
		return nil, err
	}

	return &u, nil
}

// CreateUser inserts a new user with the given email and
// username. The api_key and id are generated automatically
// by the database (see the DEFAULT values in the schema).
func CreateUser(ctx context.Context, pool *pgxpool.Pool, email, username string) (*User, error) {
	var u User

	err := pool.QueryRow(ctx, `
		INSERT INTO users (email, username)
		VALUES ($1, $2)
		RETURNING id, username, email, api_key, country, bio, website,
		          avatar_url, avatar_public_id, created_at, deleted_at
	`, email, username).Scan(
		&u.ID, &u.Username, &u.Email, &u.APIKey, &u.Country, &u.Bio,
		&u.Website, &u.AvatarURL, &u.AvatarPublicID, &u.CreatedAt, &u.DeletedAt,
	)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
