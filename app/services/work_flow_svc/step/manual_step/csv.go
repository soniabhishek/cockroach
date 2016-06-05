package manual_step

import (
	"encoding/csv"
	"errors"
	"github.com/lib/pq"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/utilities"
	"os"
	"time"
)

const timeFormat = time.RFC3339

func DownloadCsv(manualStepId uuid.UUID, projectId uuid.UUID) (file string) {
	flRepo := feed_line_repo.New()
	flus, err := flRepo.GetByStepId(manualStepId)
	if err != nil {
		plog.Error("Manual Step", err, manualStepId)
		return
	}
	plog.Info("manual step flus going to be downloaded", flus, manualStepId)

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
		plog.Error("Manual Step", err, manualStepId, file)
		return
	}
	defer f.Close()

	// Write Unmarshaled json data to CSV file
	w := csv.NewWriter(f)
	record := []string{"Id", "ReferenceId", "Data", "Build", "Tag", "ProjectId", "StepId", "CreatedAt", "UpdatedAt"}
	w.Write(record)

	for _, obj := range flus {
		record = make([]string, 0)
		record = append(record, obj.ID.String())
		record = append(record, obj.ReferenceId)
		record = append(record, obj.Data.String())
		record = append(record, obj.Build.String())
		record = append(record, obj.Tag)
		record = append(record, obj.ProjectId.String())
		record = append(record, obj.StepId.String())

		if obj.CreatedAt.Valid {
			record = append(record, obj.CreatedAt.Time.Format(timeFormat))
		} else {
			record = append(record, "")
		}

		if obj.UpdatedAt.Valid {
			record = append(record, obj.UpdatedAt.Time.Format(timeFormat))
		} else {
			record = append(record, "")
		}

		w.Write(record)
	}
	w.Flush()
	return
}

func UploadCsv(filepath string, manualStepId uuid.UUID, projectId uuid.UUID) error {
	file := filepath + manualStepId.String() + utilities.Hyphen + projectId.String() + ".csv"
	csvfile, err := os.Open(file)
	if err != nil {
		plog.Error("Manual Step", err, "csv opening")
		return err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = -1 // see the Reader struct information below
	rawCSVdata, err := reader.ReadAll()

	if err != nil {
		plog.Error("Manual Step", err, "csv reading error")
		return err
	}

	flus := make([]models.FeedLineUnit, 0)
	// sanity check, display to standard output
	for i, row := range rawCSVdata {
		if i != 0 {

			// CSV Headers should be :
			// "Id", "ReferenceId", "Data", "Build", "Tag", "ProjectId", "StepId", "CreatedAt"
			id, err := uuid.FromString(row[0])
			if err != nil {
				return err
			}
			referenceId := row[1]
			if referenceId == "" {
				return errors.New("reference id cant be empty")
			}
			// skip the data scanning part in future
			// user should not be able to change data
			data := models.JsonFake{}
			if err := data.Scan(row[2]); err != nil {
				return err
			}

			build := models.JsonFake{}
			if err := build.Scan(row[3]); err != nil {
				return err
			}

			tag := row[4]
			if tag == "" {
				return errors.New("tag cant be empty")
			}

			projectId, err := uuid.FromString(row[5])
			if err != nil {
				return err
			}

			stepId, err := uuid.FromString(row[6])
			if err != nil {
				return err
			}

			createdAtTime, err := time.Parse(timeFormat, row[7])
			if err != nil {
				return errors.New("Time format not valid")
			}
			createdAt := pq.NullTime{createdAtTime, true}

			updatedAt := pq.NullTime{time.Now(), true}

			flu := models.FeedLineUnit{
				ID:          id,
				ReferenceId: referenceId,
				Data:        data,
				Build:       build,
				Tag:         tag,
				ProjectId:   projectId,
				StepId:      stepId,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}
			flus = append(flus, flu)
		}
	}
	plog.Info("Manual Step", "Flus going to be updated from csv upload", flus)

	flRepo := feed_line_repo.New()
	err = flRepo.BulkUpdate(flus)
	return err
}
