package purge

import (
	"fmt"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type coredataCleaner struct {
	baseUrl string
}

// NewCoredataCleaner creates an instance of CoreDataCleaner
func NewCoredataCleaner() Purgeable {
	fmt.Println("\n * core-data")
	return &coredataCleaner{
		baseUrl: config.Conf.Clients["CoreData"].Url(),
	}
}
func (d *coredataCleaner) Purge() {
	d.cleanEventsAndReadings()
	d.cleanValueDescriptors()
}

func (d *coredataCleaner) cleanValueDescriptors() {
	url := d.baseUrl + clients.ApiValueDescriptorRoute
	var valueDescriptors []models.ValueDescriptor
	err := request.Get(url, valueDescriptors)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	var count int
	for _, valueDescriptor := range valueDescriptors {
		err = request.Delete(url + config.PathId + valueDescriptor.Id)
		if err == nil {
			count = count + 1
		}
	}
	fmt.Printf("Removed %d Value Descriptors from %d \n", count, len(valueDescriptors))
}

func (d *coredataCleaner) cleanEventsAndReadings() {
	url := d.baseUrl + clients.ApiEventRoute + "/scruball"
	err := request.Delete(url)
	if err == nil {
		fmt.Print("All Events and Readings have been removed \n")
	}
}
