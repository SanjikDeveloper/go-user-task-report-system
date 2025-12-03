package application

import (
	"context"
	"tasks-service/internal/models"
	"tasks-service/internal/repository/postgres"
)

type TaskService struct {
	repo *postgres.TaskPostgres
}

func NewTaskService(repo *postgres.TaskPostgres) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(ctx context.Context, task models.Task) (int, error) {
	return s.repo.CreateTask(ctx, task)
}

func (s *TaskService) GetTaskById(ctx context.Context, id, userID int) (models.Task, error) {
	return s.repo.GetTaskById(ctx, id, userID)
}

func (s *TaskService) DeleteTask(ctx context.Context, id, userID int) error {
	return s.repo.DeleteTask(ctx, id, userID)
}

func (s *TaskService) UpdateTask(ctx context.Context, task models.Task) error {
	return s.repo.UpdateTask(ctx, task)
}
