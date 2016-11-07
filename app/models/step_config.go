package models

import "github.com/crowdflux/angel/app/models/uuid"

type CrowdsourcingConfig struct {
	MicroTaskId uuid.UUID
	AnswerKey   string
}

type InternalSourcingConfig struct {
	MicroTaskId uuid.UUID
}

type TransformationConfig struct {
	TemplateId string
}

type UnificationConfig struct {
	Multiplication int
}

type BifurcationConfig struct {
	Multiplication int
}

type StartStepConfig struct {
	ImageFieldKey string
}
