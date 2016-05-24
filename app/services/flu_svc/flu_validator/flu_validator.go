package flu_validator

import (
	"errors"

	"gitlab.com/playment-main/angel/app/DAL/repositories/flu_validator_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
)

type fluValidator struct {
	fluValidatorRepo flu_validator_repo.IFluValidatorRepo
}

var _ IFluValidatorService = &fluValidator{}

func (i *fluValidator) GetValidators(macroTaskId uuid.UUID, tag string) (fvs []models.FLUValidator, err error) {

	fvs, err = i.fluValidatorRepo.GetValidatorsForMacroTask(macroTaskId, tag)
	return
}

func (i *fluValidator) SaveValidator(fv *models.FLUValidator) (err error) {

	existingFvs, err := i.GetValidators(fv.MacroTaskId, fv.Tag)

	if err != nil {
		return err
	}

	for _, existingFv := range existingFvs {
		if existingFv.FieldName == fv.FieldName {
			fv.ID = existingFv.ID
			break
		}
	}

	if fv.FieldName == "" {
		return errors.New("Field Name cant be empty")
	}

	if fv.Type != "STRING" {
		return errors.New("Invalid type")
	}

	if fv.Tag == "" {
		return errors.New("Tag cant be empty")
	}

	err = i.fluValidatorRepo.Save(fv)
	if err != nil {
		return
	}
	return nil
}

// not implemented
func (i *fluValidator) DeleteValidator(fv models.FLUValidator) (err error) {

	if fv.MacroTaskId == uuid.Nil {
		return errors.New("MacroTaskId cant be empty")
	}

	if fv.FieldName == "" {
		return errors.New("Field Name cant be empty")
	}

	if fv.Tag == "" {
		return errors.New("Tag cant be empty")
	}

	//err = i.fluValidatorRepo.Delete(fv)
	panic("Not implemented")
	return
}

func (i *fluValidator) Validate(flu models.FeedLineUnit) (IsValid bool, err error) {
	return validateFlu(i.fluValidatorRepo, flu)
}
