package purge

import (
	"fmt"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type SchedulerCleaner interface {
	Purge()
	cleanIntervals()
	cleanIntervalActions()
}

type schedulerCleaner struct {
	baseUrl string
}

// NewSchedulerCleaner creates an instance of SchedulerCleaner
func NewSchedulerCleaner() SchedulerCleaner {
	fmt.Println("\n * Scheduler")
	return &schedulerCleaner{
		baseUrl: config.Conf.Clients["Scheduler"].Url(),
	}
}

func (d *schedulerCleaner) Purge() {
	d.cleanIntervals()
	d.cleanIntervalActions()
}

func (d *schedulerCleaner) cleanIntervals(){
	url := d.baseUrl + clients.ApiIntervalRoute
	var intervals []models.Interval
	err := request.Get(url, &intervals)
	if err != nil {
		return
	}

	var count int
	for _, interval := range intervals {
		err = request.Delete(url + "/" + interval.ID)
		if err == nil {
			count = count +1
		}
	}
	fmt.Printf("Removed %d Intervals from %d \n", count, len(intervals))
}

func (d *schedulerCleaner) cleanIntervalActions() {
	url := d.baseUrl + clients.ApiIntervalActionRoute
	var intervalActions []models.IntervalAction
	err := request.Get(url, &intervalActions)
	if err != nil {
		return
	}

	var count int
	for _, intervalAction := range intervalActions {
		err = request.Delete(url + "/" + intervalAction.ID)
		if err == nil {
			count = count +1
		}
	}
	fmt.Printf("Removed %d Interval Actions from %d \n", count, len(intervalActions))
}
