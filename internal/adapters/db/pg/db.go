package pg

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"auth/internal/models"
	"auth/internal/adapters/store/userstore"
)

type pgdb struct {
	pool *pgxpool.Pool
}

func NewPgDb(ctx context.Context, conn string) (userstore.User, error) {
	pool, err := pgxpool.Connect(ctx, conn)
	if err != nil {
		return nil, err
	}

	return &pgdb{
		pool: pool,
	}, nil
}

func (db *pgdb) Get(ctx context.Context, login string) (*models.User, error) {
	const query = `
	SELECT * FROM users
	WHERE login = $1
`
	out := &models.User{}

	err := db.pool.QueryRow(ctx, query, login).Scan(
		&out.Login,
		&out.Password,
		&out.Role,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return out, fmt.Errorf("can not get user: %w", err)
	}
	return out, nil
}
/*
func (db *pgdb) Add(ctx context.Context, login string, password string, role string) (uint64, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback(ctx)

	const query = `
	INSERT INTO tasks (login, password, role)
	VALUES ($1, $2, $3)
	RETURNING id
`
	var id uint64
	err = tx.QueryRow(ctx, query, login, password, role).Scan(&id)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)

	return id, err
}
*/