package manual_step

import (
	"encoding/csv"
	"errors"
	"gitlab.com/playment-main/angel/app/DAL/repositories/feed_line_repo"
	"gitlab.com/playment-main/angel/app/config"
	"gitlab.com/playment-main/angel/app/models"
	"gitlab.com/playment-main/angel/app/models/uuid"
	"gitlab.com/playment-main/angel/app/plog"
	"gitlab.com/playment-main/angel/utilities/constants"
	"io"
	"os"
	"time"
)

const timeFormat = time.RFC3339

type megatronJson struct {
	Jsons []models.JsonFake `json:jsons`
}

func DownloadCsv(manualStepId uuid.UUID) (string, error) {
	flRepo := feed_line_repo.New()
	flus, err := flRepo.GetByStepId(manualStepId)
	if err != nil {
		plog.Error("Manual Step", err, manualStepId)
		return constants.Empty, err
	}
	plog.Info("manual step flus going to be downloaded", flus, manualStepId)

	path := config.Get(config.DOWNLOAD_PATH)
	//file, err := createCSV(flus, path, manualStepId)

	file, numOfLines, err := createJSONFile(flus, path, manualStepId)
	if err != nil {
		plog.Error("Write file error", err, manualStepId)
		return constants.Empty, err
	}

	if numOfLines == 0 {
		return constants.Empty, errors.New("No Data to show.")
	}

	url := config.Get(config.MEGATRON_API)
	filename, err := FlattenCSV(file, url, manualStepId)
	if err != nil {
		plog.Error("Transformation error", err, manualStepId)
		return constants.Empty, errors.New("Transformation Error [" + err.Error() + "]")
	}
	return url + filename, nil
}

func createJSONFile(flus []models.FeedLineUnit, path string, manualStepId uuid.UUID) (file string, numOfLines int, err error) {

	file = path + string(os.PathSeparator) + manualStepId.String() + ".txt"
	err = createFile(file)
	if err != nil {
		plog.Error("Create file error", err, manualStepId)
		return constants.Empty, 0, nil
	}

	csvBuff := megatronJson{make([]models.JsonFake, 0)}
	for _, obj := range flus {
		var jsMap models.JsonFake = make(map[string]interface{})
		jsMap[ID] = obj.ID.String()
		jsMap[REF_ID] = obj.ReferenceId
		jsMap[DATA] = obj.Data.String()

		jsMap[BUILD] = obj.Build.String()
		jsMap[TAG] = obj.Tag
		jsMap[PROJECT_ID] = obj.ProjectId.String()

		jsMap[STEP_ID] = obj.StepId.String()

		if obj.CreatedAt.Valid {
			jsMap[CREATED_ID] = obj.CreatedAt.Time.Format(timeFormat)
		} else {
			jsMap[CREATED_ID] = constants.Empty
		}

		if obj.UpdatedAt.Valid {
			jsMap[UPDATED_AT] = obj.UpdatedAt.Time.Format(timeFormat)
		} else {
			jsMap[UPDATED_AT] = constants.Empty
		}

		csvBuff.Jsons = append(csvBuff.Jsons, jsMap)
	}

	// Write unmarshaled json data to CSV file
	err = writeFile(file, csvBuff)
	return file, len(csvBuff.Jsons), err

}

func createCSV(flus []models.FeedLineUnit, path string, manualStepId uuid.UUID) (file string, err error) {

	file = path + string(os.PathSeparator) + manualStepId.String() + ".csv"
	err = createFile(file)
	if err != nil {
		plog.Error("Create file error", err, manualStepId)
		return constants.Empty, nil
	}
	csvBuff := [][]string{{ID, REF_ID, DATA, BUILD, TAG, PROJECT_ID, STEP_ID, CREATED_ID, UPDATED_AT}}

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
			record = append(record, constants.Empty)
		}

		if obj.UpdatedAt.Valid {
			record = append(record, obj.UpdatedAt.Time.Format(timeFormat))
		} else {
			record = append(record, constants.Empty)
		}

		csvBuff = append(csvBuff, record)
	}

	// Write unmarshaled json data to CSV file
	err = writeCSV(file, csvBuff)
	return file, err
}

func UploadCsv(filename string) error {
	file := TEMP_FOLDER + string(os.PathSeparator) + filename
	csvFile, err := os.Open(file)
	if err != nil {
		plog.Error("Manual Step", err, "csv opening")
		return err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
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
	err = flRepo.BulkFluBuildUpdate(flus)
	if err != nil {
		plog.Info(err.Error())
	}
	return err
}

func getFlu(row []string) (flu models.FeedLineUnit, err error) {
	fluId := row[FLU_ID_INDEX]
	id, err := uuid.FromString(fluId)
	if err != nil {
		plog.Error("Error ID:", err)
		return flu, errors.New("ID is not valid. [" + fluId + "]")
	}

	build := models.JsonFake{}
	buildVal := row[BUILD_INDEX]
	if err := build.Scan(buildVal); err != nil {
		plog.Error("Error Build:", err)
		return flu, errors.New("Build field is not valid. [" + buildVal + "]")
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

func writeFile(filepath string, records megatronJson) error {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(filepath, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(`{"jsons" : [ `)
	if err != nil {
		return err
	}

	l := len(records.Jsons)
	c := 0
	// write some text to file
	for _, mj := range records.Jsons {
		c++
		stringToWrite := ""
		if c < l {
			stringToWrite = mj.String() + COMMA
		} else {
			stringToWrite = mj.String()
		}
		_, err = file.WriteString(stringToWrite)
		if err != nil {
			return err
		}
	}
	_, err = file.WriteString(`]}`)
	if err != nil {
		return err
	}

	// save changes
	err = file.Sync()
	return err
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
