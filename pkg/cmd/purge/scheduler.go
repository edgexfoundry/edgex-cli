package purge

import (
	"context"
	"fmt"

	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/v2/models"
)

type schedulerCleaner struct {
	baseUrl string
	ctx     context.Context
}

// NewSchedulerCleaner creates an instance of SchedulerCleaner
func NewSchedulerCleaner(ctx context.Context) Purgeable {
	fmt.Println("\n * Scheduler")
	return &schedulerCleaner{
		baseUrl: config.Conf.Clients["Scheduler"].Url(),
		ctx:     ctx,
	}
}

func (d *schedulerCleaner) Purge() {
	d.cleanIntervals()
	d.cleanIntervalActions()
}

func (d *schedulerCleaner) cleanIntervals() {
	url := d.baseUrl + clients.ApiIntervalRoute
	var intervals []models.Interval
	err := request.Get(d.ctx, url, &intervals)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	var count int
	for _, interval := range intervals {
		err = request.Delete(d.ctx, url+"/"+interval.Id)
		if err == nil {
			count = count + 1
		}
	}
	fmt.Printf("Removed %d Intervals from %d \n", count, len(intervals))
}

func (d *schedulerCleaner) cleanIntervalActions() {
	url := d.baseUrl + clients.ApiIntervalActionRoute
	var intervalActions []models.IntervalAction
	err := request.Get(d.ctx, url, &intervalActions)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	var count int
	for _, intervalAction := range intervalActions {
		err = request.Delete(d.ctx, url+"/"+intervalAction.Id)
		if err == nil {
			count = count + 1
		}
	}
	fmt.Printf("Removed %d Interval Actions from %d \n", count, len(intervalActions))
}
