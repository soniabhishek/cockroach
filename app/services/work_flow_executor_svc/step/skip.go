package step

import (
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
)

func getFluIdFromString(fluIdsString []string) []uuid.UUID {

	fluIds := make([]uuid.UUID, len(fluIdsString))

	for i, fluIdString := range fluIdsString {

		fluIds[i] = uuid.FromStringOrNil(fluIdString)
	}
	return fluIds
}

var FluIDsToSkip = []uuid.UUID{}

func init() {
	FluIDsToSkip = getFluIdFromString(fluIDToSkipString)
}

func IsSkipped(fluID uuid.UUID) bool {

	for _, fluIDToSkip := range FluIDsToSkip {

		if fluID == fluIDToSkip {
			plog.Info("skipping flus", "fluId: "+fluID.String())
			return true
		}
	}
	return false
}
