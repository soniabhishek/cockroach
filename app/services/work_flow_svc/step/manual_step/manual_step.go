package manual_step

import (
	"encoding/csv"
	"fmt"
	"github.com/lib/pq"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/feed_line"
	"gitlab.com/playment-main/angel/app/services/work_flow_svc/step"
	"gitlab.com/playment-main/angel/utilities"
	"os"
	"reflect"
	"time"
)

type manualStep struct {
	step.Step
	id int
}

func (ms manualStep) processFlu(flu feed_line.FLU) {
	ms.saveJson(flu)
}

func (ms *manualStep) saveJson(flu feed_line.FLU) {
	flus := models.FeedLineUnit{ID: flu.ID,
		ReferenceId: flu.ReferenceId,
		Data:        flu.Data,
		Tag:         flu.Tag,
		MacroTaskId: flu.MacroTaskId,
		CreatedAt:   flu.CreatedAt,
		UpdatedAt:   flu.UpdatedAt,
		StepId:      flu.StepId,
	}
	flRepo := feed_line_repo.New()
	flRepo.Update(flus)
}

func (ms *manualStep) DownloadCsv(manualStepId uuid.UUID, projectId uuid.UUID) (file string) {
	flRepo := feed_line_repo.New()
	flus, err := flRepo.GetByStepId(manualStepId)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(flus)

	//TODO get the data from DB
	/*stream := "json data"
	data := []byte(stream)
	// Unmarshal JSON data
	var flus []feed_line.FLU
	err := json.Unmarshal([]byte(data), &flus)
	if err != nil {
		fmt.Println(err)
	}*/

	//TODO get the path
	path := "./"
	file = path + manualStepId.String() + utilities.Hyphen + projectId.String() + ".csv"
	// Create a csv file
	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	// Write Unmarshaled json data to CSV file
	w := csv.NewWriter(f)
	var record []string

	s := reflect.ValueOf(&models.FeedLineUnit{}).Elem()
	typeOfT := s.Type()
	for i := 0; i < s.NumField(); i++ {
		record = append(record, typeOfT.Field(i).Name)
	}
	w.Write(record)

	for _, obj := range flus {
		record = make([]string, 0)
		record = append(record, obj.ID.String())
		record = append(record, obj.ReferenceId)
		record = append(record, obj.Data.String())

		record = append(record, obj.Tag)
		record = append(record, obj.MacroTaskId.String())

		record = append(record, obj.CreatedAt.Time.String())
		record = append(record, obj.UpdatedAt.Time.String())

		w.Write(record)
	}
	w.Flush()
	return
}

func (ma *manualStep) UploadCsv(filepath string, manualStepId uuid.UUID, projectId uuid.UUID) {
	file := filepath + manualStepId.String() + utilities.Hyphen + projectId.String() + ".csv"
	csvfile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below
	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	flus := make([]models.FeedLineUnit, 0)
	// sanity check, display to standard output
	for i, each := range rawCSVdata {
		if i != 0 {
			id, _ := uuid.FromString(each[0])
			macroTaskId, _ := uuid.FromString(each[4])
			loc, _ := time.LoadLocation("India")
			ptime, _ := pq.ParseTimestamp(loc, each[5])
			createdTime := pq.NullTime{Time: ptime, Valid: true}
			ptime, _ = pq.ParseTimestamp(loc, each[6])
			updatedTime := pq.NullTime{Time: ptime, Valid: true}
			flu := models.FeedLineUnit{
				ID:          id,
				ReferenceId: each[1],
				//Data:        each[2],
				Data:        models.JsonFake{},
				Tag:         each[3],
				MacroTaskId: macroTaskId,
				CreatedAt:   createdTime,
				UpdatedAt:   updatedTime,
			}
			flus = append(flus, flu)
		}
	}
	fmt.Println(flus)

	flRepo := feed_line_repo.New()
	for _, each := range flus {
		flRepo.Update(each)
	}
}
