package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/labsession"
)

type LabSessionRepository struct {
	db *pgxpool.Pool
}

func NewLabSessionRepository(
	db *pgxpool.Pool,
) *LabSessionRepository {
	return &LabSessionRepository{
		db: db,
	}
}

func (r *LabSessionRepository) Create(
	ctx context.Context,
	session labsession.Session,
) (labsession.Session, error) {
	var created labsession.Session

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO lab_sessions (
			id,
			lab_id,
			user_id,
			status,
			expires_at
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING
			id,
			lab_id,
			user_id,
			status,
			namespace,
			pod_name,
			created_at,
			started_at,
			expires_at,
			finished_at,
			failure_reason
		`,
		session.ID,
		session.LabID,
		session.UserID,
		session.Status,
		session.ExpiresAt,
	).Scan(
		&created.ID,
		&created.LabID,
		&created.UserID,
		&created.Status,
		&created.Namespace,
		&created.PodName,
		&created.CreatedAt,
		&created.StartedAt,
		&created.ExpiresAt,
		&created.FinishedAt,
		&created.FailureReason,
	)
	if err != nil {
		return labsession.Session{},
			fmt.Errorf("create lab session: %w", err)
	}

	return created, nil
}

func (r *LabSessionRepository) FindByID(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) (labsession.Session, error) {
	var session labsession.Session

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			lab_id,
			user_id,
			status,
			namespace,
			pod_name,
			created_at,
			started_at,
			expires_at,
			finished_at,
			failure_reason
		FROM lab_sessions
		WHERE id = $1
		  AND user_id = $2
		`,
		id,
		userID,
	).Scan(
		&session.ID,
		&session.LabID,
		&session.UserID,
		&session.Status,
		&session.Namespace,
		&session.PodName,
		&session.CreatedAt,
		&session.StartedAt,
		&session.ExpiresAt,
		&session.FinishedAt,
		&session.FailureReason,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return labsession.Session{}, labsession.ErrNotFound
	}

	if err != nil {
		return labsession.Session{},
			fmt.Errorf("find lab session: %w", err)
	}

	return session, nil
}

func (r *LabSessionRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (labsession.Session, error) {
	var session labsession.Session

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			lab_id,
			user_id,
			status,
			namespace,
			pod_name,
			created_at,
			started_at,
			expires_at,
			finished_at,
			failure_reason
		FROM lab_sessions
		WHERE id = $1
		`,
		id,
	).Scan(
		&session.ID,
		&session.LabID,
		&session.UserID,
		&session.Status,
		&session.Namespace,
		&session.PodName,
		&session.CreatedAt,
		&session.StartedAt,
		&session.ExpiresAt,
		&session.FinishedAt,
		&session.FailureReason,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return labsession.Session{}, labsession.ErrNotFound
	}

	if err != nil {
		return labsession.Session{},
			fmt.Errorf("find lab session: %w", err)
	}

	return session, nil
}

func (r *LabSessionRepository) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]labsession.Session, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			lab_id,
			user_id,
			status,
			namespace,
			pod_name,
			created_at,
			started_at,
			expires_at,
			finished_at,
			failure_reason
		FROM lab_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
		`,
		userID,
	)
	if err != nil {
		return nil,
			fmt.Errorf("list lab sessions: %w", err)
	}
	defer rows.Close()

	sessions := make([]labsession.Session, 0)

	for rows.Next() {
		var session labsession.Session

		if err := rows.Scan(
			&session.ID,
			&session.LabID,
			&session.UserID,
			&session.Status,
			&session.Namespace,
			&session.PodName,
			&session.CreatedAt,
			&session.StartedAt,
			&session.ExpiresAt,
			&session.FinishedAt,
			&session.FailureReason,
		); err != nil {
			return nil,
				fmt.Errorf("scan lab session: %w", err)
		}

		sessions = append(sessions, session)
	}

	if err := rows.Err(); err != nil {
		return nil,
			fmt.Errorf("iterate lab sessions: %w", err)
	}

	return sessions, nil
}

func (r *LabSessionRepository) Stop(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
) error {
	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE lab_sessions
		SET
			status = 'stopped',
			finished_at = COALESCE(finished_at, NOW())
		WHERE id = $1
		  AND user_id = $2
		`,
		id,
		userID,
	)
	if err != nil {
		return fmt.Errorf("stop lab session: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return labsession.ErrNotFound
	}

	return nil
}

func (r *LabSessionRepository) MarkProvisioning(
	ctx context.Context,
	id uuid.UUID,
) error {
	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE lab_sessions
		SET status = $2
		WHERE id = $1
		  AND status = $3
		`,
		id,
		labsession.StatusProvisioning,
		labsession.StatusPending,
	)
	if err != nil {
		return fmt.Errorf(
			"mark lab session provisioning: %w",
			err,
		)
	}

	if tag.RowsAffected() == 0 {
		return labsession.ErrNotFound
	}

	return nil
}

func (r *LabSessionRepository) MarkRunning(
	ctx context.Context,
	id uuid.UUID,
) error {
	tag, err := r.db.Exec(
		ctx,
		`
		UPDATE lab_sessions
		SET status = $2
		WHERE id = $1
		  AND status = $3
		`,
		id,
		labsession.StatusRunning,
		labsession.StatusProvisioning,
	)
	if err != nil {
		return fmt.Errorf(
			"mark lab session running: %w",
			err,
		)
	}

	if tag.RowsAffected() == 0 {
		return labsession.ErrNotFound
	}

	return nil
}
