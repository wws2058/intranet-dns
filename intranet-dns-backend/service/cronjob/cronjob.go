package cronjob

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/tswcbyy1107/intranet-dns/config"
	"github.com/tswcbyy1107/intranet-dns/database"
	"github.com/tswcbyy1107/intranet-dns/models"
	"github.com/tswcbyy1107/intranet-dns/service/redis"
	"github.com/tswcbyy1107/intranet-dns/utils"
)

var (
	CronjobRefreshTs = "global_cronjob-register-ts"
	currentTs        int64
	cronManager      *cron.Cron
)

type RunJob struct {
	c *models.Cronjob
}

func (r *RunJob) Run() {
	if !r.c.Started {
		return
	}

	key := fmt.Sprintf("%s.cronjob.%s", config.GlobalConfig.App.Name, r.c.Name)
	if redis.Lock(key, time.Hour) {
		return
	}
	defer redis.UnlockNow(key)

	err := models.TemplateQuery(&models.DaoDBReq{
		Dst: r.c,
		ModelFilter: models.Cronjob{
			BaseModel: models.BaseModel{Id: r.c.Id},
		},
	})
	if err != nil {
		return
	}

	if r.c.History == nil {
		initHistory := models.TaskHistory{}
		r.c.History = initHistory
	}

	uid := utils.GenUUID()
	callAt := time.Now()
	var outerErr error
	defer func() {
		record := models.TaskRecord{
			UID:     uid,
			CallAt:  models.JsonTime(callAt),
			Succeed: true,
		}
		r.c.LastSucceed = true
		if outerErr != nil {
			record.Succeed = false
			record.Error = outerErr.Error()
			r.c.LastSucceed = false
			logrus.WithField("uid", uid).WithField("name", r.c.Name).Errorf("cronjob called failed: %v", err)
		} else {
			logrus.WithField("uid", uid).WithField("name", r.c.Name).Info("cronjob called succeed")
		}
		r.c.History.Add(record)
		models.TemplateCreate(r.c)
	}()

	switch r.c.TaskType {
	case models.HttpType:
		url := r.c.TaskArgs.Url
		resp, err := http.Get(url)
		if err != nil {
			outerErr = err
			return
		}
		body, err := io.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			outerErr = err
			return
		}
		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("result:%s %v", string(body), resp.StatusCode)
			outerErr = err
			return
		}
	case models.FuncType:
		f, ok := internalFunctionMaps[r.c.TaskArgs.FunctionName]
		if !ok {
			err := fmt.Errorf("no such function:%s", r.c.TaskArgs.FunctionName)
			outerErr = err
			return
		}
		err := f()
		if err != nil {
			outerErr = err
		}
	}
}

// background goroutine update cronjobs
func InitCronJob() {
	registerCronJob()
	currentTs = time.Now().Unix()
	go func() {
		for {
			time.Sleep(10 * time.Second)
			var cacheTs int64
			err := redis.LoadCache(CronjobRefreshTs, &cacheTs)
			if err != nil {
				return
			}
			if currentTs < cacheTs {
				registerCronJob()
				currentTs = cacheTs
			}
		}
	}()
}

// update cronjob register ts version
func RefreshRegisterJobTs() {
	redis.Cache(CronjobRefreshTs, time.Now().Unix(), time.Minute*5)
}

// db cronjobs -> cronjobs
func registerCronJob() {
	if cronManager != nil {
		cronManager.Stop()
	}

	cronManager = cron.New(cron.WithChain(cron.Recover(cron.DefaultLogger), cron.DelayIfStillRunning(cron.DefaultLogger)))
	cronjobs := []*models.Cronjob{}
	err := database.DB.Find(&cronjobs).Error
	if err != nil {
		os.Exit(0)
		return
	}

	hostname, _ := os.Hostname()
	for _, c := range cronjobs {
		if !c.Started {
			continue
		}
		entryID, err := cronManager.AddJob(c.Spec, &RunJob{c: c})
		if err != nil {
			os.Exit(0)
		}
		logrus.Infof("cronjob init: %s %s %s %v", hostname, c.Name, c.Spec, entryID)
	}
	cronManager.Start()
}
