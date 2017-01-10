package flu_svc

import (
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_errors"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/models"
	"errors"
	"os"
	"encoding/csv"
	"sync"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"github.com/crowdflux/angel/app/DAL/imdb"
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

func getFileAndReader(filePath string) (file *os.File, fileReader *csv.Reader, err error) {
	file, err = os.Open(filePath)
	if err != nil {
		return
	}
	fileReader = csv.NewReader(file)
	return
}

func generateFileAndWriter(filePath string) (file *os.File, fileWriter *csv.Writer, err error) {
	file, err = os.Create(filePath)
	if err != nil {
		return
	}
	fileWriter = csv.NewWriter(file)
	return
}

func validateFluInBatch(fv flu_validator.IFluValidatorService, errorWriterChan chan []string, batcherChan, validatorChan chan models.FeedLineUnit){
	var wg sync.WaitGroup
	for i := 0; i < 50 ; i ++ {
		wg.Add(1)
		go func(valChan chan models.FeedLineUnit) {
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
	wg.Wait()
	close(batcherChan)
}

func processFluFromCSV(fileReader *csv.Reader, projectId uuid.UUID, errorWriterChan chan []string, validatorChan chan models.FeedLineUnit){
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
			row = append(row, err.Error())
			errorWriterChan <- row
			continue
		}
		validatorChan <- flu
	}
}