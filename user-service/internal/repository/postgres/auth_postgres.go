package postgres

import (
	"time"
	"user-service/internal/models"

	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

func (r *Repository) CreateUser(ctx context.Context, user models.User) (int, error) {
	var id int

	insert := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("users").
		Columns("username", "password", "role", "date_create").
		Suffix("RETURNING id").
		Values(user.Username, user.Password, user.Role, time.Now())

	sql, args, err := insert.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "insert.ToSql()")
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "CreateUser()")
	}

	return id, nil
}

func (r *Repository) GetUser(ctx context.Context, username, password string) (models.User, error) {
	var user models.User

	sel := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "username", "password", "role", "date_create").
		From("users").
		Where(squirrel.Eq{"username": username, "password": password})

	sql, args, err := sel.ToSql()
	if err != nil {
		return models.User{}, errors.Wrap(err, "sel.ToSql()")
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Username, &user.Password, &user.Role, &user.DateCreate)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, models.ErrNotFound
		}
		return models.User{}, errors.Wrap(err, "GetUser()")
	}

	return user, nil
}
