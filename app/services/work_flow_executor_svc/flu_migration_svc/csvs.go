package flu_migration_svc

import (
	"encoding/csv"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/support"
	"io"
	"os"
	"strconv"
)

func GetFluIDsFromCSVReader(fluIdReadCloser io.ReadCloser) (uids []uuid.UUID, err error) {

	r := csv.NewReader(fluIdReadCloser)
	defer fluIdReadCloser.Close()

	records, err := r.ReadAll()
	if err != nil {
		return
	}

	// Ignore 1st row as it will contain MasterFluID header
	records = records[1:]

	for _, record := range records {

		uid, err := uuid.FromString(record[0])
		if err != nil {
			return uids, err
		}

		uids = append(uids, uid)
	}
	return
}

type FluMigrationCSVDetails struct {
	CrowdsourcingBufferDeleteFile     *os.File
	CrowdsourcingBufferDeleteFileName string

	UnificationBufferDeleteFile     *os.File
	UnificationBufferDeleteFileName string

	DeactivateFluFile     *os.File
	DeactivateFluFileName string
}

func WriteFluMigrationInfoCSV(fmi FluMigrationInfo, migrationRefName string) (fluMigrationCSVDetails FluMigrationCSVDetails, err error) {

	crowdSourcingBufferToDelete := fmi.FluBufferToDelete[step_type.CrowdSourcing]
	if len(crowdSourcingBufferToDelete) > 0 {

		name := "flu_migration_info_crowd_buffer_to_delete_" + migrationRefName + ".csv"
		crowdBuffDelFile, err := os.Create(support.GetExposedDir() + "/" + name)
		if err != nil {
			return fluMigrationCSVDetails, err
		}
		defer crowdBuffDelFile.Close()
		writeFluInfoCSVFile(crowdSourcingBufferToDelete, crowdBuffDelFile)

		fluMigrationCSVDetails.CrowdsourcingBufferDeleteFile = crowdBuffDelFile
		fluMigrationCSVDetails.CrowdsourcingBufferDeleteFileName = name
	}

	unificationBufferToDelete := fmi.FluBufferToDelete[step_type.CrowdSourcing]
	if len(unificationBufferToDelete) > 0 {
		name := "flu_migration_info_unification_buffer_to_delete_" + migrationRefName + ".csv"
		unificationBuffDelFile, err := os.Create(support.GetExposedDir() + "/")
		if err != nil {
			return fluMigrationCSVDetails, err
		}
		defer unificationBuffDelFile.Close()
		writeFluInfoCSVFile(unificationBufferToDelete, unificationBuffDelFile)

		fluMigrationCSVDetails.CrowdsourcingBufferDeleteFile = unificationBuffDelFile
		fluMigrationCSVDetails.CrowdsourcingBufferDeleteFileName = name
	}

	if len(fmi.FlusToDeactivate) > 0 {
		name := "flu_migration_info_flu_to_deactivate_" + migrationRefName + ".csv"
		deactivateFile, err := os.Create(support.GetExposedDir() + "/")
		if err != nil {
			return fluMigrationCSVDetails, err
		}
		defer deactivateFile.Close()
		writeFluIDsCSVFile(fmi.FlusToDeactivate, deactivateFile)

		fluMigrationCSVDetails.CrowdsourcingBufferDeleteFile = deactivateFile
		fluMigrationCSVDetails.CrowdsourcingBufferDeleteFileName = name
	}

	return
}

func writeFluInfoCSVFile(fluInfos []fluInfo, file *os.File) {
	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"FluID", "StepID", "IsMaster", "MasterID"})
	for _, fluInfo := range fluInfos {
		csvWriter.Write([]string{fluInfo.FluID.String(), fluInfo.StepID.String(), strconv.FormatBool(fluInfo.IsMaster), fluInfo.MasterID.String()})
	}
}

func writeFluIDsCSVFile(fluIDs []uuid.UUID, file *os.File) {
	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"FluID"})
	for _, fluID := range fluIDs {
		csvWriter.Write([]string{fluID.String()})
	}
}
