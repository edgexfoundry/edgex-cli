package version

import (
	"encoding/json"
	"net/http"

	"github.com/edgexfoundry-holding/edgex-cli/config"

	"github.com/edgexfoundry/go-mod-core-contracts/clients"
)

type Version struct {
	Version string `json:"version" yaml:"version,omitempty"`
}
// GetEdgeXVersion returns the version of core-data microservice
func GetEdgeXVersion() (version Version, err error) {
	url := config.Conf.Clients["CoreData"].Url() + clients.ApiVersionRoute
	resp, err := http.Get(url)
	if err != nil {
		return version, err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&version)
	if err != nil {
		return Version{}, err
	}

	return version, nil
}
