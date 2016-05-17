package flu_svc

import (
	"errors"

	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/models/uuid"
)

type extendedFluService struct {
	fluService
}

var _ IValidatorService = &extendedFluService{}

func (i *extendedFluService) GetValidators(macroTaskId uuid.UUID, tag string) (fvs []models.FLUValidator, err error) {

	fvs, err = i.fluValidatorRepo.GetValidatorsForMacroTask(macroTaskId, tag)
	return
}

func (i *extendedFluService) SaveValidator(fv *models.FLUValidator) (err error) {

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

func (i *extendedFluService) DeleteValidator(fv models.FLUValidator) (err error) {

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
	return
}
