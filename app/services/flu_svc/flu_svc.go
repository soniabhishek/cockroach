package flu_svc

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/crowdflux/angel/app/DAL/imdb"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/config"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/flu_upload_status"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_errors"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc"
	"github.com/crowdflux/angel/utilities"
	"io"
	"mime/multipart"
	"os"
	"strconv"
	"time"
)

type fluService struct {
	fluRepo      feed_line_repo.IFluRepo
	fluValidator flu_validator.IFluValidatorService
	projectsRepo projects_repo.IProjectsRepo
	workFlowSvc  work_flow_executor_svc.IWorkFlowSvc
}

var _ IFluService = &fluService{}

var timeout_sec = services.AtoiOrDefault(config.FEEDLINE_API_TIMEOUT_SEC.Get(), 10)

func (i *fluService) AddFeedLineUnit(flu *models.FeedLineUnit) error {

	timedOut := time.After(time.Duration(timeout_sec) * time.Second)

	errChan := make(chan error, 1)

	go func() {

		flu.Build = flu.Data.Copy()
		_, err := i.fluValidator.Validate(flu)
		if err != nil {
			errChan <- err
			return
		}

		err = checkProjectExists(i.projectsRepo, flu.ProjectId)
		if err != nil {
			errChan <- err
			return
		}

		errChan <- nil

	}()

	select {
	case <-timedOut:
		return flu_errors.ErrRequestTimedOut
	case err := <-errChan:
		if err != nil {
			return err
		}
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

		plog.Error("Flu_svc", err, plog.Message("Error occured while getting data"))
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
			plog.Error("Flu_svc", err, plog.Message("Bulk insert failed"))
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
			plog.Error("Flu_svc", err, plog.Message("Changing queue status failed"))
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

func (i *fluService) CsvCheckBasicValidation(file multipart.File, fileName string, projectId uuid.UUID) error {

	allowed, err := allowCsvUpload(projectId.String())
	if !allowed {
		return err
	}
	err = checkProjectExists(i.projectsRepo, projectId)
	if err != nil {
		return err
	}

	updateUploadStatus(projectId, flu_upload_status.Pending)
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = 3 // so the reader will always check how many records are present in each row sequence ['Reference Id', 'Tag', 'Body']

	filePath := fmt.Sprintf("./uploads/%s_%s.csv", strconv.Itoa(int(time.Now().UnixNano())), projectId.String())
	uploadFile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	defer uploadFile.Close()

	uploadWriter := csv.NewWriter(uploadFile)

	referenceIdMapper := make(map[string]struct{})

	var cnt int = -1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		cnt++

		if err != nil {
			plog.Error("Flu_svc", err, plog.M(" csv reading error"))
			return err
		}

		wrongCol, err := utilities.IsValidUTF8(row)
		if wrongCol != -1 {
			return err
		}

		//to skip the header row
		if cnt == 0 {
			continue
		}

		if _, ok := referenceIdMapper[row[0]]; ok {
			err = errors.New("duplicate Reference Id uploaded")
			return err
		} else {
			referenceIdMapper[row[0]] = struct{}{}
		}

		uploadWriter.Write(row)
	}
	uploadWriter.Flush()

	fls, err := i.GetUploadStatus(projectId.String())
	if err != nil {
		panic(err)
	}
	fls.TotalFluCount = cnt
	fls.Status = flu_upload_status.Processing
	setUploadStatus(projectId, fls)

	go processCSV(i.fluValidator, filePath, fileName, projectId)
	return nil
}

func (i *fluService) GetUploadStatus(projectId string) (models.FluUploadStats, error) {
	return imdb.FluUploadCache.Get(projectId)
}

//--------------------------------------------------------------------------------//
//CHECK PROJECT
//--------------------------------------------------------------------------------//

func checkProjectExists(r projects_repo.IProjectsRepo, mId uuid.UUID) error {
	_, err := r.GetById(mId)
	return err
}
