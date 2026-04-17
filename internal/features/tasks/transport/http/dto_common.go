package tasks_transport_http

import (
	"time"

	"github.com/FeeRro/golang-todoApp/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"is_complited"`
	CreatedAt    time.Time  `json:"created_at"`
	ComplitedAt  *time.Time `json:"complited_at"`
	AuthorUserID int        `json:"author_user_id"`
}

func taskDTOFromDomain(task domain.Task) TaskDTOResponse {
	return TaskDTOResponse{
		ID:           task.ID,
		Version:      task.Version,
		Title:        task.Title,
		Description:  task.Description,
		Completed:    task.Completed,
		CreatedAt:    task.CreatedAt,
		ComplitedAt:  task.CompletedAt,
		AuthorUserID: task.AuthorUserID,
	}
}

func tasksDTOFromDomains(tasks []domain.Task) []TaskDTOResponse {
	dtos := make([]TaskDTOResponse, len(tasks))

	for i, task := range tasks {
		dtos[i] = taskDTOFromDomain(task)
	}

	return dtos
}
