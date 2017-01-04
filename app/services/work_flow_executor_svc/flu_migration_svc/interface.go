package flu_migration_svc

import (
	"io"
)

type IFluMigrationService interface {
	GetFluMigrationDetails(fluIdReadCloser io.ReadCloser, migrationRefName string) (fluMigrationCSVDetails FluMigrationCSVDetails, err error)
}
