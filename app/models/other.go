package models

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/crowdflux/angel/app/models/flu_upload_status"
	"github.com/crowdflux/angel/app/models/uuid"
)

type ImageContainer struct {
	Id  string
	Url string
}

type ImageDictionaryNew struct {
	ID             uuid.UUID
	Label          string
	OriginalUrl    string
	IsCached       bool
	CdnUrl         string
	ResizedCdnUrl  string
	CreatedAt      pq.NullTime
	BatchProcessId uuid.UUID
}

type ImageDictionary1 struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	RealUrl   string         `db:"real_url" json:"real_url"`
	CloudUrl  string         `db:"cloud_url" json:"cloud_url"`
	Extra     sql.NullString `db:"extra" json:"extra"`
	CreatedAt pq.NullTime    `db:"created_at" json:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at" json:"updated_at"`
}

type BatchProcess struct {
	ID          uuid.UUID
	Name        string
	Done        bool
	Completion  float64
	Type        int8
	CreatedAt   pq.NullTime
	MacroTaskId uuid.UUID
}

type RouteWithLogicGate struct {
	Route
	LogicGate LogicGate
}

type WorkflowContainer struct {
	WorkFlow
	Steps  []Step                  `json:"steps"`
	Routes []Route                 `json:"routes"`
	Tags   []WorkFlowTagAssociator `json:"tags"`
}

type WorkFlowCloneModel struct {
	ClientId   uuid.UUID               `json:"client_id"`
	ProjectId  uuid.UUID               `json:"project_id"`
	WorkFlowId uuid.UUID               `json:"workFlow_id"`
	Label      string                  `json:"label"`
	Tags       []WorkFlowTagAssociator `json:"tags"`
}

type FluUploadStats struct {
	Status            flu_upload_status.FluUploadStatus
	TotalFluCount     int
	CompletedFluCount int
	ErrorFluCount     int
}
