package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(
	ctx context.Context,
	input user.User,
	roleName role.Name,
) (user.User, error) {
	var created user.User

	err := r.db.QueryRow(
		ctx,
		`
		INSERT INTO users (
			id,
			email,
			password_hash,
			role_id
		)
		SELECT
			$1,
			$2,
			$3,
			r.id
		FROM roles r
		WHERE r.name = $4
		RETURNING
			id,
			email,
			password_hash,
			role_id,
			created_at,
			updated_at
		`,
		input.ID,
		input.Email,
		input.PasswordHash,
		roleName,
	).Scan(
		&created.ID,
		&created.Email,
		&created.PasswordHash,
		&created.RoleID,
		&created.CreatedAt,
		&created.UpdatedAt,
	)
	if err != nil {
		return user.User{}, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}

func (r *UserRepository) FindByEmail(
	ctx context.Context,
	email string,
) (user.User, error) {
	var result user.User

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			email,
			password_hash,
			role_id,
			created_at,
			updated_at
		FROM users
		WHERE LOWER(email) = LOWER($1)
		`,
		email,
	).Scan(
		&result.ID,
		&result.Email,
		&result.PasswordHash,
		&result.RoleID,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if err != nil {
		return user.User{}, fmt.Errorf("find user by email: %w", err)
	}

	return result, nil
}
