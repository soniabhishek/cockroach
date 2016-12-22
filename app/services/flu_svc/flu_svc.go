package flu_svc

import (
	"encoding/csv"
	"errors"
	"github.com/crowdflux/angel/app/DAL/imdb"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/flu_upload_status"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_errors"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"github.com/crowdflux/angel/app/services/plerrors"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc"
	"github.com/crowdflux/angel/utilities"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

type fluService struct {
	fluRepo      feed_line_repo.IFluRepo
	fluValidator flu_validator.IFluValidatorService
	projectsRepo projects_repo.IProjectsRepo
	workFlowSvc  work_flow_executor_svc.IWorkFlowSvc
}

var _ IFluService = &fluService{}

func (i *fluService) AddFeedLineUnit(flu *models.FeedLineUnit) error {

	flu.Build = flu.Data.Copy()
	_, err := i.fluValidator.Validate(flu)
	if err != nil {
		return err
	}
	err = i.CheckProjectExists(flu.ProjectId)
	if err != nil {
		return err
	}

	fin := feed_line_repo.NewInputQueue()
	id, err := fin.Add(*flu)
	flu.ID = id
	if err != nil && err == feed_line_repo.ErrDuplicateReferenceId {
		err = flu_errors.ErrDuplicateReferenceId
	}

	return err
}

func (i *fluService) SyncInputFeedLine() error {

	fluInputQueue := feed_line_repo.NewInputQueue()

	flus, err := fluInputQueue.GetQueued()

	if err != nil {

		plog.Error("Error occured while getting data", err)
		return err
	}

	if len(flus) > 0 {

		for i, _ := range flus {

			flus[i].MasterId = flus[i].ID
			flus[i].IsActive = true
			flus[i].IsMaster = true
		}

		err = i.fluRepo.BulkInsert(flus)

		if err != nil {
			plog.Error("Bulk insert failed", err)
			return err
		}

		// start adding to workFlowSvc in another go routine
		go func() {

			for _, flu := range flus {

				i.workFlowSvc.AddFLU(flu)
			}
		}()

		err = fluInputQueue.MarkFinished(flus)

		if err != nil {
			plog.Error("Changing queue status failed", err)
			return err
		}
		//plog.Info(len(flus), "flus processed")

	}

	return nil
}

func (i *fluService) GetFeedLineUnit(fluId uuid.UUID) (models.FeedLineUnit, error) {

	flu, err := i.fluRepo.GetById(fluId)
	if err != nil && err == feed_line_repo.ErrFLUNotFoundInInputQueue {
		err = flu_errors.ErrFluNotFound
	}
	return flu, err
}

//This Services Read From File Descriptor and fetch contents row by row
func (i *fluService) BulkAddFeedLineUnit(file multipart.File, errorCsv *os.File, fileName string, projectId uuid.UUID) {
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3 // so the reader will always check how many records are present in each row sequence ['Reference Id', 'Tag', 'Body']

	errorWriter := csv.NewWriter(errorCsv)

	batcherChan := make(chan models.FeedLineUnit) //this channel will be used to batch the Flus
	errorChan := make(chan plerrors.ChildError)   //this channel will be used to receive errors
	errorWriterChan := make(chan []string)        //this channel is for writing data to error csv

	//check header for utf-8 and csv format
	row, err := reader.Read()
	if err != nil {
		//Put Failed Status
		plog.Error("Not in correct format.", errors.New("Invalid Upload"))
		i.UpdateUploadStatus(projectId, flu_upload_status.Failure)
		return
	}
	wrongCol, err := utilities.IsValidUTF8(row)
	if wrongCol != -1 {
		//Put Failed Status
		plog.Error("Not in correct encoding[UTF-8].", errors.New("Invalid Upload"))
		i.UpdateUploadStatus(projectId, flu_upload_status.Failure)
		return
	}

	go writeCsvError(errorWriter, errorWriterChan, projectId.String())
	//This will be Used to collect flus from batcherChann and Create a batch of given size and returns bulk error in errorChann
	go mongoBatcher(batcherChan, errorChan, 3000)

	defer close(batcherChan)

	//This Go routine will be used to fetch errors in bulk insert and will write them in error file
	go func() {
		receiveBulkError(errorChan, errorWriterChan)
		close(errorWriterChan)
		errorWriter.Flush()
		errorCsv.Close()
	}()

	i.UpdateUploadStatus(projectId, flu_upload_status.Processing)
	var cnt int = 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		cnt++
		//Updating cache Status

		project := projectId.String()
		fus, _ := imdb.FluUploadCache.Get(project)
		fus.CompletedFluCount = cnt
		imdb.FluUploadCache.Set(project, fus)

		if err != nil {
			plog.Error(" csv reading error", err)
			row = append(row, err.Error())
			errorWriterChan <- row
			continue
		}

		wrongCol, err := utilities.IsValidUTF8(row)
		if wrongCol != -1 {
			plog.Error("Not in correct encoding[UTF-8]. [Row:"+strconv.Itoa(cnt)+", Col:"+strconv.Itoa(wrongCol)+"]", err)
			row = append(row, err.Error())
			errorWriterChan <- row
			continue
		}

		flu, err := getFlu(row, projectId)
		if err != nil {
			plog.Error(" csv reading error", err)
			row = append(row, err.Error())
			errorWriterChan <- row
			continue
		}
		isValid, err := i.fluValidator.Validate(&flu)
		if err != nil {
			row = append(row, err.Error())
			errorWriterChan <- row
			continue
		}
		if !isValid {
			row = append(row, err.Error())
			errorWriterChan <- row
			continue
		}
		batcherChan <- flu
	}
}

//--------------------------------------------------------------------------------//
//CHECK PROJECT
//--------------------------------------------------------------------------------//

func (i *fluService) CheckProjectExists(mId uuid.UUID) error {
	_, err := i.projectsRepo.GetById(mId)
	return err
}

func (i *fluService) CheckCsvUploaded(projectId string) (bool, error) {
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
func (i *fluService) UpdateUploadStatus(projectId uuid.UUID, status flu_upload_status.FluUploadStatus) {
	val := models.FluUploadStats{}
	val.Status = status
	imdb.FluUploadCache.Set(projectId.String(), val)
}

func (i *fluService) GetUploadStatus(projectId string) (models.FluUploadStats, error) {
	return imdb.FluUploadCache.Get(projectId)
}
