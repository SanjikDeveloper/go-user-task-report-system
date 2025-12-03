package postgres

import (
	"context"
	"fmt"
	"tasks-service/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type TaskPostgres struct {
	db *pgxpool.Pool
}

func NewTaskPostgres(db *pgxpool.Pool) *TaskPostgres {
	return &TaskPostgres{db: db}
}

func (r *TaskPostgres) CreateTask(ctx context.Context, task models.Task) (int, error) {
	// Сначала проверяем существование пользователя
	var userExists bool
	err := r.db.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)",
		task.UserID,
	).Scan(&userExists)

	if err != nil {
		return 0, errors.Wrap(err, fmt.Sprintf("failed to check user existence for user_id=%d: %v", task.UserID, err))
	}

	if !userExists {
		// Дополнительная проверка - может быть проблема с подключением к другой базе
		var totalUsers int
		var dbName string
		_ = r.db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&totalUsers)
		_ = r.db.QueryRow(ctx, "SELECT current_database()").Scan(&dbName)
		return 0, errors.New(fmt.Sprintf("user with id %d does not exist in database '%s' (found %d total users). Please make sure user-service and tasks-service use the same database (tasks_service_db), or sign in again", task.UserID, dbName, totalUsers))
	}

	var id int

	insert := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Insert("tasks").
		Columns("user_id", "title", "description", "status", "created_at", "updated_at").
		Values(task.UserID, task.Title, task.Description, task.Status, task.CreatedAt, task.UpdatedAt).
		Suffix("RETURNING id")

	sql, args, err := insert.ToSql()
	if err != nil {
		return 0, errors.Wrap(err, "insert.ToSql()")
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		return 0, errors.Wrap(err, "failed to create task")
	}

	return id, nil
}

func (r *TaskPostgres) GetTaskById(ctx context.Context, id, userID int) (models.Task, error) {
	var task models.Task

	sel := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("id", "user_id", "title", "description", "status", "created_at", "updated_at").
		From("tasks").
		Where(squirrel.Eq{"id": id, "user_id": userID})

	sql, args, err := sel.ToSql()
	if err != nil {
		return models.Task{}, errors.Wrap(err, "insert.ToSql()")
	}

	err = r.db.QueryRow(ctx, sql, args...).Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return models.Task{}, errors.New("task not found or unauthorized access")
	}

	return task, nil
}

func (r *TaskPostgres) UpdateTask(ctx context.Context, task models.Task) error {
	update := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Update("tasks").
		Set("updated_at", task.UpdatedAt).
		Where(squirrel.And{
			squirrel.Eq{"id": task.ID},
			squirrel.Eq{"user_id": task.UserID},
		})

	if task.Title != "" {
		update = update.Set("title", task.Title)
	}
	if task.Description != "" {
		update = update.Set("description", task.Description)
	}
	if task.Status != "" {
		update = update.Set("status", task.Status)
	}

	sql, args, err := update.ToSql()
	if err != nil {
		return errors.Wrap(err, "updateToSql()")
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskPostgres) DeleteTask(ctx context.Context, id, userID int) error {

	del := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Delete("tasks").
		Where(squirrel.Eq{"id": id, "user_id": userID})

	sql, args, err := del.ToSql()
	if err != nil {
		return errors.Wrap(err, "del.ToSql()")
	}

	res, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return err
	}

	rowsAffected := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no task found")
	}

	return nil
}
