package manual_step

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"io"
	"os"
	"time"
)

const timeFormat = time.RFC3339

func DownloadCsv(manualStepId uuid.UUID) (file string) {
	flRepo := feed_line_repo.New()
	flus, err := flRepo.GetByStepId(manualStepId)
	if err != nil {
		plog.Error("Manual Step", err, manualStepId)
		return
	}
	plog.Info("manual step flus going to be downloaded", flus, manualStepId)

	path := config.Get(config.DOWNLOAD_PATH)
	file = path + string(os.PathSeparator) + manualStepId.String() + ".csv"
	createFile(file)

	// Write Unmarshaled json data to CSV file
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
	writeCSV(file, writeBuff)
	return
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
	// TODO Refactor the code below
	// CSV Headers should be :
	// "Id", "ReferenceId", "Data", "Build", "Tag", "ProjectId", "StepId", "CreatedAt"
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

func createFile(filepath string) {
	// detect if file exists
	var _, err = os.Stat(filepath)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(filepath)
		checkError(err)
		defer file.Close()
	}
}

func writeFile(filepath string, buffer []string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(filepath, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// write some text to file
	for _, line := range buffer {
		_, err = file.WriteString(line)
		checkError(err)
	}

	// save changes
	err = file.Sync()
	checkError(err)
}

func readFile(filepath string) {
	// re-open file
	var file, err = os.OpenFile(filepath, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()

	// read file
	var text = make([]byte, 1024)
	for {
		n, err := file.Read(text)
		if err != io.EOF {
			checkError(err)
		}
		if n == 0 {
			break
		}
	}
	fmt.Println(string(text))
	checkError(err)
}

func deleteFile(filepath string) {
	// delete file
	var err = os.Remove(filepath)
	checkError(err)
}

func writeCSV(filepath string, records [][]string) {
	var file, err = os.OpenFile(filepath, os.O_RDWR, 0644)
	checkError(err)
	defer file.Close()
	writer := csv.NewWriter(file)
	for _, record := range records {
		err := writer.Write(record)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
	writer.Flush()
}

func checkError(err error) {
	if err != nil {
		plog.Error("Error File Operations:", err)
		//os.Exit(0)
	}
}

func getFlu_old(row []string) (flu models.FeedLineUnit, err error) {
	// TODO Refactor the code below
	// CSV Headers should be :
	// "Id", "ReferenceId", "Data", "Build", "Tag", "ProjectId", "StepId", "CreatedAt"
	id, err := uuid.FromString(row[0])
	if err != nil {
		plog.Error("Error ID:", err)
		return flu, errors.New("ID is not valid.")
	}
	referenceId := row[1]
	if referenceId == "" {
		plog.Error("Error Ref-ID:", errors.New("reference id cant be empty"))
		return flu, errors.New("reference id cant be empty")
	}
	// skip the data scanning part in future
	// user should not be able to change data
	data := models.JsonFake{}
	if err := data.Scan(row[2]); err != nil {
		plog.Error("Error Data:", err)
		return flu, errors.New("Data field is not valid.")
	}

	build := models.JsonFake{}
	if err := build.Scan(row[3]); err != nil {
		plog.Error("Error Build:", err)
		return flu, errors.New("Build field is not valid.")
	}

	tag := row[4]
	if tag == "" {
		plog.Error("Error Tag:", errors.New("tag cant be empty"))
		return flu, errors.New("Tag field is not valid")
	}

	projectId, err := uuid.FromString(row[5])
	if err != nil {
		plog.Error("Error ProjectID:", err)
		return flu, errors.New("Project ID is not valid")
	}

	stepId, err := uuid.FromString(row[6])
	if err != nil {
		plog.Error("Error StepID:", err)
		return flu, errors.New("StepID is not valid")
	}

	createdAtTime, err := time.Parse(timeFormat, row[7])
	if err != nil {
		plog.Error("Error Time:", errors.New("Time format not valid "+row[7]))
		//return errors.New("CreatedTime : Time format is not valid " + row[7])
	}
	createdAt := pq.NullTime{createdAtTime, true}

	updatedAt := pq.NullTime{time.Now(), true}

	flu = models.FeedLineUnit{
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
	return flu, nil
}
