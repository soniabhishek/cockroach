package flu_validator_repo

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/playment-main/angel/app/DAL/clients"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"testing"
)

func TestSome(t *testing.T) {
	pg := clients.GetPostgresClient()

	var macroTask models.MacroTask

	mm := []interface{}{&macroTask, &macroTask.Project}

	err := pg.SelectOne(mm,
		`select m.*, p.* from macro_tasks m
		inner join projects p on p.id = m.project_id
		where m.id = $1`,
		uuid.FromStringOrNil("6242dc7d-9b61-4f69-855e-52011c97611a"))

	assert.NoError(t, err)
	t.Log(macroTask)
	t.Log(macroTask.Project)
}
