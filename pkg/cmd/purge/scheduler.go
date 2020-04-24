package purge

import (
	"fmt"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	client "github.com/edgexfoundry-holding/edgex-cli/pkg"
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
	err := client.ListHelper(url, &intervals)
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int
	for _, interval := range intervals {
		_, err = client.DeleteItem(url + "/" + interval.ID)
		if err != nil {
			fmt.Printf("Failed to delete Internals with id %s because of error: %s", interval.ID, err)
		} else {
			count = count +1
		}
	}
	fmt.Printf("Removed %d Intervals from %d \n", count, len(intervals))
}

func (d *schedulerCleaner) cleanIntervalActions() {
	url := d.baseUrl + clients.ApiIntervalActionRoute
	var intervalActions []models.IntervalAction
	err := client.ListHelper(url, &intervalActions)
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int
	for _, intervalAction := range intervalActions {
		_, err = client.DeleteItem(url + "/" + intervalAction.ID)

		if err != nil {
			fmt.Printf("Failed to delete Internal Actions with id %s because of error: %s", intervalAction.ID, err)
		} else {
			count = count +1
		}
	}
	fmt.Printf("Removed %d Interval Actions from %d \n", count, len(intervalActions))
}
