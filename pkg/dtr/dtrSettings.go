package dtr

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/joeabbey/diver/pkg/dtr/types"
)

//DTRCreateRepoOnPush - This will toggle the functionality to enable repositories to be created on push
func (c *Client) DTRCreateRepoOnPush(enabled bool) error {

	url := fmt.Sprintf("%s/api/v0/meta/settings?refresh_token=%s", c.DTRURL, c.Token)

	var dtrSettings struct {
		EnableCreateOnPush bool `json:"createRepositoryOnPush"`
	}

	dtrSettings.EnableCreateOnPush = enabled
	log.Debugf("Setting EnableCreateOnPush set to [%t]", enabled)
	b, err := json.Marshal(dtrSettings)
	if err != nil {
		return nil
	}
	_, err = c.postRequest(url, b)
	if err != nil {
		return err
	}

	return nil
}

//DTRScanningEnable - This will toggle the functionality to enable Image Scanning
func (c *Client) DTRScanningEnable(enabled bool) error {

	url := fmt.Sprintf("%s/api/v0/meta/settings?refresh_token=%s", c.DTRURL, c.Token)

	var dtrSettings struct {
		EnableCreateOnPush bool `json:"scanningEnabled"`
	}

	dtrSettings.EnableCreateOnPush = enabled
	log.Debugf("Setting scanningEnabled set to [%t]", enabled)
	b, err := json.Marshal(dtrSettings)
	if err != nil {
		return nil
	}
	_, err = c.postRequest(url, b)
	if err != nil {
		return err
	}

	return nil
}

//DTROnlineScan - This will toggle the functionality to Sync online
func (c *Client) DTROnlineScan(enabled bool) error {

	url := fmt.Sprintf("%s/api/v0/meta/settings?refresh_token=%s", c.DTRURL, c.Token)

	var dtrSettings struct {
		EnableCreateOnPush bool `json:"scanningSyncOnline"`
	}

	dtrSettings.EnableCreateOnPush = enabled
	log.Debugf("Setting sync online set to [%t]", enabled)
	b, err := json.Marshal(dtrSettings)
	if err != nil {
		return nil
	}
	_, err = c.postRequest(url, b)
	if err != nil {
		return err
	}
	return nil
}

//DTRGetSettings - Return a struct of all DTR settings
func (c *Client) DTRGetSettings() (*dtrtypes.DTRSettings, error) {
	url := fmt.Sprintf("%s/api/v0/meta/settings?refresh_token=%s", c.DTRURL, c.Token)

	var settings dtrtypes.DTRSettings
	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(response, &settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}
