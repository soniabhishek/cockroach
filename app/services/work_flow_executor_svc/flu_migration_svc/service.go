package flu_migration_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/pkg/errors"
	"io"
)

type fluMigrationService struct {
	fluRepo feed_line_repo.IFluRepo
}

func (fms *fluMigrationService) GetFluMigrationDetails(fluIdReadCloser io.ReadCloser, migrationRefName string) (fluMigrationCSVDetails FluMigrationCSVDetails, err error) {

	masterIDs, err := GetFluIDsFromCSVReader(fluIdReadCloser)
	if err != nil {
		return
	}

	if len(masterIDs) == 0 {
		return fluMigrationCSVDetails, errors.New("Please pass flu_ids in csv")
	}

	fluMigrationInfo, err := FluMigrationHelper(fms.fluRepo, masterIDs)
	if err != nil {
		return
	}

	return WriteFluMigrationInfoCSV(fluMigrationInfo, migrationRefName)
}
