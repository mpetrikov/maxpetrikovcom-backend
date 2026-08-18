package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/lab"
)

type LabRepository struct {
	db *pgxpool.Pool
}

func NewLabRepository(
	db *pgxpool.Pool,
) *LabRepository {
	return &LabRepository{
		db: db,
	}
}

func (r *LabRepository) Create(
	ctx context.Context,
	input lab.Lab,
) (lab.Lab, error) {
	var created lab.Lab

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO labs (
			id,
			slug,
			title,
			description,
			difficulty,
			timeout_minutes,
			image,
			cpu_limit,
			memory_limit,
			is_published
		)
		VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10
		)
		RETURNING
			id,
			slug,
			title,
			description,
			difficulty,
			timeout_minutes,
			image,
			cpu_limit,
			memory_limit,
			is_published,
			created_at,
			updated_at
		`,
		input.ID,
		input.Slug,
		input.Title,
		input.Description,
		input.Difficulty,
		input.TimeoutMinutes,
		input.Image,
		input.CPULimit,
		input.MemoryLimit,
		input.IsPublished,
	).Scan(
		&created.ID,
		&created.Slug,
		&created.Title,
		&created.Description,
		&created.Difficulty,
		&created.TimeoutMinutes,
		&created.Image,
		&created.CPULimit,
		&created.MemoryLimit,
		&created.IsPublished,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) &&
			pgErr.Code == pgUniqueViolation &&
			pgErr.ConstraintName == labsSlugUniqueConstraint {

			return lab.Lab{}, lab.ErrSlugAlreadyExists
		}

		return lab.Lab{}, fmt.Errorf("create lab: %w", err)
	}

	return created, nil
}

func (r *LabRepository) FindBySlug(
	ctx context.Context,
	slug string,
) (lab.Lab, error) {
	var result lab.Lab

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			slug,
			title,
			description,
			difficulty,
			timeout_minutes,
			image,
			cpu_limit,
			memory_limit,
			is_published,
			created_at,
			updated_at
		FROM labs
		WHERE slug = $1
		`,
		slug,
	).Scan(
		&result.ID,
		&result.Slug,
		&result.Title,
		&result.Description,
		&result.Difficulty,
		&result.TimeoutMinutes,
		&result.Image,
		&result.CPULimit,
		&result.MemoryLimit,
		&result.IsPublished,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return lab.Lab{}, lab.ErrNotFound
	}

	if err != nil {
		return lab.Lab{}, fmt.Errorf("find lab by slug: %w", err)
	}

	return result, nil
}

func (r *LabRepository) FindByID(
	ctx context.Context,
	labId uuid.UUID,
) (lab.Lab, error) {
	var result lab.Lab

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			slug,
			title,
			description,
			difficulty,
			timeout_minutes,
			image,
			cpu_limit,
			memory_limit,
			is_published,
			created_at,
			updated_at
		FROM labs
		WHERE id = $1
		`,
		labId,
	).Scan(
		&result.ID,
		&result.Slug,
		&result.Title,
		&result.Description,
		&result.Difficulty,
		&result.TimeoutMinutes,
		&result.Image,
		&result.CPULimit,
		&result.MemoryLimit,
		&result.IsPublished,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return lab.Lab{}, lab.ErrNotFound
	}

	if err != nil {
		return lab.Lab{}, fmt.Errorf("find lab by slug: %w", err)
	}

	return result, nil
}

func (r *LabRepository) ListPublished(
	ctx context.Context,
) ([]lab.Lab, error) {
	rows, err := r.db.Query(
		ctx,
		`
		SELECT
			id,
			slug,
			title,
			description,
			difficulty,
			timeout_minutes,
			image,
			cpu_limit,
			memory_limit,
			is_published,
			created_at,
			updated_at
		FROM labs
		WHERE is_published = TRUE
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("list published labs: %w", err)
	}
	defer rows.Close()

	var labs []lab.Lab

	for rows.Next() {
		var item lab.Lab

		if err := rows.Scan(
			&item.ID,
			&item.Slug,
			&item.Title,
			&item.Description,
			&item.Difficulty,
			&item.TimeoutMinutes,
			&item.Image,
			&item.CPULimit,
			&item.MemoryLimit,
			&item.IsPublished,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan lab: %w", err)
		}

		labs = append(labs, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate labs: %w", err)
	}

	return labs, nil
}
