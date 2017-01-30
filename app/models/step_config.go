package models

import "github.com/crowdflux/angel/app/models/uuid"

type CrowdsourcingConfig struct {
	MicroTaskId uuid.UUID
	AnswerKey   string
}

type InternalSourcingConfig struct {
	MicroTaskId uuid.UUID
	AnswerKey   string
}

type TransformationConfig struct {
	TemplateId string
}

type ValidationConfig struct {
	TemplateId string
	AnswerKey  string
}

type UnificationConfig struct {
	Multiplication int
}

type BifurcationConfig struct {
	Multiplication int
}

type AlgorithmConfig struct {
	TextFieldKey   string
	AnswerKey      string
	TimeDelayStart float64
	TimeDelayStop  float64
}
