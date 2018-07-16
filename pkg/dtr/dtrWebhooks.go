package dtr

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/thebsdbox/diver/pkg/dtr/types"
)

//ListWebhooks -
func (c *Client) ListWebhooks() ([]dtrtypes.DTRWebHook, error) {

	url := fmt.Sprintf("%s/api/v0/webhooks?webhookType=any&refresh_token=%s", c.DTRURL, c.Token)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	//log.Debugf("%v", string(response))
	var info []dtrtypes.DTRWebHook

	err = json.Unmarshal(response, &info)
	if err != nil {
		return nil, err
	}
	return info, nil

}

//CreateWebhook -
func (c *Client) CreateWebhook(webhook dtrtypes.DTRWebHook) error {

	url := fmt.Sprintf("%s/api/v0/webhooks?refresh_token=%s", c.DTRURL, c.Token)

	b, err := json.Marshal(webhook)
	if err != nil {
		return err
	}
	log.Debugf("%s", b)
	response, err := c.postRequest(url, b)
	if err != nil {
		log.Debugf("%s", response)
		return err
	}

	return nil

}

//DeleteWebhook -
func (c *Client) DeleteWebhook(id string) error {

	url := fmt.Sprintf("%s/api/v0/webhooks/%s?refresh_token=%s", c.DTRURL, id, c.Token)

	response, err := c.delRequest(url, nil)
	if err != nil {
		log.Debugf("%s", response)
		return err
	}

	return nil

}
