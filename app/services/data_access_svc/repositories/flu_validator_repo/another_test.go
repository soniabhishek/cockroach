package flu_validator_repo

import (
	"fmt"
	"gitlab.com/playment-main/support/app/models"
	"gitlab.com/playment-main/support/app/services/data_access_svc/clients"
	"testing"
)

func TestSome(t *testing.T) {
	pg := clients.GetPostgresClient()

	m := models.MacroTask{}

	err := pg.SelectOne(&m, `select macro_tasks.*, projects.* from macro_tasks
	inner join projects on projects.id = macro_tasks.project_id
	limit 1`)

	if err != nil {
		t.Fatalf("%#v", err)
	} else {
		fmt.Printf("%#v", &m)

	}

}
