package flu_svc

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/models/uuid"
)

type fakeValidatorRepo struct {
}

func (f *fakeValidatorRepo) GetValidatorsForMacroTask(macroTaskId uuid.UUID, tag string) ([]models.FLUValidator, error) {
	return []models.FLUValidator{
		models.FLUValidator{
			FieldName:   "brand",
			Type:        "STRING",
			IsMandatory: true,
		},
		models.FLUValidator{
			FieldName:   "color",
			Type:        "STRING",
			IsMandatory: true,
		},
		models.FLUValidator{
			FieldName:   "category_id",
			Type:        "STRING",
			IsMandatory: false,
		}}, nil
}

func (f *fakeValidatorRepo) Save(*models.FLUValidator) error {
	return nil
}

func TestValidateFluEmptyValidator(t *testing.T) {

	flu := models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonFake{
			"product_id":  "40843808",
			"category_id": "t_shirt_12",
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       "XYZ",
			"color":       "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, flu)

	assert.NoError(t, err, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")
}

func TestValidateFluPerfectFlu(t *testing.T) {

	flu := models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonFake{
			"product_id":  "40843808",
			"category_id": "t_shirt_12",
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       "XYZ",
			"color":       "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, flu)

	assert.NoError(t, err, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")
}

func TestValidateFluForFieldNotFound(t *testing.T) {

	flu := models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonFake{
			"product_id": "40843808",
			"name":       "XYZ Men's Gold T-Shirt",
			"brand":      "XYZ",
			"color":      "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, flu)
	validationErrs := err.(DataValidationError).Validations

	assert.Error(t, err, "Error occured while validating")
	assert.False(t, isValid, "Expected inValid flu but found valid")
	assert.NotEmpty(t, validationErrs, "Validations errors were empty for inValid flu")
	assert.Equal(t, 1, len(validationErrs), "More than one validation Error found")
	assert.Equal(t, fieldNotFoundVCode, validationErrs[0].ValidationCode, fieldNotFoundVCode+" error was expected")
	assert.Equal(t, []string{"category_id"}, validationErrs[0].MetaData.Fields, "only category_id was expected")
}

func TestValidateFluForWrongDataType(t *testing.T) {

	flu := models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonFake{
			"product_id":  "40843808",
			"category_id": "t_shirt_12",
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       []string{"123", "1233"},
			"color":       "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, flu)
	validationErrs := err.(DataValidationError).Validations

	assert.Error(t, err, "Error occured while validating")
	assert.False(t, isValid, "Expected inValid flu but found valid")
	assert.NotEmpty(t, validationErrs, "Validations errors were empty for inValid flu")
	assert.Equal(t, 1, len(validationErrs), "More than one validation Error found")
	assert.Equal(t, wrongDataTypeVCode, validationErrs[0].ValidationCode, wrongDataTypeVCode+" error was expected")
	assert.Equal(t, []string{"brand"}, validationErrs[0].MetaData.Fields, "only brand was expected")
}

func TestValidateFluForMandatoryField(t *testing.T) {

	flu := models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonFake{
			"product_id":  "40843808",
			"category_id": "t_shirt_12",
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       "XYZ",
			"color":       "",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, flu)
	validationErrs := err.(DataValidationError).Validations

	assert.Error(t, err, "Error occured while validating")
	assert.False(t, isValid, "Expected inValid flu but found valid")
	assert.NotEmpty(t, validationErrs, "Validations errors were empty for inValid flu")
	assert.Equal(t, 1, len(validationErrs), "More than one validation Error found")
	assert.Equal(t, mandatoryFieldEmptyVCode, validationErrs[0].ValidationCode, mandatoryFieldEmptyVCode+" error was expected")
	assert.Equal(t, []string{"color"}, validationErrs[0].MetaData.Fields, "only color was expected")
}
