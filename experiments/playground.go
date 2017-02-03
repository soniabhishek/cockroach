package main

import (
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_output"
)

func main() {

	stepId := uuid.FromStringOrNil("057b5478-78a5-4fc4-8e67-b8db2c2e9490")
	projectId := uuid.FromStringOrNil("b256f731-2eca-4aaa-99d8-83ca24bd6b2c")
	flu_output.ForceSendBackInQps(stepId, projectId, 10)

	//stepId = uuid.FromStringOrNil("6a6662ba-685b-4b35-9057-0dbd493451cc")
	//projectId = uuid.FromStringOrNil("b256f731-2eca-4aaa-99d8-83ca24bd6b2c")
	//flu_output.ForceSendBackInQps(stepId, projectId, 10)

	//stepId = uuid.FromStringOrNil("d977389b-bf55-4cc3-914d-8c2b4eb1f8d5")
	//projectId = uuid.FromStringOrNil("b256f731-2eca-4aaa-99d8-83ca24bd6b2c")
	//flu_output.ForceSendBackInQps(stepId, projectId, 10)

}
