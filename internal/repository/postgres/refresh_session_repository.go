package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/refreshsession"
)

type RefreshSessionRepository struct {
	db *pgxpool.Pool
}

func NewRefreshSessionRepository(
	db *pgxpool.Pool,
) *RefreshSessionRepository {
	return &RefreshSessionRepository{
		db: db,
	}
}

func (r *RefreshSessionRepository) Create(
	ctx context.Context,
	session refreshsession.Session,
) error {
	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO refresh_sessions (
			id,
			user_id,
			token_hash,
			expires_at
		)
		VALUES ($1, $2, $3, $4)
		`,
		session.ID,
		session.UserID,
		session.TokenHash,
		session.ExpiresAt,
	)
	if err != nil {
		return fmt.Errorf("create refresh session: %w", err)
	}

	return nil
}

func (r *RefreshSessionRepository) FindByTokenHash(
	ctx context.Context,
	tokenHash string,
) (refreshsession.Session, error) {
	var session refreshsession.Session

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			token_hash,
			expires_at,
			created_at,
			revoked_at
		FROM refresh_sessions
		WHERE token_hash = $1
		`,
		tokenHash,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.TokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.RevokedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return refreshsession.Session{}, refreshsession.ErrNotFound
	}

	if err != nil {
		return refreshsession.Session{},
			fmt.Errorf("find refresh session: %w", err)
	}

	return session, nil
}

func (r *RefreshSessionRepository) Revoke(
	ctx context.Context,
	tokenHash string,
) error {
	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE refresh_sessions
		SET revoked_at = NOW()
		WHERE token_hash = $1
		  AND revoked_at IS NULL
		`,
		tokenHash,
	)
	if err != nil {
		return fmt.Errorf("revoke refresh session: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return refreshsession.ErrNotFound
	}

	return nil
}
