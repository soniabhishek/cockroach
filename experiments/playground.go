package main

import (
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_output"
)

func main() {

	stepId := uuid.FromStringOrNil("43dacf34-83fc-4628-b276-23f05d23e6f7")
	projectId := uuid.FromStringOrNil("dff44216-f42e-484d-aade-8004ccb1bf79")

	flu_output.ForceSendBackInQps(stepId, projectId, 1)

}
