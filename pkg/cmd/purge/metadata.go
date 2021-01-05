package purge

import (
	"context"
	"fmt"

	"github.com/edgexfoundry/edgex-cli/config"
	request "github.com/edgexfoundry/edgex-cli/pkg"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/metadata"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/urlclient/local"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
)

type metadataCleaner struct {
	baseUrl string
	ctx     context.Context
}

// NewMetadataCleaner creates an instance of MetadataCleaner
func NewMetadataCleaner(ctx context.Context) Purgeable {
	fmt.Println("\n * core-metadata")
	return &metadataCleaner{
		baseUrl: config.Conf.Clients["Metadata"].Url(),
		ctx:     ctx,
	}
}

func (d *metadataCleaner) Purge() {
	d.cleanDevices()
	d.cleanDeviceServices()
	d.cleanDeviceProfiles()
	d.cleanAddressables()
}

func (d *metadataCleaner) cleanDevices() {
	url := d.baseUrl + clients.ApiDeviceRoute
	mdc := metadata.NewDeviceClient(local.New(url))

	devices, err := mdc.Devices(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	var count int
	for _, device := range devices {
		err = request.Delete(d.ctx, url+config.PathId+device.Id)
		if err == nil {
			count = count + 1
		}
	}
	fmt.Printf("Removed %d Devices from %d \n", count, len(devices))
}

func (d *metadataCleaner) cleanDeviceServices() {
	url := d.baseUrl + clients.ApiDeviceServiceRoute
	var deviceServices []models.DeviceService
	err := request.Get(d.ctx, url, &deviceServices)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	var count int
	for _, deviceService := range deviceServices {
		err = request.Delete(d.ctx, url+config.PathId+deviceService.Id)
		if err == nil {
			count = count + 1
		}
	}

	fmt.Printf("Removed %d Device Services from %d \n", count, len(deviceServices))
}

func (d *metadataCleaner) cleanDeviceProfiles() {
	url := d.baseUrl + clients.ApiDeviceProfileRoute
	var deviceProfiles []models.DeviceProfile
	err := request.Get(d.ctx, url, &deviceProfiles)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	var count int
	for _, deviceProfile := range deviceProfiles {
		err = request.Delete(d.ctx, url+config.PathId+deviceProfile.Id)
		if err == nil {
			count = count + 1
		}
	}
	fmt.Printf("Removed %d Device Profiles from %d \n", count, len(deviceProfiles))
}

func (d *metadataCleaner) cleanAddressables() {
	url := d.baseUrl + clients.ApiAddressableRoute
	var addressables []models.Addressable
	err := request.Get(d.ctx, url, &addressables)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	var count int
	for _, addr := range addressables {
		err = request.Delete(d.ctx, url+config.PathId+addr.Id)
		if err == nil {
			count = count + 1
		}
	}
	fmt.Printf("Removed %d Addressables from %d \n", count, len(addressables))
}
