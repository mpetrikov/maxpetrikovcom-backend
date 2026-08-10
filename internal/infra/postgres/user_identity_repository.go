package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/role"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/user"
	"github.com/maxpetrikov/maxpetrikovcom-backend/internal/domain/useridentity"
)

type UserIdentityRepository struct {
	db *pgxpool.Pool
}

func NewUserIdentityRepository(
	db *pgxpool.Pool,
) *UserIdentityRepository {
	return &UserIdentityRepository{
		db: db,
	}
}

func (r *UserIdentityRepository) FindUserByProvider(
	ctx context.Context,
	provider useridentity.Provider,
	providerUserID string,
) (user.User, error) {
	var result user.User

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			u.id,
			u.email,
			u.password_hash,
			u.role_id,
			u.created_at,
			u.updated_at
		FROM user_identities ui
		JOIN users u ON u.id = ui.user_id
		WHERE ui.provider = $1
		  AND ui.provider_user_id = $2
		`,
		provider,
		providerUserID,
	).Scan(
		&result.ID,
		&result.Email,
		&result.PasswordHash,
		&result.RoleID,
		&result.CreatedAt,
		&result.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return user.User{}, useridentity.ErrNotFound
	}

	if err != nil {
		return user.User{}, fmt.Errorf(
			"find user by identity: %w",
			err,
		)
	}

	return result, nil
}

func (r *UserIdentityRepository) CreateUserWithIdentity(
	ctx context.Context,
	inputUser user.User,
	inputIdentity useridentity.Identity,
	roleName role.Name,
) (user.User, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return user.User{}, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	var created user.User

	err = tx.QueryRow(
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
		inputUser.ID,
		inputUser.Email,
		inputUser.PasswordHash,
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
		return user.User{}, fmt.Errorf("create oauth user: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		`
		INSERT INTO user_identities (
			id,
			user_id,
			provider,
			provider_user_id,
			email
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		inputIdentity.ID,
		created.ID,
		inputIdentity.Provider,
		inputIdentity.ProviderUserID,
		inputIdentity.Email,
	)
	if err != nil {
		return user.User{}, fmt.Errorf(
			"create user identity: %w",
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return user.User{}, fmt.Errorf("commit transaction: %w", err)
	}

	return created, nil
}
