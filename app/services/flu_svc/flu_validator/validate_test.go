package flu_validator

import (
	"testing"

	"github.com/asaskevich/govalidator"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/stretchr/testify/assert"
)

type fakeValidatorRepo struct {
}

func (f *fakeValidatorRepo) GetValidatorsForProject(projectId uuid.UUID, tag string) ([]models.FLUValidator, error) {
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
			FieldName:   "image_url",
			Type:        "IMAGE_ARRAY",
			IsMandatory: false,
		},
		models.FLUValidator{
			FieldName:   "category_id",
			Type:        "STRING",
			IsMandatory: false,
		}}, nil

}

var image_url_valid = []string{"https://s3-ap-southeast-1.amazonaws.com/playmentproduction/public/B00X0X3AKG_2.jpg", "https://s3-ap-southeast-1.amazonaws.com/playmentproduction/public/B00PU0DELW_2.jpg"}

func (f *fakeValidatorRepo) Save(*models.FLUValidator) error {
	return nil
}

func TestValidateFluEmptyValidator(t *testing.T) {
	flu := models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonF{
			"product_id":  "40843808",
			"category_id": "t_shirt_12",
			"image_url":   image_url_valid,
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       "XYZ",
			"color":       "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, &flu)

	assert.NoError(t, err, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")
}

func TestValidateFluPerfectFlu(t *testing.T) {
	flu := models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonF{
			"product_id":  "40843808",
			"category_id": "t_shirt_12",
			"image_url":   image_url_valid,
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       "XYZ",
			"color":       "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, &flu)

	assert.NoError(t, err, "Error occured while validating")
	assert.True(t, isValid, "Expected valid flu but found inValid")
	assert.Empty(t, err, "Validations errors were non-empty for valid flu")
}

func TestValidateFluForFieldNotFound(t *testing.T) {

	flu := models.FeedLineUnit{
		ReferenceId: "PAYTM_123",
		Data: models.JsonF{
			"product_id": "40843808",
			"name":       "XYZ Men's Gold T-Shirt",
			"image_url":  image_url_valid,
			"brand":      "XYZ",
			"color":      "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, &flu)
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
		Data: models.JsonF{
			"product_id":  "40843808",
			"category_id": "t_shirt_12",
			"image_url":   image_url_valid,
			"name":        "XYZ Men's Gold T-Shirt",
			"brand":       []string{"123", "1233"},
			"color":       "Gold",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, &flu)
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
		Data: models.JsonF{
			"product_id":  "40843808",
			"image_url":   image_url_valid,
			"name":        "XYZ Men's Gold T-Shirt",
			"category_id": "t_shirt_12",
			"brand":       "XYZ",
			"color":       "",
		},
		Tag: "PAYTM_TSHIRT",
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, &flu)
	validationErrs := err.(DataValidationError).Validations

	assert.Error(t, err, "Error occured while validating")
	assert.False(t, isValid, "Expected inValid flu but found valid")
	assert.NotEmpty(t, validationErrs, "Validations errors were empty for inValid flu")
	assert.Equal(t, 1, len(validationErrs), "More than one validation Error found")
	assert.Equal(t, mandatoryFieldEmptyVCode, validationErrs[0].ValidationCode, mandatoryFieldEmptyVCode+" error was expected")
	assert.Equal(t, []string{"color"}, validationErrs[0].MetaData.Fields, "only color was expected")
}

func TestEncryptionForValidImageUrls(t *testing.T) {

	var fluId = uuid.NewV4()

	var flu = models.FeedLineUnit{
		ID:          fluId,
		ReferenceId: "Ola123",
		Tag:         "Ola",
		Data: models.JsonF{
			"brand":       "Sony",
			"category_id": "t_shirt_12",
			"color":       "Gold",
			"image_url":   image_url_valid},
		Build: models.JsonF{},
	}

	isValid, err := validateFlu(&fakeValidatorRepo{}, &flu)

	returnedUrlList := flu.Data["image_url"].([]string)
	assert.True(t, govalidator.IsURL(returnedUrlList[0]))
	assert.True(t, govalidator.IsURL(returnedUrlList[1]))
	assert.True(t, isValid)
	assert.Nil(t, err)
	assert.EqualValues(t, len(returnedUrlList), 2)
	assert.NotEqual(t, returnedUrlList, image_url_valid)

}

func Test_for_invalid_urls(t *testing.T) {

	var fluId = uuid.NewV4()
	var image_url_invalid = []string{"https://s3-ap-southea-1.amazonaws.com/playmentproduction/public/B00X0X3AKG_2.jpg", "https://s3-ap-southeast-1.amazonaws.com/playmentproduction/public/B00PU0DELW_2.jpg"}

	var flu = models.FeedLineUnit{
		ID:          fluId,
		ReferenceId: "Ola123",
		Tag:         "Ola",
		Data: models.JsonF{
			"brand":       "Sony",
			"category_id": "t_shirt_12",
			"color":       "Gold",
			"image_url":   image_url_invalid},
		Build: models.JsonF{},
	}

	returnedUrlList := flu.Data["image_url"].([]string)
	isValid, err := validateFlu(&fakeValidatorRepo{}, &flu)
	validationErrs := err.(DataValidationError).Validations

	assert.Error(t, err, "Error occured while validating")
	assert.False(t, isValid, "Expected inValid flu but found valid")
	assert.NotEmpty(t, validationErrs, "Validations errors were empty for inValid flu")
	assert.Equal(t, 1, len(validationErrs), "More than one validation Error found")
	assert.Equal(t, invalidImageLinkVCode, validationErrs[0].ValidationCode, invalidImageLinkVCode+" error was expected")
	assert.Equal(t, []string{"image_url"}, validationErrs[0].MetaData.Fields, "only image_url was expected")
	assert.Equal(t, returnedUrlList, image_url_invalid)
}
