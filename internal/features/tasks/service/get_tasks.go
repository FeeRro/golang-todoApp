package tasks_service

import (
	"context"
	"fmt"

	"github.com/FeeRro/golang-todoApp/internal/core/domain"
	core_error "github.com/FeeRro/golang-todoApp/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"limit must be non-negative: %w",
			core_error.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"offset must be non-negative: %w",
			core_error.ErrInvalidArgument,
		)
	}

	Tasks, err := s.tasksRepository.GetTasks(
		ctx,
		userID,
		limit,
		offset,
	)
	if err != nil {
		return nil, fmt.Errorf("get users from repository: %w", err)
	}

	return Tasks, nil
}
