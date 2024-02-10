package database

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizface/go-ms-systemd/ms-user/user"
)

type User struct {
	dbPool *pgxpool.Pool
}

func NewUser(dbPool *pgxpool.Pool) *User {
	return &User{
		dbPool: dbPool,
	}
}

func (u *User) Insert(ctx context.Context, usr user.User) (user.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return user.User{}, err
	}

	sql := `
		insert into users (id, name, username, password) values (
			$1, $2, $3, $4
		)
	`

	_, err = conn.Exec(ctx, sql, usr.Id, usr.Name, usr.Username, usr.Password)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return user.User{}, err
	}

	return usr, nil
}

func (u *User) GetByUsername(ctx context.Context, username string) (user.User, error) {
	conn, err := u.dbPool.Acquire(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return user.User{}, err
	}

	defer conn.Release()

	var result user.User

	err = conn.QueryRow(ctx, `select id, name, username, password from users where username = $1`, username).Scan(
		&result.Id, &result.Name, &result.Username, &result.Password,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return result, user.ErrNotFound
	}

	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return result, err
	}

	return result, nil
}
