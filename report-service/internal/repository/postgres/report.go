package postgres

import (
	"context"
	"report-service/internal/models"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (r *Repository) ListByRange(ctx context.Context, userId int64, start, end time.Time) ([]models.Task, error) {
	sel := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "user_id", "title", "description", "status", "created_at", "updated_at").
		From("tasks").
		Where(squirrel.Eq{"user_id": userId}).
		Where(squirrel.And{squirrel.GtOrEq{"created_at": start}, squirrel.LtOrEq{"created_at": end}}).
		OrderBy("created_at ASC")

	sql, args, err := sel.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ListByRange ToSql")
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "ListByRange Query")
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "ListByRange Rows err")
	}

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err = rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "ListByRange Scan")
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *Repository) ListByDate(ctx context.Context, userId int64, date time.Time) ([]models.Task, error) {
	sel := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "user_id", "title", "description", "status", "created_at", "updated_at").
		From("tasks").
		Where(squirrel.Eq{"user_id": userId, "created_at": date}).
		OrderBy("created_at ASC")

	sql, args, err := sel.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "ListByDate ToSql")
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "ListByDate Query")
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err = rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Description, &t.Status, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, errors.Wrap(err, "ListByDate Scan")
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}
