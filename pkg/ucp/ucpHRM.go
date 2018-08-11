package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/thebsdbox/diver/pkg/ucp/types"
)

// ConfigureHRM - This will toggle the functionality to enable Image Scanning
func (c *Client) ConfigureHRM(interlockConfig ucptypes.InterlockConfig) error {
	var hrmConfig ucptypes.HRMConfig
	url := fmt.Sprintf("%s/api/hrm", c.UCPURL)

	log.Debugf("Setting HRM set to [%t]", interlockConfig.InterlockEnabled)

	hrmConfig.Arch = interlockConfig.Arch
	hrmConfig.HTTPPort = interlockConfig.HTTPPort
	hrmConfig.HTTPSPort = interlockConfig.HTTPSPort

	b, err := json.Marshal(hrmConfig)
	if err != nil {
		return nil
	}

	if !interlockConfig.InterlockEnabled {
		_, err = c.delRequest(url, nil)
	} else {
		_, err = c.postRequest(url, b)
	}

	if err != nil {
		return err
	}

	return nil
}
