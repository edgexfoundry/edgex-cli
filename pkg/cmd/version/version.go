package version

import (
	"encoding/json"
	"net/http"

	"github.com/spf13/viper"
)

//TODO move Version struct into go-mod-core-contracts
type Version struct {
	Version string `json:"version" yaml:"version,omitempty"`
}
// GetEdgeXVersion returns the version of core-data microservice
func GetEdgeXVersion(port string) (version Version, err error) {
	host := viper.GetString("dataservice.host")
	url := "http://" + host + ":" + port + "/api/version"

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
