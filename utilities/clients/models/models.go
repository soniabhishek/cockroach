package utilModels

import (
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type Client struct {
	UserName        string          `json:"username"`
	Password        string          `json:"password"`
	ClientId        uuid.UUID       `json:"clientId"`
	ClientSecretId  uuid.UUID       `json:"clientSecret"`
	ClientSecretStr string          `json:"clientSecretString"`
	ProjectId       uuid.UUID       `json:"projectId"`
	ProjectLabel    string          `json:"projectLabel"`
	ProjectName     string          `json:"projectName"`
	Url             string          `json:"url"`
	Header          models.JsonFake `json:"header"`
	Steps           []int           `json:"steps"`

	Gender    string `json:"gender"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
}
