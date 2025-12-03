package application

import (
	"context"
	"fmt"
	"report-service/internal/models"
	"time"
)

func (s *Service) ReportByRange(ctx context.Context, userID int64, username string, start, end time.Time) (models.ReportVM, error) {
	tasks, err := s.repo.ListByRange(ctx, userID, start, end)
	if err != nil {
		return models.ReportVM{}, err
	}

	return s.PrepareReportData(username, start, end, tasks), nil
}

func (s *Service) PrepareReportData(username string, start, end time.Time, tasks []models.Task) models.ReportVM {
	vm := models.ReportVM{
		Period: fmt.Sprintf("%s..%s", start.Format("2006-01-02"), end.Format("2006-01-02")),
	}

	for _, task := range tasks {
		vm.Tasks = append(vm.Tasks, models.Task{
			ID:          task.ID,
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
		})

		if task.Status == "completed" {
			vm.DoneCount++
		} else {
			vm.OtherCount++
		}
	}

	return vm
}
