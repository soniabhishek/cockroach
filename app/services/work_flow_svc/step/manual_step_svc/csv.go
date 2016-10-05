package manual_step_svc

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/step_type"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/work_flow_svc/feed_line"
	"github.com/crowdflux/angel/utilities"
	"github.com/crowdflux/angel/utilities/constants"
)

const timeFormat = time.RFC3339

func DownloadCsv(manualStepId uuid.UUID) (string, error) {
	flRepo := feed_line_repo.New()
	flus, err := flRepo.GetByStepId(manualStepId)
	if err != nil {
		plog.Error("Manual Step", err, manualStepId)
		return constants.Empty, err
	}
	plog.Info("manual step flus going to be downloaded", len(flus), manualStepId)

	path := config.DOWNLOAD_PATH.Get()
	//file, err := createCSV(flus, path, manualStepId)

	file, numOfLines, err := createJSONFile(flus, path, manualStepId)
	if err != nil {
		plog.Error("Write file error", err, manualStepId)
		return constants.Empty, err
	}

	if numOfLines == 0 {
		return constants.Empty, errors.New("No Data to show.")
	}

	url := config.MEGATRON_API.Get() + "/flats"
	filename, err := FlattenCSV(file, url, manualStepId)
	if err != nil {
		plog.Error("Transformation error", err, manualStepId)
		return constants.Empty, errors.New("Transformation Error [" + err.Error() + "]")
	}
	return config.MEGATRON_API.Get() + filename, nil
}

func createJSONFile(flus []models.FeedLineUnit, path string, manualStepId uuid.UUID) (filePath string, numOfLines int, err error) {

	filePath = path + string(os.PathSeparator) + manualStepId.String() + ".txt"

	// creates a file , overwrites if exists
	file, err := os.Create(filePath)
	if err != nil {
		plog.Error("Create file error", err, manualStepId)
		return constants.Empty, 0, err
	}
	defer file.Close()

	type megatronJson struct {
		Jsons []models.FeedLineUnit `json:"jsons"`
	}

	bty, err := json.Marshal(megatronJson{flus})
	if err != nil {
		plog.Error("manual step", err, "Unable to create megatron json, manual step id : "+manualStepId.String())
		return
	}

	l, err := file.Write(bty)
	if err != nil {
		plog.Error("manual step", err, "error writing megatron json file for manual step id: "+manualStepId.String())
	}
	return filePath, l, err
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

		wrongCol, err := utilities.IsValidUTF8(row)
		if wrongCol != -1 {
			plog.Error("Manual Step", err, " csv is not in correct encoding[UTF-8]. [Row:"+strconv.Itoa(cnt)+", Col:"+strconv.Itoa(wrongCol)+"]")
			return err
		}

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
	plog.Info("Manual Step", "Flus going to be updated from csv upload ", len(flus), " first flu ", flus[0])

	flRepo := feed_line_repo.New()
	updatedFlus, err := flRepo.BulkFluBuildUpdateByStepType(flus, step_type.Manual)
	if err != nil {
		if err != feed_line_repo.ErrPartiallyUpdatedFlus {
			return err
		}
	}

	go func() {

		for _, flu := range updatedFlus {
			StdManualStep.finishFlu(feed_line.FLU{FeedLineUnit: flu})
		}
	}()

	// Returning err
	// will return nil if no error or ErrPartiallyUpdatedFlus if partially updated
	if err == feed_line_repo.ErrPartiallyUpdatedFlus {
		return errors.New("Partially updated flus. Count: " + strconv.Itoa(len(updatedFlus)))
	}
	return nil
}

func getFlu(row []string) (flu models.FeedLineUnit, err error) {
	fluId := row[FLU_ID_INDEX]
	id, err := uuid.FromString(fluId)
	if err != nil {
		plog.Error("Error ID:", err)
		return flu, errors.New("ID is not valid. [" + fluId + "]")
	}

	build := models.JsonF{}
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
