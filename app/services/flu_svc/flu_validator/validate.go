package flu_validator

import (
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/repositories/flu_validator_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_errors"
	"github.com/crowdflux/angel/app/services/plerrors"
	"strings"
)

//Does the data validation for incoming flu
//Third parameter has been called macroTaskId here & it will be called projectId in future once we have refactored the schema
func validateFlu(v flu_validator_repo.IFluValidatorRepo, fluOb *models.FeedLineUnit) (IsValid bool, err error) {

	// fluVs -> flu validators
	flu := *fluOb
	fluVs, err := v.GetValidatorsForProject(flu.ProjectId, flu.Tag)

	if err != nil {
		return
	}

	// Validating ReferenceID
	if flu.ReferenceId == "" {
		err = flu_errors.ErrReferenceIdMissing
		return
	}

	// Validating TAG
	if flu.Tag == "" {
		err = flu_errors.ErrTagMissing
		return
	}

	// Validating Data
	if flu.Data == nil {
		err = flu_errors.ErrDataMissing
		return
	}

	fieldNotFound := validationError{ValidationCode: fieldNotFoundVCode}
	wrongDataType := validationError{ValidationCode: wrongDataTypeVCode}
	mandatoryFieldEmpty := validationError{ValidationCode: mandatoryFieldEmptyVCode}

	for _, fluV := range fluVs {

		name := fluV.FieldName

		// Check if field value is present or not
		fieldVal, ok := flu.Data[name]
		if !ok {
			fieldNotFound.AddMetaDataField(name)
			continue
		}

		// Check if field value is string or not
		fieldValStr, ok := fieldVal.(string)
		if !ok {
			wrongDataType.AddMetaDataField(name)
			continue
		}

		// Check if field is mandatory & not empty
		if fluV.IsMandatory && fieldValStr == "" {
			mandatoryFieldEmpty.AddMetaDataField(name)
			continue
		}
	}

	success := true
	var vErrs []validationError

	// Loop over all the possible errors
	for _, v := range []validationError{fieldNotFound, wrongDataType, mandatoryFieldEmpty} {

		// Check if any error occurred
		if len(v.MetaData.Fields) > 0 {
			vErrs = append(vErrs, v)
			success = false
		}
	}

	err = imageUrlEncryptor(fluOb, fluVs)
	if err != nil {
		plog.Error("image encryption error", err)
	}

	if success {
		return true, nil
	} else {
		return false, DataValidationError{errDataValidation, vErrs}
	}

}

//--------------------------------------------------------------------------------//

type validationMetaData struct {
	Fields []string `json:"fields"`
}

type validationError struct {
	ValidationCode string             `json:"validation_code"`
	MetaData       validationMetaData `json:"meta_data"`
}

func (v *validationError) AddMetaDataField(field string) {
	v.MetaData.Fields = append(v.MetaData.Fields, field)
}

//--------------------------------------------------------------------------------//

const fieldNotFoundVCode = "FIELD_NOT_FOUND"
const wrongDataTypeVCode = "WRONG_DATA_TYPE"
const mandatoryFieldEmptyVCode = "MANDATORY_FIELD_EMPTY"

func imageUrlEncryptor(flu *models.FeedLineUnit, input_config []models.FLUValidator) (err error) {
	img_config := ""
	for _, item := range input_config {
		if strings.ToLower(item.Type) == "image_array" {
			if img_config != "" {
				err = plerrors.ServiceError{"GE_0002", "Multiple image_url configurations for Project and tag combination"}
				return
			}
			img_config = item.FieldName
		}
	}

	if img_config == "" || flu.Data[img_config] == nil {
		err = plerrors.ServiceError{"GE_0003", "Invalid image_url config for the flu received"}
		return
	}

	var img_urls = flu.Data[img_config].([]string)

	if err != nil || len(img_urls) == 0 {
		return flu_errors.ErrDataMissing
	}

	//Image encryption
	urlSlice, err := GetEncryptedUrls(img_urls)

	if err != nil {
		return
	}

	flu.Data.Merge(models.JsonF{img_config: urlSlice})
	return
}

func GetEncryptedUrls(imageField []string) (urlSlice []string, err error) {

	var encResult map[string]clients.LuigiResponse
	encResult, err = clients.GetLuigiClient().GetEncryptedUrls(imageField)
	if err != nil {
		return
	}
	for _, item := range imageField {
		returnItem := encResult[item]
		if returnItem.Value == false {
			err = flu_errors.ErrImageNotValid
			return
		}

		urlSlice = append(urlSlice, returnItem.PlaymentUrl)
	}
	return
}
