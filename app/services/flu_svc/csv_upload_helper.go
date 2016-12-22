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
	"github.com/crowdflux/angel/app/services/plerrors"
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

func mongoBatcher(c chan models.FeedLineUnit, err chan plerrors.ChildError, batchSize int) {
	fin := feed_line_repo.NewInputQueue()
	fluBatch := make([]models.FeedLineUnit, 0, batchSize)
	defer close(err)

	count := 0
	for {
		flu, ok := <-c
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

func writeCsvError(csvWrite *csv.Writer, c chan []string, projectId string) {
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
	}
	imdb.FluUploadCache.Set(projectId, fus)
}

func receiveBulkError(errChannel chan plerrors.ChildError, c chan []string) {
	for {
		childErr, ok := <-errChannel
		if !ok {
			break
		} else {
			c <- []string{childErr.Flu.ReferenceId, childErr.Flu.Tag, childErr.Flu.Build.String(), childErr.Message}
		}
	}
}
