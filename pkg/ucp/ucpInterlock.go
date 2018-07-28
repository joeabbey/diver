package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/thebsdbox/diver/pkg/ucp/types"
)

// ConfigureInterlock - This will toggle the functionality to enable Image Scanning
func (c *Client) ConfigureInterlock(config ucptypes.InterlockConfig) error {
	url := fmt.Sprintf("%s/api/interlock", c.UCPURL)

	log.Debugf("Setting interlock set to [%t]", config.InterlockEnabled)

	b, err := json.Marshal(config)
	if err != nil {
		return nil
	}

	if !config.InterlockEnabled {
		_, err = c.delRequest(url, nil)
	} else {
		_, err = c.postRequest(url, b)
	}

	if err != nil {
		return err
	}

	return nil
}
