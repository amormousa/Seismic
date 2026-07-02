package models

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MagicLink struct {
	ID        string
	Email     string
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}

// CreateMagicLink inserts a new magic link row for the given
// email, set to expire 15 minutes from now. The token itself
// is generated automatically by the database.
func CreateMagicLink(ctx context.Context, pool *pgxpool.Pool, email string) (*MagicLink, error) {
	var m MagicLink
	expiresAt := time.Now().Add(15 * time.Minute)

	err := pool.QueryRow(ctx, `
		INSERT INTO magic_links (email, expires_at)
		VALUES ($1, $2)
		RETURNING id, email, token, expires_at, used, created_at
	`, email, expiresAt).Scan(
		&m.ID, &m.Email, &m.Token, &m.ExpiresAt, &m.Used, &m.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

// FindMagicLinkByToken looks up a magic link by its token.
// Returns nil (no error) if no matching token was found.
func FindMagicLinkByToken(ctx context.Context, pool *pgxpool.Pool, token string) (*MagicLink, error) {
	var m MagicLink

	err := pool.QueryRow(ctx, `
		SELECT id, email, token, expires_at, used, created_at
		FROM magic_links
		WHERE token = $1
	`, token).Scan(
		&m.ID, &m.Email, &m.Token, &m.ExpiresAt, &m.Used, &m.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func MarkMagicLinkUsed(ctx context.Context, pool *pgxpool.Pool, id string) error {
	_, err := pool.Exec(ctx, `
		UPDATE magic_links SET used = true WHERE id = $1
	`, id)
	return err
}

func (m *MagicLink) IsExpired() bool {
	return time.Now().After(m.ExpiresAt)
}
