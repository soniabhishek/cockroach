package task_repo

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
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
