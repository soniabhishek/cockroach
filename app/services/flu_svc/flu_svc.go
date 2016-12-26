package flu_svc

import (
	"encoding/csv"
	"errors"
	"fmt"
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
	"time"
	"sync"
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

//--------------------------------------------------------------------------------//
//CHECK PROJECT
//--------------------------------------------------------------------------------//

func (i *fluService) CsvCheckBasicValidation(file multipart.File, fileName string, projectId uuid.UUID) error {

	valid, err := checkCsvUploaded(projectId.String())
	if !valid {
		plog.Error("Already Exist", err, projectId.String())
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
			plog.Error(" csv reading error", err)
			return err
		}

		wrongCol, err := utilities.IsValidUTF8(row)
		if wrongCol != -1 {
			plog.Error("Not in correct encoding[UTF-8]. [Row:"+strconv.Itoa(cnt)+", Col:"+strconv.Itoa(wrongCol)+"]", err)
			return err
		}

		if cnt == 0 {
			continue
		}

		if _, ok := referenceIdMapper[row[0]]; ok {
			err = errors.New("duplicate Reference Id uploaded")
			plog.Error("Reference ID Duplicate", err)
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
	setUploadStatus(projectId, fls)

	go i.startRowsProcessing(filePath, fileName, projectId)
	return nil
}

func (fs *fluService) startRowsProcessing(filePath string, fileName string, projectId uuid.UUID) {
	readerFile, err := os.Open(filePath)
	if err != nil {
		panic(err)
		return
	}
	fileReader := csv.NewReader(readerFile)

	errFilePath := fmt.Sprintf("./uploads/error_%s_%s.csv", strconv.Itoa(int(time.Now().UnixNano())), projectId.String() )

	errorCsv, err := os.Create(errFilePath)
	if err != nil {
		panic(err)
		return
	}
	errorWriter := csv.NewWriter(errorCsv)

	batcherChan := make(chan models.FeedLineUnit) //this channel will be used to batch the Flus
	errorChan := make(chan plerrors.ChildError)   //this channel will be used to receive errors
	errorWriterChan := make(chan []string)        //this channel is for writing data to error csv
	validatorChan := make(chan models.FeedLineUnit)

	go writeCsvError(errorWriter, errorWriterChan, projectId.String())
	//This will be Used to collect flus from batcherChann and Create a batch of given size and returns bulk error in errorChann
	go mongoBatcher(batcherChan, errorChan, 3000)

	defer func() {
		readerFile.Close()
		err := os.Remove(filePath)
		if err != nil {
			panic(err)
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < 50 ; i ++ {
		go func(valChan chan models.FeedLineUnit) {
			wg.Add(1)
			for {
				flu, ok := <-valChan
				if ok {
					isValid, err := fs.fluValidator.Validate(&flu)
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

		project := projectId.String()
		fus, _ := fs.GetUploadStatus(project)
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

func (i *fluService) CheckProjectExists(mId uuid.UUID) error {
	_, err := i.projectsRepo.GetById(mId)
	return err
}

func (i *fluService) GetUploadStatus(projectId string) (models.FluUploadStats, error) {
	return imdb.FluUploadCache.Get(projectId)
}
