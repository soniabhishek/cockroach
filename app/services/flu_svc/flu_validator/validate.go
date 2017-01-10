package flu_validator

import (
	"github.com/asaskevich/govalidator"
	"github.com/crowdflux/angel/app/DAL/clients"
	"github.com/crowdflux/angel/app/DAL/repositories/flu_validator_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/plog/log_tags"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_errors"
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
	invalidImageLink := validationError{ValidationCode: invalidImageLinkVCode}

	for _, fluV := range fluVs {

		name := fluV.FieldName

		// Check if field value is present or not
		fieldVal, ok := flu.Data[name]
		if !ok {
			if fluV.IsMandatory {
				mandatoryFieldEmpty.AddMetaDataField(name)
			}
			continue
		}

		switch fluV.Type {
		case "STRING":
			// Check if field value is string or not
			_, ok := fieldVal.(string)
			if !ok {
				wrongDataType.AddMetaDataField(name)
				continue
			}

		case "IMAGE_ARRAY":

			// Check if field value is string or not
			fieldValArray, ok := fieldVal.([]interface{})
			if !ok {
				wrongDataType.AddMetaDataField(name)
				continue
			}
			fieldValImgArray := make([]string, len(fieldValArray))
			success := true
			for i, val := range fieldValArray {

				fieldValImgArray[i], ok = val.(string)
				if !ok || !govalidator.IsURL(fieldValImgArray[i]) {
					invalidImageLink.AddMetaDataField(name)
					success = false
					break
				}

			}
			if !success {
				continue
			}
			// Check if field is mandatory & not empty
			if fluV.IsMandatory && len(fieldValImgArray) == 0 {
				mandatoryFieldEmpty.AddMetaDataField(name)
				continue
			}

			//Image encryption
			encUrls, err := GetEncryptedUrls(fieldValImgArray)
			if err != nil {
				plog.Error("Flu Validator", err, plog.Message("Error in Luigi Encryption"), plog.MessageWithParam(log_tags.FLU, flu))
				invalidImageLink.AddMetaDataField(name)
				continue
			}

			//Edit the flu
			flu.Build[name] = encUrls

		case "IMAGE":
			// Check if field value is string or not
			fieldValString, ok := fieldVal.(string)
			if !ok {
				wrongDataType.AddMetaDataField(name)
				continue
			}

			if !govalidator.IsURL(fieldValString) {
				invalidImageLink.AddMetaDataField(name)
				continue
			}

			strArray := []string{fieldValString}

			//Image encryption
			encUrls, err := GetEncryptedUrls(strArray)
			if err != nil {
				invalidImageLink.AddMetaDataField(name)
				continue
			}

			//Edit the flu
			flu.Build[name] = encUrls[0]

		}
	}

	success := true
	var vErrs []validationError

	// Loop over all the possible errors
	for _, v := range []validationError{fieldNotFound, wrongDataType, mandatoryFieldEmpty, invalidImageLink} {

		// Check if any error occurred
		if len(v.MetaData.Fields) > 0 {
			vErrs = append(vErrs, v)
			success = false
		}
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
const invalidImageLinkVCode = "INVALID_IMAGE_LINK"

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
