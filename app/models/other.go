package models

import (
	"database/sql"

	"github.com/lib/pq"

	"gitlab.com/playment-main/angel/app/models/uuid"
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
