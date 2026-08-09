package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

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

func (r *UserIdentityRepository) Create(
	ctx context.Context,
	input useridentity.Identity,
) (useridentity.Identity, error) {
	var created useridentity.Identity

	err := r.db.QueryRow(
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
		RETURNING
			id,
			user_id,
			provider,
			provider_user_id,
			email,
			created_at
		`,
		input.ID,
		input.UserID,
		input.Provider,
		input.ProviderUserID,
		input.Email,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.Provider,
		&created.ProviderUserID,
		&created.Email,
		&created.CreatedAt,
	)
	if err != nil {
		return useridentity.Identity{},
			fmt.Errorf("create user identity: %w", err)
	}

	return created, nil
}

func (r *UserIdentityRepository) FindByProvider(
	ctx context.Context,
	provider useridentity.Provider,
	providerUserID string,
) (useridentity.Identity, error) {
	var result useridentity.Identity

	err := r.db.QueryRow(
		ctx,
		`
		SELECT
			id,
			user_id,
			provider,
			provider_user_id,
			email,
			created_at
		FROM user_identities
		WHERE provider = $1
		  AND provider_user_id = $2
		`,
		provider,
		providerUserID,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.Provider,
		&result.ProviderUserID,
		&result.Email,
		&result.CreatedAt,
	)
	if err != nil {
		return useridentity.Identity{},
			fmt.Errorf("find user identity: %w", err)
	}

	return result, nil
}
