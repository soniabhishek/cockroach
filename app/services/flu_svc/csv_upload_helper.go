package flu_svc

import (
	"encoding/csv"
	"errors"
	"github.com/crowdflux/angel/app/DAL/imdb"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/flu_upload_status"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"os"
	"fmt"
	"strconv"
	"time"
)

func writeCsvError(csvWrite *csv.Writer, c chan []string, projectId, errFilePath string) {
	count := 0
	for {
		row, ok := <-c
		if !ok {
			break
		} else {
			count++

			fus, _ := imdb.FluUploadCache.Get(projectId)
			fus.ErrorFluCount = count
			imdb.FluUploadCache.Set(projectId, fus)

			if err := csvWrite.Write(row); err != nil {
				plog.Error("CSVERR001", err, "error while writning csv")
			}
		}
	}
	fus, _ := imdb.FluUploadCache.Get(projectId)
	if count != 0 {
		fus.Status = flu_upload_status.PartialUpload
	} else {
		fus.Status = flu_upload_status.Success
		err := os.Remove(errFilePath)
		if err != nil{
			panic(err)
		}
	}
	imdb.FluUploadCache.Set(projectId, fus)
}

func receiveBulkError(bulkErrChan chan feed_line_repo.BulkError, errRowChan chan []string) {
	for {
		childErr, ok := <-bulkErrChan
		if !ok {
			break
		} else {
			errRowChan <- []string{childErr.Flu.ReferenceId, childErr.Flu.Tag, childErr.Flu.Build.String(), childErr.Message}
		}
	}
}

func processCSV(fv flu_validator.IFluValidatorService, filePath string, fileName string, projectId uuid.UUID) {
	//This will Open file from disk and we are defering close and remove call to this file
	csvFile, fileReader, err := getFileAndReader(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		csvFile.Close()
		err := os.Remove(filePath)
		if err != nil {
			panic(err)
		}
	}()

	//Error file creation
	errFilePath := fmt.Sprintf("./uploads/error_%s_%s.csv", strconv.Itoa(int(time.Now().UnixNano())), projectId.String() )
	errorCsv, errorWriter, err := generateFileAndWriter(errFilePath)
	if err != nil {
		panic(err)
	}

	//Channels to communication between goroutines
	batcherChan := make(chan models.FeedLineUnit) //this channel will be used to batch the Flus
	errorChan := make(chan feed_line_repo.BulkError)   //this channel will be used to receive errors
	errorWriterChan := make(chan []string)        //this channel is for writing data to error csv
	validatorChan := make(chan models.FeedLineUnit)

	go writeCsvError(errorWriter, errorWriterChan, projectId.String(), errFilePath)
	//This will be Used to collect flus from batcherChann and Create a batch of given size and returns bulk error in errorChann
	go mongoBatcher(batcherChan, errorChan, 3000)
	//starting concurrent validation request this go-routine finishes on closing validatorChan
	go validateFluInBatch(fv, errorWriterChan, batcherChan, validatorChan)

	//This Go routine will be used to fetch errors in bulk insert and will write them in error file
	go func() {
		receiveBulkError(errorChan, errorWriterChan)
		close(errorWriterChan)
		errorWriter.Flush()
		errorCsv.Close()
		fmt.Println("final", time.Now())
	}()

	//this func will start fetching row from csv, will update upload status and then push to validator chan
	processFluFromCSV(fileReader,projectId, errorWriterChan, validatorChan)
	close(validatorChan)
}

func allowCsvUpload(projectId string) (bool, error) {
	fus, err := imdb.FluUploadCache.Get(projectId)
	if err != nil {
		return true, nil
	} else {
		if fus.Status == flu_upload_status.Failure || fus.Status == flu_upload_status.Success || fus.Status == flu_upload_status.PartialUpload {
			plog.Info("CHECHK CSV UPLOAD", "overriding existing status")
			return true, nil
		}
	}
	return false, errors.New("File Upload Alreay in Progress")
}

//This Will Override cached data. in case of success, failure, partial upload.
func updateUploadStatus(projectId uuid.UUID, status flu_upload_status.FluUploadStatus) {
	val := models.FluUploadStats{}
	val.Status = status
	imdb.FluUploadCache.Set(projectId.String(), val)
}

func setUploadStatus(projectId uuid.UUID, fus models.FluUploadStats) {
	imdb.FluUploadCache.Set(projectId.String(), fus)
}

