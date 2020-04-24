package purge

import (
	"context"
	"fmt"

	"github.com/edgexfoundry-holding/edgex-cli/config"
	request "github.com/edgexfoundry-holding/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type MetadataCleaner interface {
	Purge()
	cleanDevices()
	cleanDeviceServices()
	cleanDeviceProfiles()
	cleanAddressables()
}

type metadataCleaner struct {
	baseUrl string
}

// NewMetadataCleaner creates an instance of MetadataCleaner
func NewMetadataCleaner() MetadataCleaner {
	fmt.Println("\n * core-metadata")
	return &metadataCleaner{
		baseUrl: config.Conf.Clients["Metadata"].Url(),
	}
}

func (d *metadataCleaner) Purge() {
	d.cleanDevices()
	d.cleanDeviceServices()
	d.cleanDeviceProfiles()
	d.cleanAddressables()
}

func (d *metadataCleaner) cleanDevices() {
	ctx, _ := context.WithCancel(context.Background())
	url := d.baseUrl + clients.ApiDeviceRoute
	mdc := metadata.NewDeviceClient(local.New(url))

	devices, err := mdc.Devices(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int
	for _, device := range devices {
		err = request.Delete(url + config.PathId + device.Id)
		if err != nil {
			fmt.Printf("Failed to delete Device with id %s because of error: %s", device.Id, err)
		} else {
			count = count + 1
		}
	}
	fmt.Printf("Removed %d Devices from %d \n", count, len(devices))
}

func (d *metadataCleaner) cleanDeviceServices() {
	url := d.baseUrl + clients.ApiDeviceServiceRoute
	var deviceServices []models.DeviceService
	err := request.Get(url, &deviceServices)
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int
	for _, deviceService := range deviceServices {
		err = request.Delete(url + config.PathId + deviceService.Id)
		if err != nil {
			fmt.Printf("Failed to delete Device Service with id %s because of error: %s", deviceService.Id, err)
		} else {
			count = count + 1
		}
	}

	fmt.Printf("Removed %d Device Service from %d \n", count, len(deviceServices))
}

func (d *metadataCleaner) cleanDeviceProfiles() {
	url := d.baseUrl + clients.ApiDeviceProfileRoute
	var deviceProfiles []models.DeviceProfile
	err := request.Get(url, &deviceProfiles)
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int
	for _, deviceProfile := range deviceProfiles {
		err = request.Delete(url + config.PathId + deviceProfile.Id)
		if err != nil {
			fmt.Printf("Failed to delete Device Profile with id %s because of error: %s", deviceProfile.Id, err)
		} else {
			count = count + 1
		}
	}
	fmt.Printf("Removed %d Device Profile from %d \n", count, len(deviceProfiles))
}

func (d *metadataCleaner) cleanAddressables() {
	url := d.baseUrl + clients.ApiAddressableRoute
	var addressables []models.Addressable
	err := request.Get(url, &addressables)
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int
	for _, addr := range addressables {
		err = request.Delete(url + config.PathId + addr.Id)
		if err != nil {
			fmt.Printf("Failed to delete Addressable with id %s because of error: %s", addr.Id, err)
		} else {
			count = count + 1
		}
	}
	fmt.Printf("Removed %d Addressable from %d \n", count, len(addressables))
}
