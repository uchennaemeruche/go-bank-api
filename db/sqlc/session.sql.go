// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.15.0
// source: session.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions(
    id,
    username,
    refresh_token,
    user_agent,
    client_ip,
    is_blocked,
    expires_at
) VALUES(
    $1, $2, $3, $4, $5, $6, $7
) RETURNING id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
`

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	RefreshToken string    `json:"refresh_token"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.ID,
		arg.Username,
		arg.RefreshToken,
		arg.UserAgent,
		arg.ClientIp,
		arg.IsBlocked,
		arg.ExpiresAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const expireSession = `-- name: ExpireSession :exec
UPDATE sessions
SET expires_at = $2
WHERE id = $1
`

type ExpireSessionParams struct {
	ID        uuid.UUID `json:"id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) ExpireSession(ctx context.Context, arg ExpireSessionParams) error {
	_, err := q.db.ExecContext(ctx, expireSession, arg.ID, arg.ExpiresAt)
	return err
}

const getSession = `-- name: GetSession :one
SELECT id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at FROM sessions
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSession(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSession, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const toggleBlockSession = `-- name: ToggleBlockSession :exec
UPDATE sessions
SET is_blocked = $2
WHERE id = $1
`

type ToggleBlockSessionParams struct {
	ID        uuid.UUID `json:"id"`
	IsBlocked bool      `json:"is_blocked"`
}

func (q *Queries) ToggleBlockSession(ctx context.Context, arg ToggleBlockSessionParams) error {
	_, err := q.db.ExecContext(ctx, toggleBlockSession, arg.ID, arg.IsBlocked)
	return err
}

const updateSession = `-- name: UpdateSession :one
UPDATE sessions
SET
    is_blocked = CASE WHEN $1::boolean
    THEN $2::bool ELSE is_blocked END,

    expires_at = CASE WHEN $3::boolean
    THEN $4::timestamp ELSE expires_at END
WHERE 
    id = $5
RETURNING id, username, refresh_token, user_agent, client_ip, is_blocked, expires_at, created_at
`

type UpdateSessionParams struct {
	IsBlockedToUpdate bool      `json:"is_blocked_to_update"`
	IsBlocked         bool      `json:"is_blocked"`
	ExpiresAtToUpdate bool      `json:"expires_at_to_update"`
	ExpiresAt         time.Time `json:"expires_at"`
	ID                uuid.UUID `json:"id"`
}

func (q *Queries) UpdateSession(ctx context.Context, arg UpdateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, updateSession,
		arg.IsBlockedToUpdate,
		arg.IsBlocked,
		arg.ExpiresAtToUpdate,
		arg.ExpiresAt,
		arg.ID,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.RefreshToken,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}
