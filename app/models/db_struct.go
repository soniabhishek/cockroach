package models

import (
	"database/sql"

	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/lib/pq"
)

type Availabilities struct {
	ID              uuid.UUID      `db:"id" json:"id" bson:"_id"`
	EntityId        uuid.UUID      `db:"entity_id" json:"entity_id" bson:"entity_id"`
	EntityType      sql.NullString `db:"entity_type" json:"entity_type" bson:"entity_type"`
	ActivatedAt     pq.NullTime    `db:"activated_at" json:"activated_at" bson:"activated_at"`
	ActivatorId     uuid.UUID      `db:"activator_id" json:"activator_id" bson:"activator_id"`
	ActivatorType   sql.NullString `db:"activator_type" json:"activator_type" bson:"activator_type"`
	DeactivatedAt   pq.NullTime    `db:"deactivated_at" json:"deactivated_at" bson:"deactivated_at"`
	DeactivatorId   uuid.UUID      `db:"deactivator_id" json:"deactivator_id" bson:"deactivator_id"`
	DeactivatorType sql.NullString `db:"deactivator_type" json:"deactivator_type" bson:"deactivator_type"`
}

type BatchProces struct {
	ID          uuid.UUID      `db:"id" json:"id" bson:"_id"`
	Name        sql.NullString `db:"name" json:"name" bson:"name"`
	Done        sql.NullBool   `db:"done" json:"done" bson:"done"`
	Aborted     sql.NullBool   `db:"aborted" json:"aborted" bson:"aborted"`
	Completion  sql.NullInt64  `db:"completion" json:"completion" bson:"completion"`
	Type        sql.NullInt64  `db:"type" json:"type" bson:"type"`
	CreatedAt   pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
	MacroTaskId uuid.UUID      `db:"macro_task_id" json:"macro_task_id" bson:"macro_task_id"`
}

type Client struct {
	ID               uuid.UUID   `db:"id" json:"id" bson:"_id"`
	UserId           uuid.UUID   `db:"user_id" json:"user_id" bson:"user_id"`
	ClientSecretUuid uuid.UUID   `db:"client_secret_uuid" json:"client_secret_uuid" bson:"client_secret_uuid"`
	CreatedAt        pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt        pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	Options          JsonF       `db:"options" json:"options" bson:"options"`
}

type Comment struct {
	ID         uuid.UUID      `db:"id" json:"id" bson:"_id"`
	CreatorId  uuid.UUID      `db:"creator_id" json:"creator_id" bson:"creator_id"`
	Body       sql.NullString `db:"body" json:"body" bson:"body"`
	EntityId   uuid.UUID      `db:"entity_id" json:"entity_id" bson:"entity_id"`
	EntityType sql.NullString `db:"entity_type" json:"entity_type" bson:"entity_type"`
	CreatedAt  pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt  pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type ContactRequest struct {
	ID        uuid.UUID      `db:"id" json:"id" bson:"_id"`
	UserId    uuid.UUID      `db:"user_id" json:"user_id" bson:"user_id"`
	Email     sql.NullString `db:"email" json:"email" bson:"email"`
	Name      sql.NullString `db:"name" json:"name" bson:"name"`
	Subject   sql.NullString `db:"subject" json:"subject" bson:"subject"`
	Message   sql.NullString `db:"message" json:"message" bson:"message"`
	CreatedAt pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type CouponTransaction struct {
	ID                 uuid.UUID      `db:"id" json:"id" bson:"_id"`
	CouponId           uuid.UUID      `db:"coupon_id" json:"coupon_id" bson:"coupon_id"`
	UserId             uuid.UUID      `db:"user_id" json:"user_id" bson:"user_id"`
	CreatedAt          pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt          pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
	Count              int            `db:"count" json:"count" bson:"count"`
	IsServed           sql.NullBool   `db:"is_served" json:"is_served" bson:"is_served"`
	ServedAt           pq.NullTime    `db:"served_at" json:"served_at" bson:"served_at"`
	ServedBy           uuid.UUID      `db:"served_by" json:"served_by" bson:"served_by"`
	EmailId            sql.NullString `db:"email_id" json:"email_id" bson:"email_id"`
	TransactionDetails JsonF          `db:"transaction_details" json:"transaction_details" bson:"transaction_details"`
	MobileNo           sql.NullString `db:"mobile_no" json:"mobile_no" bson:"mobile_no"`
	IsReverted         sql.NullBool   `db:"is_reverted" json:"is_reverted" bson:"is_reverted"`
}

type Coupon struct {
	ID                    uuid.UUID   `db:"id" json:"id" bson:"_id"`
	Points                int         `db:"points" json:"points" bson:"points"`
	TimesRedeemed         int         `db:"times_redeemed" json:"times_redeemed" bson:"times_redeemed"`
	IntegrationProviderId uuid.UUID   `db:"integration_provider_id" json:"integration_provider_id" bson:"integration_provider_id"`
	CreatedAt             pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt             pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	Quantity              int         `db:"quantity" json:"quantity" bson:"quantity"`
	IsHidden              bool        `db:"is_hidden" json:"is_hidden" bson:"is_hidden"`
}

type CrowdsourcingFluBuffer struct {
	ID         uuid.UUID   `db:"id" json:"id" bson:"_id"`
	FluId      uuid.UUID   `db:"flu_id" json:"flu_id" bson:"flu_id"`
	QuestionId uuid.UUID   `db:"question_id" json:"question_id" bson:"question_id"`
	CreatedAt  pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt  pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type CrowdsourcingStepConfiguration struct {
	ID          uuid.UUID   `db:"id" json:"id" bson:"_id"`
	StepId      uuid.UUID   `db:"step_id" json:"step_id" bson:"step_id"`
	MicroTaskId uuid.UUID   `db:"micro_task_id" json:"micro_task_id" bson:"micro_task_id"`
	CreatedAt   pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type TransformationStepConfiguration struct {
	ID         int         `db:"id" json:"id" bson:"_id"`
	StepId     uuid.UUID   `db:"step_id" json:"step_id" bson:"step_id"`
	TemplateId string      `db:"template_id" json:"template_id" bson:"template_id"`
	CreatedAt  pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt  pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type Email struct {
	ID        uuid.UUID      `db:"id" json:"id" bson:"_id"`
	Email     sql.NullString `db:"email" json:"email" bson:"email"`
	UserId    uuid.UUID      `db:"user_id" json:"user_id" bson:"user_id"`
	CreatedAt pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type ExternalAccount struct {
	ID                    uuid.UUID      `db:"id" json:"id" bson:"_id"`
	IntegrationProviderId uuid.UUID      `db:"integration_provider_id" json:"integration_provider_id" bson:"integration_provider_id"`
	EmailId               uuid.UUID      `db:"email_id" json:"email_id" bson:"email_id"`
	CreatedAt             pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt             pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
	ProfileInfo           JsonF          `db:"profile_info" json:"profile_info" bson:"profile_info"`
	ExternalId            sql.NullString `db:"external_id" json:"external_id" bson:"external_id"`
	UserId                uuid.UUID      `db:"user_id" json:"user_id" bson:"user_id"`
}

type FeedLineUnit struct {
	ID          uuid.UUID   `db:"id" json:"id" bson:"_id"`
	ReferenceId string      `db:"reference_id" json:"reference_id" bson:"reference_id"`
	Data        JsonF       `db:"data" json:"data" bson:"data"`
	Tag         string      `db:"tag" json:"tag" bson:"tag"`
	CreatedAt   pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	StepId      uuid.UUID   `db:"step_id" json:"step_id" bson:"step_id"`
	Build       JsonF       `db:"build" json:"build" bson:"build"`
	ProjectId   uuid.UUID   `db:"project_id" json:"project_id" bson:"project_id"`
}

type FeedLineLog struct {
	ID         int            `db:"id" json:"id" bson:"_id"`
	FluId      uuid.UUID      `db:"flu_id" json:"flu_id" bson:"flu_id"`
	Message    sql.NullString `db:"message" json:"message" bson:"message"`
	MetaData   JsonF          `db:"meta_data" json:"meta_data" bson:"meta_data"`
	StepType   sql.NullInt64  `db:"step_type" json:"step_type" bson:"step_type"`
	StepEntry  sql.NullBool   `db:"step_entry" json:"step_entry" bson:"step_entry"`
	StepExit   sql.NullBool   `db:"step_exit" json:"step_exit" bson:"step_exit"`
	StepId     uuid.UUID      `db:"step_id" json:"step_id" bson:"step_id"`
	WorkFlowId uuid.UUID      `db:"work_flow_id" json:"work_flow_id" bson:"work_flow_id"`
	CreatedAt  pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
}

type Feedback struct {
	ID        uuid.UUID      `db:"id" json:"id" bson:"_id"`
	Subject   string         `db:"subject" json:"subject" bson:"subject"`
	UserId    uuid.UUID      `db:"user_id" json:"user_id" bson:"user_id"`
	Body      sql.NullString `db:"body" json:"body" bson:"body"`
	From      string         `db:"from" json:"from" bson:"from"`
	CreatedAt pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type ForceUpdateApp struct {
	ID           uuid.UUID      `db:"id" json:"id" bson:"_id"`
	Message      string         `db:"message" json:"message" bson:"message"`
	OptionalMin  sql.NullString `db:"optional_min" json:"optional_min" bson:"optional_min"`
	OptionalMax  sql.NullString `db:"optional_max" json:"optional_max" bson:"optional_max"`
	MandatoryMin sql.NullString `db:"mandatory_min" json:"mandatory_min" bson:"mandatory_min"`
	MandatoryMax sql.NullString `db:"mandatory_max" json:"mandatory_max" bson:"mandatory_max"`
	IsActive     bool           `db:"is_active" json:"is_active" bson:"is_active"`
	CreatedAt    pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt    pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type GrammarElement struct {
	ID             uuid.UUID      `db:"id" json:"id" bson:"_id"`
	Name           string         `db:"name" json:"name" bson:"name"`
	Label          string         `db:"label" json:"label" bson:"label"`
	InputTemplate  JsonF          `db:"input_template" json:"input_template" bson:"input_template"`
	GrammarVersion sql.NullString `db:"grammar_version" json:"grammar_version" bson:"grammar_version"`
	IsDeleted      bool           `db:"is_deleted" json:"is_deleted" bson:"is_deleted"`
	Description    string         `db:"description" json:"description" bson:"description"`
	CreatedAt      pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt      pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type ImageDictionary struct {
	ID        uuid.UUID      `db:"id" json:"id" bson:"_id"`
	RealUrl   string         `db:"real_url" json:"real_url" bson:"real_url"`
	CloudUrl  string         `db:"cloud_url" json:"cloud_url" bson:"cloud_url"`
	Extra     sql.NullString `db:"extra" json:"extra" bson:"extra"`
	CreatedAt pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type FLUValidator struct {
	ID          uuid.UUID   `db:"id" json:"id" bson:"_id"`
	FieldName   string      `db:"field_name" json:"field_name" bson:"field_name"`
	Type        string      `db:"type" json:"type" bson:"type"`
	IsMandatory bool        `db:"is_mandatory" json:"is_mandatory" bson:"is_mandatory"`
	Tag         string      `db:"tag" json:"tag" bson:"tag"`
	CreatedAt   pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	ProjectId   uuid.UUID   `db:"project_id" json:"project_id" bson:"project_id"`
}

type IntegrationProvider struct {
	ID        uuid.UUID      `db:"id" json:"id" bson:"_id"`
	Name      string         `db:"name" json:"name" bson:"name"`
	Label     string         `db:"label" json:"label" bson:"label"`
	Website   string         `db:"website" json:"website" bson:"website"`
	CreatedAt pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
	LogoUrl   sql.NullString `db:"logo_url" json:"logo_url" bson:"logo_url"`
}

type InvitationRequest struct {
	ID        uuid.UUID      `db:"id" json:"id" bson:"_id"`
	CreatedAt pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
	Email     sql.NullString `db:"email" json:"email" bson:"email"`
	UserId    uuid.UUID      `db:"user_id" json:"user_id" bson:"user_id"`
}

type KnexMigration struct {
	ID            int            `db:"id" json:"id" bson:"_id"`
	Name          sql.NullString `db:"name" json:"name" bson:"name"`
	Batch         sql.NullInt64  `db:"batch" json:"batch" bson:"batch"`
	MigrationTime pq.NullTime    `db:"migration_time" json:"migration_time" bson:"migration_time"`
}

type KnexMigrationsLock struct {
	IsLocked sql.NullInt64 `db:"is_locked" json:"is_locked" bson:"is_locked"`
}

type LogicGate struct {
	ID            uuid.UUID `db:"id" json:"id" bson:"_id"`
	InputTemplate JsonF     `db:"input_template" json:"input_template" bson:"input_template"`
	Formula       int       `db:"formula" json:"formula" bson:"formula"`
}

type LogicGateFormula struct {
	ID   int            `db:"id" json:"id" bson:"_id"`
	Name sql.NullString `db:"name" json:"name" bson:"name"`
}

type MacroTask struct {
	ID        uuid.UUID   `db:"id" json:"id" bson:"_id"`
	Label     string      `db:"label" json:"label" bson:"label"`
	Name      string      `db:"name" json:"name" bson:"name"`
	CreatedAt pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	ProjectId uuid.UUID   `db:"project_id" json:"project_id" bson:"project_id"`
	CreatorId uuid.UUID   `db:"creator_id" json:"creator_id" bson:"creator_id"`
}

type MicroTaskQuestionAssociator struct {
	ID          uuid.UUID   `db:"id" json:"id" bson:"_id"`
	MicroTaskId uuid.UUID   `db:"micro_task_id" json:"micro_task_id" bson:"micro_task_id"`
	QuestionId  uuid.UUID   `db:"question_id" json:"question_id" bson:"question_id"`
	CreatedAt   pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type MicroTaskResourceAssociator struct {
	ResourceId  uuid.UUID   `db:"resource_id" json:"resource_id" bson:"resource_id"`
	MicroTaskId uuid.UUID   `db:"micro_task_id" json:"micro_task_id" bson:"micro_task_id"`
	CreatedAt   pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	ID          uuid.UUID   `db:"id" json:"id" bson:"_id"`
}

type MicroTaskRewardAssociator struct {
	ID          uuid.UUID   `db:"id" json:"id" bson:"_id"`
	MicroTaskId uuid.UUID   `db:"micro_task_id" json:"micro_task_id" bson:"micro_task_id"`
	RewardId    uuid.UUID   `db:"reward_id" json:"reward_id" bson:"reward_id"`
	CreatedAt   pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type MicroTaskType struct {
	ID   int    `db:"id" json:"id" bson:"_id"`
	Name string `db:"name" json:"name" bson:"name"`
}

type MicroTask struct {
	ID                  uuid.UUID      `db:"id" json:"id" bson:"_id"`
	MacroTaskId         uuid.UUID      `db:"macro_task_id" json:"macro_task_id" bson:"macro_task_id"`
	CreatorId           uuid.UUID      `db:"creator_id" json:"creator_id" bson:"creator_id"`
	CreatedAt           pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt           pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
	Name                string         `db:"name" json:"name" bson:"name"`
	Label               string         `db:"label" json:"label" bson:"label"`
	Description         sql.NullString `db:"description" json:"description" bson:"description"`
	MetaData            JsonF          `db:"meta_data" json:"meta_data" bson:"meta_data"`
	Duration            sql.NullInt64  `db:"duration" json:"duration" bson:"duration"`
	Power               sql.NullInt64  `db:"power" json:"power" bson:"power"`
	Points              sql.NullInt64  `db:"points" json:"points" bson:"points"`
	IsDeleted           sql.NullBool   `db:"is_deleted" json:"is_deleted" bson:"is_deleted"`
	IsActive            sql.NullBool   `db:"is_active" json:"is_active" bson:"is_active"`
	FallbackMicroTaskId uuid.UUID      `db:"fallback_micro_task_id" json:"fallback_micro_task_id" bson:"fallback_micro_task_id"`
	Type                sql.NullInt64  `db:"type" json:"type" bson:"type"`
}

type MissionQuestionAssociator struct {
	MissionId  uuid.UUID   `db:"mission_id" json:"mission_id" bson:"mission_id"`
	QuestionId uuid.UUID   `db:"question_id" json:"question_id" bson:"question_id"`
	CreatedAt  pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt  pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	ID         uuid.UUID   `db:"id" json:"id" bson:"_id"`
}

type MissionSubmission struct {
	ID                         uuid.UUID     `db:"id" json:"id" bson:"_id"`
	UserId                     uuid.UUID     `db:"user_id" json:"user_id" bson:"user_id"`
	MissionId                  uuid.UUID     `db:"mission_id" json:"mission_id" bson:"mission_id"`
	CreatedAt                  pq.NullTime   `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt                  pq.NullTime   `db:"updated_at" json:"updated_at" bson:"updated_at"`
	CorrectTestQuestionCount   sql.NullInt64 `db:"correct_test_question_count" json:"correct_test_question_count" bson:"correct_test_question_count"`
	IncorrectTestQuestionCount sql.NullInt64 `db:"incorrect_test_question_count" json:"incorrect_test_question_count" bson:"incorrect_test_question_count"`
	Status                     int           `db:"status" json:"status" bson:"status"`
}

type Mission struct {
	ID                   uuid.UUID     `db:"id" json:"id" bson:"_id"`
	UserId               uuid.UUID     `db:"user_id" json:"user_id" bson:"user_id"`
	CreatedAt            pq.NullTime   `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt            pq.NullTime   `db:"updated_at" json:"updated_at" bson:"updated_at"`
	MicroTaskId          uuid.UUID     `db:"micro_task_id" json:"micro_task_id" bson:"micro_task_id"`
	UiTemplateId         uuid.UUID     `db:"ui_template_id" json:"ui_template_id" bson:"ui_template_id"`
	SubmissionTemplateId uuid.UUID     `db:"submission_template_id" json:"submission_template_id" bson:"submission_template_id"`
	Duration             sql.NullInt64 `db:"duration" json:"duration" bson:"duration"`
	Power                sql.NullInt64 `db:"power" json:"power" bson:"power"`
	Points               sql.NullInt64 `db:"points" json:"points" bson:"points"`
	GuidelinesId         uuid.UUID     `db:"guidelines_id" json:"guidelines_id" bson:"guidelines_id"`
	InstructionsId       uuid.UUID     `db:"instructions_id" json:"instructions_id" bson:"instructions_id"`
}

type NotificationRecipientAssociator struct {
	ID             uuid.UUID    `db:"id" json:"id" bson:"_id"`
	RecipientId    uuid.UUID    `db:"recipient_id" json:"recipient_id" bson:"recipient_id"`
	NotificationId uuid.UUID    `db:"notification_id" json:"notification_id" bson:"notification_id"`
	CreatedAt      pq.NullTime  `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt      pq.NullTime  `db:"updated_at" json:"updated_at" bson:"updated_at"`
	IsRead         sql.NullBool `db:"is_read" json:"is_read" bson:"is_read"`
}

type Notification struct {
	ID        uuid.UUID      `db:"id" json:"id" bson:"_id"`
	Message   sql.NullString `db:"message" json:"message" bson:"message"`
	CreatedAt pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type Permission struct {
	ID         uuid.UUID      `db:"id" json:"id" bson:"_id"`
	RoleId     uuid.UUID      `db:"role_id" json:"role_id" bson:"role_id"`
	Permission sql.NullString `db:"permission" json:"permission" bson:"permission"`
	EntityType string         `db:"entity_type" json:"entity_type" bson:"entity_type"`
	EntityId   uuid.UUID      `db:"entity_id" json:"entity_id" bson:"entity_id"`
	CreatedAt  pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt  pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type PointTransaction struct {
	ID                  uuid.UUID     `db:"id" json:"id" bson:"_id"`
	RewardTransactionId uuid.UUID     `db:"reward_transaction_id" json:"reward_transaction_id" bson:"reward_transaction_id"`
	UserId              uuid.UUID     `db:"user_id" json:"user_id" bson:"user_id"`
	CreatedAt           pq.NullTime   `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt           pq.NullTime   `db:"updated_at" json:"updated_at" bson:"updated_at"`
	AmountCredited      sql.NullInt64 `db:"amount_credited" json:"amount_credited" bson:"amount_credited"`
}

type PowerTransaction struct {
	ID                  uuid.UUID     `db:"id" json:"id" bson:"_id"`
	RewardTransactionId uuid.UUID     `db:"reward_transaction_id" json:"reward_transaction_id" bson:"reward_transaction_id"`
	AmountCredited      sql.NullInt64 `db:"amount_credited" json:"amount_credited" bson:"amount_credited"`
	UserId              uuid.UUID     `db:"user_id" json:"user_id" bson:"user_id"`
	CreatedAt           pq.NullTime   `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt           pq.NullTime   `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type ProjectConfiguration struct {
	ProjectId   uuid.UUID   `db:"project_id" json:"project_id" bson:"project_id"`
	PostBackUrl string      `db:"post_back_url" json:"post_back_url" bson:"post_back_url"`
	Headers     JsonF       `db:"headers" json:"headers" bson:"headers"`
	CreatedAt   pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	Options     JsonF       `db:"options" json:"options" bson:"options"`
}

type Project struct {
	ID        uuid.UUID   `db:"id" json:"id" bson:"_id"`
	Label     string      `db:"label" json:"label" bson:"label"`
	Name      string      `db:"name" json:"name" bson:"name"`
	ClientId  uuid.UUID   `db:"client_id" json:"client_id" bson:"client_id"`
	CreatorId uuid.UUID   `db:"creator_id" json:"creator_id" bson:"creator_id"`
	StartedAt pq.NullTime `db:"started_at" json:"started_at" bson:"started_at"`
	EndedAt   pq.NullTime `db:"ended_at" json:"ended_at" bson:"ended_at"`
	CreatedAt pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type QuestionAnswer struct {
	ID          uuid.UUID   `db:"id" json:"id" bson:"_id"`
	QuestionId  uuid.UUID   `db:"question_id" json:"question_id" bson:"question_id"`
	MicroTaskId uuid.UUID   `db:"micro_task_id" json:"micro_task_id" bson:"micro_task_id"`
	Body        JsonF       `db:"body" json:"body" bson:"body"`
	CreatedAt   pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type QuestionSubmission struct {
	ID                  uuid.UUID     `db:"id" json:"id" bson:"_id"`
	UserId              uuid.UUID     `db:"user_id" json:"user_id" bson:"user_id"`
	MissionSubmissionId uuid.UUID     `db:"mission_submission_id" json:"mission_submission_id" bson:"mission_submission_id"`
	QuestionId          uuid.UUID     `db:"question_id" json:"question_id" bson:"question_id"`
	CreatedAt           pq.NullTime   `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt           pq.NullTime   `db:"updated_at" json:"updated_at" bson:"updated_at"`
	Confidence          sql.NullInt64 `db:"confidence" json:"confidence" bson:"confidence"`
	IsTest              bool          `db:"is_test" json:"is_test" bson:"is_test"`
	Body                JsonF         `db:"body" json:"body" bson:"body"`
	Status              int           `db:"status" json:"status" bson:"status"`
}

type Question struct {
	ID        uuid.UUID    `db:"id" json:"id" bson:"_id"`
	Body      JsonF        `db:"body" json:"body" bson:"body"`
	CreatedAt pq.NullTime  `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime  `db:"updated_at" json:"updated_at" bson:"updated_at"`
	Label     string       `db:"label" json:"label" bson:"label"`
	IsTest    sql.NullBool `db:"is_test" json:"is_test" bson:"is_test"`
	CreatorId uuid.UUID    `db:"creator_id" json:"creator_id" bson:"creator_id"`
	IsActive  sql.NullBool `db:"is_active" json:"is_active" bson:"is_active"`
}

type Resources struct {
	ID        uuid.UUID      `db:"id" json:"id" bson:"_id"`
	Body      sql.NullString `db:"body" json:"body" bson:"body"`
	BodyType  string         `db:"body_type" json:"body_type" bson:"body_type"`
	Name      string         `db:"name" json:"name" bson:"name"`
	Label     string         `db:"label" json:"label" bson:"label"`
	CreatedAt pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
	CreatorId uuid.UUID      `db:"creator_id" json:"creator_id" bson:"creator_id"`
}

type RewardTransaction struct {
	ID        uuid.UUID   `db:"id" json:"id" bson:"_id"`
	UserId    uuid.UUID   `db:"user_id" json:"user_id" bson:"user_id"`
	RewardId  uuid.UUID   `db:"reward_id" json:"reward_id" bson:"reward_id"`
	MissionId uuid.UUID   `db:"mission_id" json:"mission_id" bson:"mission_id"`
	CreatedAt pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type Reward struct {
	ID        uuid.UUID   `db:"id" json:"id" bson:"_id"`
	Points    int         `db:"points" json:"points" bson:"points"`
	Power     int         `db:"power" json:"power" bson:"power"`
	CreatorId uuid.UUID   `db:"creator_id" json:"creator_id" bson:"creator_id"`
	CreatedAt pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type Roles struct {
	ID               uuid.UUID   `db:"id" json:"id" bson:"_id"`
	Name             string      `db:"name" json:"name" bson:"name"`
	Label            string      `db:"label" json:"label" bson:"label"`
	ApprovalStrategy int         `db:"approval_strategy" json:"approval_strategy" bson:"approval_strategy"`
	CreatedAt        pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt        pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type Route struct {
	ID          uuid.UUID    `db:"id" json:"id" bson:"_id"`
	StepId      uuid.UUID    `db:"step_id" json:"step_id" bson:"step_id"`
	LogicGateId uuid.UUID    `db:"logic_gate_id" json:"logic_gate_id" bson:"logic_gate_id"`
	NextStepId  uuid.UUID    `db:"next_step_id" json:"next_step_id" bson:"next_step_id"`
	IsDeleted   sql.NullBool `db:"is_deleted" json:"is_deleted" bson:"is_deleted"`
	CreatedAt   pq.NullTime  `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime  `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type Step struct {
	ID         uuid.UUID          `db:"id" json:"id" bson:"_id"`
	Type       step_type.StepType `db:"type" json:"type" bson:"type"`
	WorkFlowId uuid.UUID          `db:"work_flow_id" json:"work_flow_id" bson:"work_flow_id"`
	IsDeleted  sql.NullBool       `db:"is_deleted" json:"is_deleted" bson:"is_deleted"`
	CreatedAt  pq.NullTime        `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt  pq.NullTime        `db:"updated_at" json:"updated_at" bson:"updated_at"`
	IsStart    bool               `db:"is_start" json:"is_start" bson:"is_start"`
}

type StepTypes struct {
	ID   int            `db:"id" json:"id" bson:"_id"`
	Name sql.NullString `db:"name" json:"name" bson:"name"`
}

type Tag struct {
	ID        uuid.UUID `db:"id" json:"id" bson:"_id"`
	Name      string    `db:"name" json:"name" bson:"name"`
	Value     string    `db:"value" json:"value" bson:"value"`
	IsDeleted bool      `db:"is_deleted" json:"is_deleted" bson:"is_deleted"`
}

type TagsMicroTaskGroup struct {
	ID          uuid.UUID   `db:"id" json:"id" bson:"_id"`
	TagName     string      `db:"tag_name" json:"tag_name" bson:"tag_name"`
	TagValue    string      `db:"tag_value" json:"tag_value" bson:"tag_value"`
	MicroTaskId uuid.UUID   `db:"micro_task_id" json:"micro_task_id" bson:"micro_task_id"`
	IsDeleted   bool        `db:"is_deleted" json:"is_deleted" bson:"is_deleted"`
	CreatedAt   pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type Task struct {
	ID          uuid.UUID    `db:"id" json:"id" bson:"_id"`
	IsActive    sql.NullBool `db:"is_active" json:"is_active" bson:"is_active"`
	IsDeleted   sql.NullBool `db:"is_deleted" json:"is_deleted" bson:"is_deleted"`
	MacroTaskId uuid.UUID    `db:"macro_task_id" json:"macro_task_id" bson:"macro_task_id"`
	CreatedAt   pq.NullTime  `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   pq.NullTime  `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type UserActivity struct {
	ID        uuid.UUID    `db:"id" json:"id" bson:"_id"`
	UserId    uuid.UUID    `db:"user_id" json:"user_id" bson:"user_id"`
	Type      int          `db:"type" json:"type" bson:"type"`
	Body      JsonF        `db:"body" json:"body" bson:"body"`
	IsDeleted bool         `db:"is_deleted" json:"is_deleted" bson:"is_deleted"`
	CreatedAt pq.NullTime  `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime  `db:"updated_at" json:"updated_at" bson:"updated_at"`
	IsRead    sql.NullBool `db:"is_read" json:"is_read" bson:"is_read"`
}

type UserMicroTaskBlocker struct {
	ID               uuid.UUID     `db:"id" json:"id" bson:"_id"`
	UserId           uuid.UUID     `db:"user_id" json:"user_id" bson:"user_id"`
	MicroTaskId      uuid.UUID     `db:"micro_task_id" json:"micro_task_id" bson:"micro_task_id"`
	UnblockAfterDays sql.NullInt64 `db:"unblock_after_days" json:"unblock_after_days" bson:"unblock_after_days"`
	CreatedAt        pq.NullTime   `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt        pq.NullTime   `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type UserRoleAssociator struct {
	ID               uuid.UUID   `db:"id" json:"id" bson:"_id"`
	UserId           uuid.UUID   `db:"user_id" json:"user_id" bson:"user_id"`
	RoleId           uuid.UUID   `db:"role_id" json:"role_id" bson:"role_id"`
	ApprovedAt       pq.NullTime `db:"approved_at" json:"approved_at" bson:"approved_at"`
	ApprovalStrategy int         `db:"approval_strategy" json:"approval_strategy" bson:"approval_strategy"`
	CreatedAt        pq.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt        pq.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
}

type User struct {
	ID                      uuid.UUID      `db:"id" json:"id" bson:"_id"`
	Username                string         `db:"username" json:"username" bson:"username"`
	Password                sql.NullString `db:"password" json:"password" bson:"password"`
	CreatedAt               pq.NullTime    `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt               pq.NullTime    `db:"updated_at" json:"updated_at" bson:"updated_at"`
	Gender                  sql.NullString `db:"gender" json:"gender" bson:"gender"`
	FirstName               sql.NullString `db:"first_name" json:"first_name" bson:"first_name"`
	LastName                sql.NullString `db:"last_name" json:"last_name" bson:"last_name"`
	Locale                  sql.NullString `db:"locale" json:"locale" bson:"locale"`
	AvatarUrl               sql.NullString `db:"avatar_url" json:"avatar_url" bson:"avatar_url"`
	IncorrectQuestionsCount int            `db:"incorrect_questions_count" json:"incorrect_questions_count" bson:"incorrect_questions_count"`
	CorrectQuestionsCount   int            `db:"correct_questions_count" json:"correct_questions_count" bson:"correct_questions_count"`
	PendingQuestionsCount   int            `db:"pending_questions_count" json:"pending_questions_count" bson:"pending_questions_count"`
	CoinsCount              int            `db:"coins_count" json:"coins_count" bson:"coins_count"`
	CurrentPower            int            `db:"current_power" json:"current_power" bson:"current_power"`
	CouponRedeemedCount     int            `db:"coupon_redeemed_count" json:"coupon_redeemed_count" bson:"coupon_redeemed_count"`
	Phone                   sql.NullString `db:"phone" json:"phone" bson:"phone"`
	TotalCoinsCount         int            `db:"total_coins_count" json:"total_coins_count" bson:"total_coins_count"`
}

type WorkFlow struct {
	ID        uuid.UUID    `db:"id" json:"id" bson:"_id"`
	ProjectId uuid.UUID    `db:"project_id" json:"project_id" bson:"project_id"`
	IsDeleted sql.NullBool `db:"is_deleted" json:"is_deleted" bson:"is_deleted"`
	CreatedAt pq.NullTime  `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt pq.NullTime  `db:"updated_at" json:"updated_at" bson:"updated_at"`
}
