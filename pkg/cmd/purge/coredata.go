package purge

import (
	"context"
	"fmt"

	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients"
)

type coredataCleaner struct {
	baseUrl string
	ctx     context.Context
}

// NewCoredataCleaner creates an instance of CoreDataCleaner
func NewCoredataCleaner(ctx context.Context) Purgeable {
	fmt.Println("\n * core-data")
	return &coredataCleaner{
		baseUrl: config.Conf.Clients["CoreData"].Url(),
		ctx:     ctx,
	}
}
func (d *coredataCleaner) Purge() {
	d.cleanEventsAndReadings()
	//d.cleanValueDescriptors()
}

// TODO (jpw) - rid
// func (d *coredataCleaner) cleanValueDescriptors() {
// 	url := d.baseUrl + clients.ApiValueDescriptorRoute
// 	var valueDescriptors []models.ValueDescriptor
// 	err := request.Get(d.ctx, url, valueDescriptors)
// 	if err != nil {
// 		fmt.Printf("Error: %s\n", err.Error())
// 		return
// 	}

// 	var count int
// 	for _, valueDescriptor := range valueDescriptors {
// 		err = request.Delete(d.ctx, url+config.PathId+valueDescriptor.Id)
// 		if err == nil {
// 			count = count + 1
// 		}
// 	}
// 	fmt.Printf("Removed %d Value Descriptors from %d \n", count, len(valueDescriptors))
// }

func (d *coredataCleaner) cleanEventsAndReadings() {
	url := d.baseUrl + clients.ApiEventRoute + "/scruball"
	err := request.Delete(d.ctx, url)
	if err == nil {
		fmt.Print("All Events and Readings have been removed \n")
	}
}
