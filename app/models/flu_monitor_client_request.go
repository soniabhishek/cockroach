package models

import "github.com/crowdflux/angel/app/models/uuid"

type Flu_monitor_client_request struct {
	ID          uuid.UUID   `json:"flu_id"`
	ReferenceId string      `json:"reference_id"`
	Tag         string      `json:"tag"`
	Status      string      `json:"status"`
	Result      interface{} `json:"result"`
}
