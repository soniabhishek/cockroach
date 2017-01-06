package flu_migration_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/pkg/errors"
)

type fluInfo struct {
	FluID    uuid.UUID
	StepID   uuid.UUID
	IsMaster bool
	MasterID uuid.UUID
}

type stepTypeFluMap map[step_type.StepType][]fluInfo

type FluMigrationInfo struct {
	FluBufferToDelete stepTypeFluMap
	FlusToDeactivate  []uuid.UUID
}

func FluMigrationHelper(fluRepo feed_line_repo.IFluRepo, masterFluIDs []uuid.UUID) (fmi FluMigrationInfo, err error) {

	masterFlusWithStep, err := fluRepo.GetByIDs(masterFluIDs)
	if err != nil {
		return fmi, err
	}

	if len(masterFlusWithStep) == 0 {
		return fmi, errors.New("No flus found in db from passed IDs")
	}

	for _, flu := range masterFlusWithStep {
		if !flu.IsMaster {
			return fmi, errors.New("Child flu found in Input IDs. FluID: " + flu.ID.String())
		}
	}

	childFlusWithStep, err := fluRepo.GetChildFLusByMasterIDs(masterFluIDs)
	if err != nil {
		return fmi, err
	}

	masterStepFluMap := GetStepFluMap(masterFlusWithStep)
	childStepFluMap := GetStepFluMap(childFlusWithStep)

	crowdSourcingFluBufferToDelete := append(masterStepFluMap[step_type.CrowdSourcing], childStepFluMap[step_type.CrowdSourcing]...)
	unificationFluBufferToDelete := append(masterStepFluMap[step_type.Unification], childStepFluMap[step_type.Unification]...)

	fluBufferToDelete := make(stepTypeFluMap)
	fluBufferToDelete[step_type.CrowdSourcing] = crowdSourcingFluBufferToDelete
	fluBufferToDelete[step_type.Unification] = unificationFluBufferToDelete

	flusToDeactivate := make([]uuid.UUID, len(childFlusWithStep))
	for i, flu := range childFlusWithStep {
		flusToDeactivate[i] = flu.ID
	}

	return FluMigrationInfo{fluBufferToDelete, flusToDeactivate}, nil
}

func GetStepFluMap(flusWithStep []models.FluWithStep) stepTypeFluMap {

	stepFluIDMap := make(stepTypeFluMap)

	for _, flu := range flusWithStep {
		fluInfos, ok := stepFluIDMap[flu.Step.Type]
		if !ok {
			stepFluIDMap[flu.Step.Type] = []fluInfo{{flu.ID, flu.StepId, flu.IsMaster, flu.MasterId}}
		} else {
			fluInfos = append(fluInfos, fluInfo{flu.ID, flu.StepId, flu.IsMaster, flu.MasterId})
			stepFluIDMap[flu.Step.Type] = fluInfos
		}
	}
	return stepFluIDMap
}
