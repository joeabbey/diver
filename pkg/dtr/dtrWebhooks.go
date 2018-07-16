package dtr

import (
	"encoding/json"
	"fmt"

	"github.com/thebsdbox/diver/pkg/dtr/types"
)

//ListWebhooks -
func (c *Client) ListWebhooks() ([]dtrTypes.DTRWebHook, error) {

	url := fmt.Sprintf("%s/api/v0/webhooks?webhookType=any&refresh_token=%s", c.DTRURL, c.Token)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	//log.Debugf("%v", string(response))
	var info []dtrTypes.DTRWebHook

	err = json.Unmarshal(response, &info)
	if err != nil {
		return nil, err
	}
	return info, nil

}
