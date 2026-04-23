package tasks_transport_http

import (
	"time"

	"github.com/FeeRro/golang-todoApp/internal/core/domain"
)

type TaskDTOResponse struct {
	ID           int        `json:"id"              example:"10"`
	Version      int        `json:"version"         example:"3"`
	Title        string     `json:"title"           example:"Сделать домашнее задание"`
	Description  *string    `json:"description"     example:"Сделать матан до 18 мая"`
	Completed    bool       `json:"is_complited"    example:"false"`
	CreatedAt    time.Time  `json:"created_at"      example:"2024-05-15T14:00:00Z"`
	ComplitedAt  *time.Time `json:"complited_at"    example:"null"`
	AuthorUserID int        `json:"author_user_id"  example:"5"`
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
