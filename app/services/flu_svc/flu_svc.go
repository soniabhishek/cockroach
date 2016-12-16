package flu_svc

import (
	"encoding/csv"
	"errors"
	"github.com/crowdflux/angel/app/DAL/repositories/feed_line_repo"
	"github.com/crowdflux/angel/app/DAL/repositories/projects_repo"
	"github.com/crowdflux/angel/app/models"
	"github.com/crowdflux/angel/app/models/uuid"
	"github.com/crowdflux/angel/app/plog"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_errors"
	"github.com/crowdflux/angel/app/services/flu_svc/flu_validator"
	"github.com/crowdflux/angel/app/services/work_flow_executor_svc"
	"github.com/crowdflux/angel/utilities"
	"github.com/lib/pq"
	"io"
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

func (i *fluService) AddFeedLineUnit(flu *models.FeedLineUnit) error {

	flu.Build = flu.Data.Copy()
	_, err := i.fluValidator.Validate(flu)
	if err != nil {
		return err
	}

	err = checkProjectExists(i.projectsRepo, flu.ProjectId)
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

func (i *fluService) BulkAddFeedLineUnit(fileName string) {
	file := `./uploads/` + string(os.PathSeparator) + fileName
	csvFile, err := os.Open(file)
	if err != nil {
		plog.Error("Manual Step", err, "csv opening")
		mongoTracker()
		return
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	reader.FieldsPerRecord = 2 // so the reader will always check how many records are present in each row

	batcherChan := make(chan models.FeedLineUnit)
	errorChan := make(chan error)

	go mongoBatcher(batcherChan, 1000, errorChan)

	var cnt int = -1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			close(batcherChan)
			break
		}
		if err != nil {
			plog.Error(" csv reading error", err)
			mongoTracker()
			return
		}
		cnt++

		wrongCol, err := utilities.IsValidUTF8(row)
		if wrongCol != -1 {
			plog.Error("Not in correct encoding[UTF-8]. [Row:"+strconv.Itoa(cnt)+", Col:"+strconv.Itoa(wrongCol)+"]", err)
			mongoTracker()
			return
		}
		if cnt == 0 {
			continue
		}

		flu, err := getFlu(row)
		if err != nil {
			plog.Error(" csv reading error", err)
			mongoTracker()
			return
		}
		flu.Build = flu.Data.Copy()
		_, err = i.fluValidator.Validate(flu)
		if err != nil {
			mongoTracker()
			return
		}

		err = checkProjectExists(i.projectsRepo, flu.ProjectId)
		if err != nil {
			mongoTracker()
			return
		}

		batcherChan <- flu
	}
}

//--------------------------------------------------------------------------------//
//CHECK PROJECT
//--------------------------------------------------------------------------------//

func checkProjectExists(r projects_repo.IProjectsRepo, mId uuid.UUID) error {
	_, err := r.GetById(mId)
	return err
}

//This Should be implemented for testing purspose
func mongoTracker() {

}

func getFlu(row []string) (flu models.FeedLineUnit, err error) {
	fluId := row[0]
	id, err := uuid.FromString(fluId)
	if err != nil {
		plog.Error("Error ID:", err)
		return flu, errors.New("ID is not valid. [" + fluId + "]")
	}

	build := models.JsonF{}
	buildVal := row[1]
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

type feedLineInputModel struct {
	models.FeedLineUnit `bson:",inline"`
	RetryCount          uint                `bson:"retry_count"`
	Status              feedLineQueueStatus `bson:"status"`
	IdString            string              `bson:"id_string"`
	ProjectIdString     string              `bson:"project_id_string"`
}

type feedLineQueueStatus uint

func mongoBatcher(c chan models.FeedLineUnit, batchSize int, err chan error) {
	fin := feed_line_repo.NewInputQueue()

	bulkData := make([]interface{}, 0, batchSize)
	count := 0
	for {
		flu, ok := <-c
		if !ok {
			if len(bulkData != 0) {
				fin.BulkAdd(bulkData)
			}
			break
		} else {
			switch {
			case flu.ID == uuid.Nil:
				flu.ID = uuid.NewV4()
				fallthrough
			case flu.CreatedAt.Valid == false:
				flu.CreatedAt = pq.NullTime{time.Now(), true}
			}
			flu.IsMaster = true
			flu.IsActive = true
			flu.MasterId = flu.ID
			bulkData = append(bulkData, feedLineInputModel{
				FeedLineUnit:    flu,
				RetryCount:      0,
				Status:          0,
				IdString:        flu.ID.String(),
				ProjectIdString: flu.ProjectId.String(),
			})
			count++
			if count%batchSize == 0 {
				fin.BulkAdd(bulkData)
				bulkData = bulkData[:0]
			}
		}
	}
}
