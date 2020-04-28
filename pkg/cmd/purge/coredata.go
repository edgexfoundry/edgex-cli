package purge

import (
	"fmt"
	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type CoreDataCleaner interface {
	Purge()
	cleanEvents()
	cleanReadings()
	cleanValueDescriptors()
}

type coredataCleaner struct {
	baseUrl string
}

// NewCoredataCleaner creates an instance of CoreDataCleaner
func NewCoredataCleaner() CoreDataCleaner {
	fmt.Println("\n * core-data")
	return &coredataCleaner{
		baseUrl: config.Conf.Clients["CoreData"].Url(),
	}
}
func (d *coredataCleaner) Purge() {
	d.cleanEvents()
	d.cleanReadings()
	d.cleanValueDescriptors()
}

func (d *coredataCleaner) cleanReadings() {
	url := d.baseUrl + clients.ApiReadingRoute
	var readings []models.Reading
	err := request.Get(url, readings)
	if err != nil {
		return
	}

	var count int
	for _, reading := range readings {
		// call delete function here
		err = request.Delete(url + config.PathId + reading.Id)
		if err == nil {
			count = count +1
		}
	}
	fmt.Printf("Removed %d Reading from %d \n", count, len(readings))
}

func (d *coredataCleaner) cleanValueDescriptors() {
	url := d.baseUrl + clients.ApiValueDescriptorRoute
	var valueDescriptors []models.ValueDescriptor
	err := request.Get(url, valueDescriptors)
	if err != nil {
		return
	}

	var count int
	for _, valueDescriptor := range valueDescriptors {
		err = request.Delete(url + config.PathId + valueDescriptor.Id)
		if err == nil {
			count = count +1
		}
	}
	fmt.Printf("Removed %d Value Descriptors from %d \n", count, len(valueDescriptors))
}

func (d *coredataCleaner) cleanEvents() {
	url := d.baseUrl + clients.ApiEventRoute + "/scruball"
	err := request.Delete(url)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print("All Events have been removed \n")
	}
}
