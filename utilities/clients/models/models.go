package utilModels

import (
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
)

type Client struct {
	UserName        string       `json:"username"`
	Password        string       `json:"password"`
	ClientId        uuid.UUID    `json:"clientId"`
	ClientSecretId  uuid.UUID    `json:"clientSecret"`
	ClientSecretStr string       `json:"clientSecretString"`
	ProjectId       uuid.UUID    `json:"projectId"`
	ProjectLabel    string       `json:"projectLabel"`
	ProjectName     string       `json:"projectName"`
	Url             string       `json:"url"`
	Header          models.JsonF `json:"header"`
	Steps           []int        `json:"steps"`

	Gender    string `json:"gender"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Phone     string `json:"phone"`
}
