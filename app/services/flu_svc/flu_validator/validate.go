package flu_validator

import (
	"github.com/crowdflux/angel/app/DAL/repositories/flu_validator_repo"
	"github.com/crowdflux/angel/app/models"
)

//Does the data validation for incoming flu
//Third parameter has been called macroTaskId here & it will be called projectId in future once we have refactored the schema
func validateFlu(v flu_validator_repo.IFluValidatorRepo, flu models.FeedLineUnit) (IsValid bool, err error) {

	// fluVs -> flu validators
	fluVs, err := v.GetValidatorsForProject(flu.ProjectId, flu.Tag)

	if err != nil {
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
