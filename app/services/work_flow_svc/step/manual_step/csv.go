package manual_step

import (
	"encoding/csv"
	"errors"
	"fmt"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/utilities"
	"io"
	"os"
	"time"
)

const timeFormat = time.RFC3339

func DownloadCsv(manualStepId uuid.UUID) (string, error) {
	flRepo := feed_line_repo.New()
	flus, err := flRepo.GetByStepId(manualStepId)
	if err != nil {
		plog.Error("Manual Step", err, manualStepId)
		return utilities.Empty, err
	}
	plog.Info("manual step flus going to be downloaded", flus, manualStepId)

	path := config.Get(config.DOWNLOAD_PATH)
	file := path + string(os.PathSeparator) + manualStepId.String() + ".csv"
	err = createFile(file)
	if err != nil {
		return utilities.Empty, nil
	}

	// Write unmarshaled json data to CSV file
	writeBuff := [][]string{{"Id", "ReferenceId", "Data", "Build", "Tag", "ProjectId", "StepId", "CreatedAt", "UpdatedAt"}}

	for _, obj := range flus {
		record := make([]string, 0)
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

		writeBuff = append(writeBuff, record)
	}
	err = writeCSV(file, writeBuff)
	if err != nil {
		return utilities.Empty, err
	}
	return file, nil
}

func UploadCsv(filename string) error {
	file := "/tmp/" + string(os.PathSeparator) + filename
	csvfile, err := os.Open(file)
	if err != nil {
		plog.Error("Manual Step", err, "csv opening")
		return err
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.FieldsPerRecord = 2 // so the reader will always check how many records are present in each row

	flus := make([]models.FeedLineUnit, 0)
	var cnt int = -1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			plog.Error("Manual Step", err, " csv reading error")
			return err
		}

		cnt++
		if cnt == 0 {
			continue
		}
		flu, err := getFlu(row)
		if err != nil {
			plog.Error("Manual Step", err, " csv reading error")
			return err
		}

		flus = append(flus, flu)
	}
	plog.Info("Total lines to read", len(flus))
	plog.Info("Manual Step", "Flus going to be updated from csv upload", flus)

	flRepo := feed_line_repo.New()
	err = flRepo.BulkUpdate(flus)
	fmt.Println(err)
	return err
}

func getFlu(row []string) (flu models.FeedLineUnit, err error) {
	id, err := uuid.FromString(row[0])
	if err != nil {
		plog.Error("Error ID:", err)
		return flu, errors.New("ID is not valid. [" + row[0] + "]")
	}

	build := models.JsonFake{}
	if err := build.Scan(row[1]); err != nil {
		plog.Error("Error Build:", err)
		return flu, errors.New("Build field is not valid. [" + row[1] + "]")
	}

	flu = models.FeedLineUnit{
		ID:    id,
		Build: build,
	}
	return flu, nil
}

func createFile(filepath string) error {
	// detect if file exists
	var _, err = os.Stat(filepath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

func writeCSV(filepath string, records [][]string) error {
	var file, err = os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			plog.Error("Error while writing CSV", err)
			return err
		}
	}
	writer.Flush()
	return nil
}
