package flu_svc

import (
	"gitlab.com/playment-main/angel/app/DAL/repositories/flu_validator_repo"
	"gitlab.com/playment-main/angel/app/DAL/repositories/macro_task_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

//Does the data validation for incoming flu
//Third parameter has been called macroTaskId here & it will be called projectId in future once we have refactored the schema
func validateFlu(v flu_validator_repo.IFluValidatorRepo, flu models.FeedLineUnit) (IsValid bool, err error) {

	fluVs, err := v.GetValidatorsForMacroTask(flu.MacroTaskId, flu.Tag)

	if err != nil {
		return
	}

	fieldNotFound := validationError{ValidationCode: fieldNotFoundVCode}
	wrongDataType := validationError{ValidationCode: wrongDataTypeVCode}
	mandatoryFieldEmpty := validationError{ValidationCode: mandatoryFieldEmptyVCode}

	for _, fluV := range fluVs {

		name := fluV.FieldName

		fieldVal, ok := flu.Data[name]

		if !ok {
			fieldNotFound.AddMetaDataField(name)
			continue
		}

		fieldValStr, ok := fieldVal.(string)

		if !ok {
			wrongDataType.AddMetaDataField(name)
			continue
		}

		if fluV.IsMandatory && fieldValStr == "" {
			mandatoryFieldEmpty.AddMetaDataField(name)
			continue
		}
	}

	success := true
	var vErrs []validationError

	for _, v := range []validationError{fieldNotFound, wrongDataType, mandatoryFieldEmpty} {
		if len(v.MetaData.Fields) > 0 {
			vErrs = append(vErrs, v)
			success = false
		}
	}

	if success {
		return true, nil
	} else {
		return false, DataValidationError{ErrDataValidation, vErrs}
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

//--------------------------------------------------------------------------------//
//CHECK MACRO_TASK_ID
//--------------------------------------------------------------------------------//

func checkMacroTaskExists(r macro_task_repo.IMacroTaskRepo, mId uuid.UUID) error {
	_, err := r.Get(mId)
	return err
}
