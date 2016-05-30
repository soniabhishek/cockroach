package models

import (
	"database/sql"

	"gitlab.com/playment-main/angel/app/models/uuid"
	"gopkg.in/gorp.v1"
)

// CREATE FUNCTION
// CREATE FUNCTION
type Availabilities struct {
	ID              uuid.UUID      `db:"id" json:"id"`
	EntityId        uuid.UUID      `db:"entity_id" json:"entity_id"`
	EntityType      sql.NullString `db:"entity_type" json:"entity_type"`
	ActivatedAt     gorp.NullTime  `db:"activated_at" json:"activated_at"`
	ActivatorId     uuid.UUID      `db:"activator_id" json:"activator_id"`
	ActivatorType   sql.NullString `db:"activator_type" json:"activator_type"`
	DeactivatedAt   gorp.NullTime  `db:"deactivated_at" json:"deactivated_at"`
	DeactivatorId   uuid.UUID      `db:"deactivator_id" json:"deactivator_id"`
	DeactivatorType sql.NullString `db:"deactivator_type" json:"deactivator_type"`
}

type Comment struct {
	ID         uuid.UUID      `db:"id" json:"id"`
	CreatorId  uuid.UUID      `db:"creator_id" json:"creator_id"`
	Body       sql.NullString `db:"body" json:"body"`
	EntityId   uuid.UUID      `db:"entity_id" json:"entity_id"`
	EntityType sql.NullString `db:"entity_type" json:"entity_type"`
	CreatedAt  gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt  gorp.NullTime  `db:"updated_at" json:"updated_at"`
}

type ContactRequest struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	UserId    uuid.UUID      `db:"user_id" json:"user_id"`
	Email     sql.NullString `db:"email" json:"email"`
	Name      sql.NullString `db:"name" json:"name"`
	Subject   sql.NullString `db:"subject" json:"subject"`
	Message   sql.NullString `db:"message" json:"message"`
	CreatedAt gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime  `db:"updated_at" json:"updated_at"`
}

type CouponTransaction struct {
	ID                 uuid.UUID      `db:"id" json:"id"`
	CouponId           uuid.UUID      `db:"coupon_id" json:"coupon_id"`
	UserId             uuid.UUID      `db:"user_id" json:"user_id"`
	CreatedAt          gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt          gorp.NullTime  `db:"updated_at" json:"updated_at"`
	Count              int            `db:"count" json:"count"`
	IsServed           sql.NullBool   `db:"is_served" json:"is_served"`
	ServedAt           gorp.NullTime  `db:"served_at" json:"served_at"`
	ServedBy           uuid.UUID      `db:"served_by" json:"served_by"`
	EmailId            sql.NullString `db:"email_id" json:"email_id"`
	TransactionDetails JsonFake       `db:"transaction_details" json:"transaction_details"`
	MobileNo           sql.NullString `db:"mobile_no" json:"mobile_no"`
}

type Coupon struct {
	ID                    uuid.UUID     `db:"id" json:"id"`
	Points                int           `db:"points" json:"points"`
	TimesRedeemed         int           `db:"times_redeemed" json:"times_redeemed"`
	IntegrationProviderId uuid.UUID     `db:"integration_provider_id" json:"integration_provider_id"`
	CreatedAt             gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt             gorp.NullTime `db:"updated_at" json:"updated_at"`
	Quantity              int           `db:"quantity" json:"quantity"`
	IsHidden              bool          `db:"is_hidden" json:"is_hidden"`
}

type Email struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	Email     sql.NullString `db:"email" json:"email"`
	UserId    uuid.UUID      `db:"user_id" json:"user_id"`
	CreatedAt gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime  `db:"updated_at" json:"updated_at"`
}

type ExternalAccount struct {
	ID                    uuid.UUID      `db:"id" json:"id"`
	IntegrationProviderId uuid.UUID      `db:"integration_provider_id" json:"integration_provider_id"`
	EmailId               uuid.UUID      `db:"email_id" json:"email_id"`
	CreatedAt             gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt             gorp.NullTime  `db:"updated_at" json:"updated_at"`
	ProfileInfo           JsonFake       `db:"profile_info" json:"profile_info"`
	ExternalId            sql.NullString `db:"external_id" json:"external_id"`
}

type Feedback struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	Subject   string         `db:"subject" json:"subject"`
	UserId    uuid.UUID      `db:"user_id" json:"user_id"`
	Body      sql.NullString `db:"body" json:"body"`
	From      string         `db:"from" json:"from"`
	CreatedAt gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime  `db:"updated_at" json:"updated_at"`
}

type FeedLineUnit struct {
	ID          uuid.UUID     `db:"id" json:"id" bson:"_id"`
	ReferenceId string        `db:"reference_id" json:"reference_id" bson:"reference_id"`
	Data        JsonFake      `db:"data" json:"data" bson:"data"`
	Tag         string        `db:"tag" json:"tag" bson:"tag"`
	MacroTaskId uuid.UUID     `db:"macro_task_id" json:"macro_task_id" bson:"macro_task_id"`
	CreatedAt   gorp.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt   gorp.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	ProjectID uuid.UUID `db:"-"`
}

type ForceUpdateApp struct {
	ID           uuid.UUID      `db:"id" json:"id"`
	Message      string         `db:"message" json:"message"`
	OptionalMin  sql.NullString `db:"optional_min" json:"optional_min"`
	OptionalMax  sql.NullString `db:"optional_max" json:"optional_max"`
	MandatoryMin sql.NullString `db:"mandatory_min" json:"mandatory_min"`
	MandatoryMax sql.NullString `db:"mandatory_max" json:"mandatory_max"`
	IsActive     bool           `db:"is_active" json:"is_active"`
	CreatedAt    gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt    gorp.NullTime  `db:"updated_at" json:"updated_at"`
}

type GrammarElement struct {
	ID             uuid.UUID      `db:"id" json:"id"`
	Name           string         `db:"name" json:"name"`
	Label          string         `db:"label" json:"label"`
	InputTemplate  JsonFake       `db:"input_template" json:"input_template"`
	GrammarVersion sql.NullString `db:"grammar_version" json:"grammar_version"`
	IsDeleted      bool           `db:"is_deleted" json:"is_deleted"`
	Description    string         `db:"description" json:"description"`
	CreatedAt      gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt      gorp.NullTime  `db:"updated_at" json:"updated_at"`
}

type IntegrationProvider struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	Name      string         `db:"name" json:"name"`
	Label     string         `db:"label" json:"label"`
	Website   string         `db:"website" json:"website"`
	CreatedAt gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime  `db:"updated_at" json:"updated_at"`
	LogoUrl   sql.NullString `db:"logo_url" json:"logo_url"`
}

type InvitationRequest struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	CreatedAt gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime  `db:"updated_at" json:"updated_at"`
	Email     sql.NullString `db:"email" json:"email"`
	UserId    uuid.UUID      `db:"user_id" json:"user_id"`
}

type KnexMigration struct {
	ID            int            `db:"id" json:"id"`
	Name          sql.NullString `db:"name" json:"name"`
	Batch         sql.NullInt64  `db:"batch" json:"batch"`
	MigrationTime gorp.NullTime  `db:"migration_time" json:"migration_time"`
}

type KnexMigrationsLock struct {
	IsLocked sql.NullInt64 `db:"is_locked" json:"is_locked"`
}

type MacroTask struct {
	ID        uuid.UUID     `db:"id" json:"id" bson:"_id"`
	Label     string        `db:"label" json:"label" bson:"label"`
	Name      string        `db:"name" json:"name" bson:"name"`
	CreatedAt gorp.NullTime `db:"created_at" json:"created_at" bson:"created_at"`
	UpdatedAt gorp.NullTime `db:"updated_at" json:"updated_at" bson:"updated_at"`
	ProjectId uuid.UUID     `db:"project_id" json:"project_id" bson:"project_id"`
	CreatorId uuid.UUID     `db:"creator_id" json:"creator_id" bson:"creator_id"`
}

type MicroTaskQuestionAssociator struct {
	ID          uuid.UUID     `db:"id" json:"id"`
	MicroTaskId uuid.UUID     `db:"micro_task_id" json:"micro_task_id"`
	QuestionId  uuid.UUID     `db:"question_id" json:"question_id"`
	CreatedAt   gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt   gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type MicroTaskResourceAssociator struct {
	ResourceId  uuid.UUID     `db:"resource_id" json:"resource_id"`
	MicroTaskId uuid.UUID     `db:"micro_task_id" json:"micro_task_id"`
	CreatedAt   gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt   gorp.NullTime `db:"updated_at" json:"updated_at"`
	ID          uuid.UUID     `db:"id" json:"id"`
}

type MicroTaskRewardAssociator struct {
	ID          uuid.UUID     `db:"id" json:"id"`
	MicroTaskId uuid.UUID     `db:"micro_task_id" json:"micro_task_id"`
	RewardId    uuid.UUID     `db:"reward_id" json:"reward_id"`
	CreatedAt   gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt   gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type MicroTask struct {
	ID                  uuid.UUID      `db:"id" json:"id"`
	MacroTaskId         uuid.UUID      `db:"macro_task_id" json:"macro_task_id"`
	CreatorId           uuid.UUID      `db:"creator_id" json:"creator_id"`
	CreatedAt           gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt           gorp.NullTime  `db:"updated_at" json:"updated_at"`
	Name                string         `db:"name" json:"name"`
	Label               string         `db:"label" json:"label"`
	Description         sql.NullString `db:"description" json:"description"`
	MetaData            JsonFake       `db:"meta_data" json:"meta_data"`
	Duration            sql.NullInt64  `db:"duration" json:"duration"`
	Power               sql.NullInt64  `db:"power" json:"power"`
	Points              sql.NullInt64  `db:"points" json:"points"`
	IsDeleted           sql.NullBool   `db:"is_deleted" json:"is_deleted"`
	IsActive            sql.NullBool   `db:"is_active" json:"is_active"`
	FallbackMicroTaskId *uuid.UUID     `db:"fallback_micro_task_id" json:"fallback_micro_task_id"`
	Type                int            `db:"type" json:"type"`
}

type MissionQuestionAssociator struct {
	MissionId  uuid.UUID     `db:"mission_id" json:"mission_id"`
	QuestionId uuid.UUID     `db:"question_id" json:"question_id"`
	CreatedAt  gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt  gorp.NullTime `db:"updated_at" json:"updated_at"`
	ID         uuid.UUID     `db:"id" json:"id"`
}

type MissionSubmission struct {
	ID                         uuid.UUID     `db:"id" json:"id"`
	UserId                     uuid.UUID     `db:"user_id" json:"user_id"`
	MissionId                  uuid.UUID     `db:"mission_id" json:"mission_id"`
	CreatedAt                  gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt                  gorp.NullTime `db:"updated_at" json:"updated_at"`
	CorrectTestQuestionCount   sql.NullInt64 `db:"correct_test_question_count" json:"correct_test_question_count"`
	IncorrectTestQuestionCount sql.NullInt64 `db:"incorrect_test_question_count" json:"incorrect_test_question_count"`
	Status                     int           `db:"status" json:"status"`
}

type Mission struct {
	ID                   uuid.UUID     `db:"id" json:"id"`
	UserId               uuid.UUID     `db:"user_id" json:"user_id"`
	CreatedAt            gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt            gorp.NullTime `db:"updated_at" json:"updated_at"`
	MicroTaskId          uuid.UUID     `db:"micro_task_id" json:"micro_task_id"`
	UiTemplateId         uuid.UUID     `db:"ui_template_id" json:"ui_template_id"`
	SubmissionTemplateId uuid.UUID     `db:"submission_template_id" json:"submission_template_id"`
	Duration             sql.NullInt64 `db:"duration" json:"duration"`
	Power                sql.NullInt64 `db:"power" json:"power"`
	Points               sql.NullInt64 `db:"points" json:"points"`
	GuidelinesId         uuid.UUID     `db:"guidelines_id" json:"guidelines_id"`
	InstructionsId       uuid.UUID     `db:"instructions_id" json:"instructions_id"`
}

type NotificationRecipientAssociator struct {
	ID             uuid.UUID     `db:"id" json:"id"`
	RecipientId    uuid.UUID     `db:"recipient_id" json:"recipient_id"`
	NotificationId uuid.UUID     `db:"notification_id" json:"notification_id"`
	CreatedAt      gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt      gorp.NullTime `db:"updated_at" json:"updated_at"`
	IsRead         sql.NullBool  `db:"is_read" json:"is_read"`
}

type Notification struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	Message   sql.NullString `db:"message" json:"message"`
	CreatedAt gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime  `db:"updated_at" json:"updated_at"`
}

type Permission struct {
	ID         uuid.UUID      `db:"id" json:"id"`
	RoleId     uuid.UUID      `db:"role_id" json:"role_id"`
	Permission sql.NullString `db:"permission" json:"permission"`
	EntityType string         `db:"entity_type" json:"entity_type"`
	EntityId   uuid.UUID      `db:"entity_id" json:"entity_id"`
	CreatedAt  gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt  gorp.NullTime  `db:"updated_at" json:"updated_at"`
}

type PointTransaction struct {
	ID                  uuid.UUID     `db:"id" json:"id"`
	RewardTransactionId uuid.UUID     `db:"reward_transaction_id" json:"reward_transaction_id"`
	UserId              uuid.UUID     `db:"user_id" json:"user_id"`
	CreatedAt           gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt           gorp.NullTime `db:"updated_at" json:"updated_at"`
	AmountCredited      sql.NullInt64 `db:"amount_credited" json:"amount_credited"`
}

type PowerTransaction struct {
	ID                  uuid.UUID     `db:"id" json:"id"`
	RewardTransactionId uuid.UUID     `db:"reward_transaction_id" json:"reward_transaction_id"`
	AmountCredited      sql.NullInt64 `db:"amount_credited" json:"amount_credited"`
	UserId              uuid.UUID     `db:"user_id" json:"user_id"`
	CreatedAt           gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt           gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type Project struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	Label     string        `db:"label" json:"label"`
	Name      string        `db:"name" json:"name"`
	ClientId  *uuid.UUID    `db:"client_id" json:"client_id"`
	CreatorId uuid.UUID     `db:"creator_id" json:"creator_id"`
	StartedAt gorp.NullTime `db:"started_at" json:"started_at"`
	EndedAt   gorp.NullTime `db:"ended_at" json:"ended_at"`
	CreatedAt gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type QuestionSubmission struct {
	ID                  uuid.UUID     `db:"id" json:"id"`
	UserId              uuid.UUID     `db:"user_id" json:"user_id"`
	MissionSubmissionId uuid.UUID     `db:"mission_submission_id" json:"mission_submission_id"`
	QuestionId          uuid.UUID     `db:"question_id" json:"question_id"`
	CreatedAt           gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt           gorp.NullTime `db:"updated_at" json:"updated_at"`
	Confidence          sql.NullInt64 `db:"confidence" json:"confidence"`
	IsTest              bool          `db:"is_test" json:"is_test"`
	Body                JsonFake      `db:"body" json:"body"`
	Status              int           `db:"status" json:"status"`
}

type Question struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	Body      JsonFake      `db:"body" json:"body"`
	CreatedAt gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime `db:"updated_at" json:"updated_at"`
	Label     string        `db:"label" json:"label"`
	IsTest    sql.NullBool  `db:"is_test" json:"is_test"`
	CreatorId uuid.UUID     `db:"creator_id" json:"creator_id"`
	IsActive  sql.NullBool  `db:"is_active" json:"is_active"`
}

type Resources struct {
	ID        uuid.UUID      `db:"id" json:"id"`
	Body      sql.NullString `db:"body" json:"body"`
	BodyType  string         `db:"body_type" json:"body_type"`
	Name      string         `db:"name" json:"name"`
	Label     string         `db:"label" json:"label"`
	CreatedAt gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime  `db:"updated_at" json:"updated_at"`
	CreatorId uuid.UUID      `db:"creator_id" json:"creator_id"`
}

type RewardTransaction struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	UserId    uuid.UUID     `db:"user_id" json:"user_id"`
	RewardId  uuid.UUID     `db:"reward_id" json:"reward_id"`
	MissionId uuid.UUID     `db:"mission_id" json:"mission_id"`
	CreatedAt gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type Reward struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	Points    int           `db:"points" json:"points"`
	Power     int           `db:"power" json:"power"`
	CreatorId uuid.UUID     `db:"creator_id" json:"creator_id"`
	CreatedAt gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type Roles struct {
	ID               uuid.UUID     `db:"id" json:"id"`
	Name             string        `db:"name" json:"name"`
	Label            string        `db:"label" json:"label"`
	ApprovalStrategy int           `db:"approval_strategy" json:"approval_strategy"`
	CreatedAt        gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt        gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type Tag struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	Value     string    `db:"value" json:"value"`
	IsDeleted bool      `db:"is_deleted" json:"is_deleted"`
}

type TagsMicroTaskGroup struct {
	ID          uuid.UUID     `db:"id" json:"id"`
	TagName     string        `db:"tag_name" json:"tag_name"`
	TagValue    string        `db:"tag_value" json:"tag_value"`
	MicroTaskId uuid.UUID     `db:"micro_task_id" json:"micro_task_id"`
	IsDeleted   bool          `db:"is_deleted" json:"is_deleted"`
	CreatedAt   gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt   gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type UserActivity struct {
	ID        uuid.UUID     `db:"id" json:"id"`
	UserId    uuid.UUID     `db:"user_id" json:"user_id"`
	Type      int           `db:"type" json:"type"`
	Body      JsonFake      `db:"body" json:"body"`
	IsDeleted bool          `db:"is_deleted" json:"is_deleted"`
	CreatedAt gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt gorp.NullTime `db:"updated_at" json:"updated_at"`
	IsRead    sql.NullBool  `db:"is_read" json:"is_read"`
}

type UserMicroTaskBlocker struct {
	ID               uuid.UUID     `db:"id" json:"id"`
	UserId           uuid.UUID     `db:"user_id" json:"user_id"`
	MicroTaskId      uuid.UUID     `db:"micro_task_id" json:"micro_task_id"`
	UnblockAfterDays sql.NullInt64 `db:"unblock_after_days" json:"unblock_after_days"`
	CreatedAt        gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt        gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type UserRoleAssociator struct {
	ID               uuid.UUID     `db:"id" json:"id"`
	UserId           uuid.UUID     `db:"user_id" json:"user_id"`
	RoleId           uuid.UUID     `db:"role_id" json:"role_id"`
	ApprovedAt       gorp.NullTime `db:"approved_at" json:"approved_at"`
	ApprovalStrategy int           `db:"approval_strategy" json:"approval_strategy"`
	CreatedAt        gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt        gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type User struct {
	ID                      uuid.UUID      `db:"id" json:"id"`
	Username                string         `db:"username" json:"username"`
	Password                sql.NullString `db:"password" json:"password"`
	CreatedAt               gorp.NullTime  `db:"created_at" json:"created_at"`
	UpdatedAt               gorp.NullTime  `db:"updated_at" json:"updated_at"`
	Gender                  sql.NullString `db:"gender" json:"gender"`
	FirstName               sql.NullString `db:"first_name" json:"first_name"`
	LastName                sql.NullString `db:"last_name" json:"last_name"`
	Locale                  sql.NullString `db:"locale" json:"locale"`
	AvatarUrl               sql.NullString `db:"avatar_url" json:"avatar_url"`
	IncorrectQuestionsCount int            `db:"incorrect_questions_count" json:"incorrect_questions_count"`
	CorrectQuestionsCount   int            `db:"correct_questions_count" json:"correct_questions_count"`
	PendingQuestionsCount   int            `db:"pending_questions_count" json:"pending_questions_count"`
	CoinsCount              int            `db:"coins_count" json:"coins_count"`
	CurrentPower            int            `db:"current_power" json:"current_power"`
	CouponRedeemedCount     int            `db:"coupon_redeemed_count" json:"coupon_redeemed_count"`
	Phone                   sql.NullString `db:"phone" json:"phone"`
	TotalCoinsCount         int            `db:"total_coins_count" json:"total_coins_count"`
}

type FLUValidator struct {
	ID          uuid.UUID     `db:"id" json:"id"`
	FieldName   string        `db:"field_name" json:"field_name"`
	Type        string        `db:"type" json:"type"`
	IsMandatory bool          `db:"is_mandatory" json:"is_mandatory"`
	MacroTaskId uuid.UUID     `db:"macro_task_id" json:"macro_task_id"`
	Tag         string        `db:"tag" json:"tag"`
	CreatedAt   gorp.NullTime `db:"created_at" json:"created_at"`
	UpdatedAt   gorp.NullTime `db:"updated_at" json:"updated_at"`
}

type FluProjectService struct {
	ID          uuid.UUID        `db:"project_id" json:"project_id"`
	Url         string           `db:"url" json:"url"`
	Header	    string           `db:"header" json:"header"`
}
