package cronscheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

var (
	CronScheduler *gocron.Scheduler
)

func InitScheduler() {
	CronScheduler = gocron.NewScheduler(time.UTC)
	CronScheduler.TagsUnique()
	CronScheduler.StartAsync()
}
