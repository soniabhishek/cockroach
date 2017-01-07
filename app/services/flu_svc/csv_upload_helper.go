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
	"github.com/crowdflux/angel/app/services/flu_svc/flu_errors"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"os"
	"fmt"
	"strconv"
	"time"
	"sync"
	"io"
)

//This will Initialize and Return a Flu from a csv row
func getFlu(row []string, projectId uuid.UUID) (flu models.FeedLineUnit, err error) {
	fluRefId := row[0] //first column will contain the Reference ID
	if len(fluRefId) < 1 {
		return flu, errors.New("ID is not valid. [" + fluRefId + "]")
	}

	data := models.JsonF{}
	buildVal := row[2]
	//Unmarshalling
	if err := data.Scan(buildVal); err != nil {
		return flu, errors.New("Build field is not valid. [" + buildVal + "]")
	}
	flu = models.FeedLineUnit{
		ReferenceId: fluRefId,
		ProjectId:   projectId,
		Data:        data,
		Build:       data.Copy(),
		Tag:         row[1],
	}
	return flu, nil
}

func mongoBatcher(fluChan chan models.FeedLineUnit, err chan feed_line_repo.BulkError, batchSize int) {
	fin := feed_line_repo.NewInputQueue()
	fluBatch := make([]models.FeedLineUnit, 0, batchSize)
	defer close(err)

	count := 0
	for {
		flu, ok := <-fluChan
		if ok {

			fluBatch = append(fluBatch, flu)
			count++
			if count == batchSize {
				blkerror := fin.BulkAdd(fluBatch)
				if blkerror.Error != nil {
					if blkerror.Error == flu_errors.ErrBulkError {
						for _, x := range blkerror.BulkError {
							err <- x
						}
					} else {
						plog.Error("BULK_ERROR", blkerror.Error)
					}
				}
				fluBatch = fluBatch[:0]
				count = 0
			}
		} else {
			if len(fluBatch) != 0 {
				blkerror := fin.BulkAdd(fluBatch)
				if blkerror.Error != nil {
					if blkerror.Error == flu_errors.ErrBulkError {
						for _, x := range blkerror.BulkError {
							err <- x
						}
					} else {
						plog.Error("BULK_ERROR", blkerror.Error)
					}
				}
			}
			break
		}
	}
}

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

func startRowsProcessing(fv flu_validator.IFluValidatorService, filePath string, fileName string, projectId uuid.UUID) {
	//This will Open file from disk and we are defering close and remove call to this file
	readerFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer func() {
		readerFile.Close()
		err := os.Remove(filePath)
		if err != nil {
			panic(err)
		}
	}()
	fileReader := csv.NewReader(readerFile)

	//Error file creation
	errFilePath := fmt.Sprintf("./uploads/error_%s_%s.csv", strconv.Itoa(int(time.Now().UnixNano())), projectId.String() )
	errorCsv, err := os.Create(errFilePath)
	if err != nil {
		panic(err)
	}
	errorWriter := csv.NewWriter(errorCsv)


	/**
		Channels to communication between goroutines
	 */
	batcherChan := make(chan models.FeedLineUnit) //this channel will be used to batch the Flus
	errorChan := make(chan feed_line_repo.BulkError)   //this channel will be used to receive errors
	errorWriterChan := make(chan []string)        //this channel is for writing data to error csv
	validatorChan := make(chan models.FeedLineUnit)


	go writeCsvError(errorWriter, errorWriterChan, projectId.String(), errFilePath)
	//This will be Used to collect flus from batcherChann and Create a batch of given size and returns bulk error in errorChann
	go mongoBatcher(batcherChan, errorChan, 3000)


	//starting concurrent validation request
	var wg sync.WaitGroup
	for i := 0; i < 50 ; i ++ {
		go func(valChan chan models.FeedLineUnit) {
			wg.Add(1)
			for {
				flu, ok := <-valChan
				if ok {
					isValid, err := fv.Validate(&flu)
					if err != nil {
						errorWriterChan <- []string{flu.ReferenceId, flu.Tag, flu.Build.String(), err.Error()}
						continue
					}
					if !isValid {
						errorWriterChan <- []string{flu.ReferenceId, flu.Tag, flu.Build.String(), err.Error()}
						continue
					}
					batcherChan <- flu
				} else {
					break
				}
			}
			wg.Done()
		}(validatorChan)
	}

	//This Go routine will be used to fetch errors in bulk insert and will write them in error file
	go func() {
		receiveBulkError(errorChan, errorWriterChan)
		close(errorWriterChan)
		errorWriter.Flush()
		errorCsv.Close()
		fmt.Println("final", time.Now())
	}()

	cnt := 0
	for {
		row, err := fileReader.Read()
		if err == io.EOF {
			break
		}

		cnt++

		//Updating cache Status
		fus, _ := imdb.FluUploadCache.Get(projectId.String())
		fus.CompletedFluCount = cnt
		setUploadStatus(projectId, fus)

		flu, err := getFlu(row, projectId)
		if err != nil {
			plog.Error(" csv reading error", err)
			row = append(row, err.Error())
			errorWriterChan <- row
			continue
		}
		validatorChan <- flu
	}
	close(validatorChan)
	wg.Wait()
	close(batcherChan)
}

func checkCsvUploaded(projectId string) (bool, error) {
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

