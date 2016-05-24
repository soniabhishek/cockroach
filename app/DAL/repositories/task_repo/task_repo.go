package task_repo

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type ITaskListRepo interface {
	GetTasks() ([]models.MicroTask, error)
	GetTasksForRole([]models.Roles) ([]models.MicroTask, error)
	GetTasksForTag([]models.Tag) ([]models.MicroTask, error)
	GetTasksForUser(userId uuid.UUID) ([]models.MicroTask, error)
}

type ITaskRepo interface {
	CreateTaskSubmission(models.MicroTask) error
	CreateTaskAttempt(userId uuid.UUID) (models.MicroTask, error)
}
