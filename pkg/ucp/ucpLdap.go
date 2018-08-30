package ucp

import (
	"encoding/json"
	"fmt"

	"github.com/thebsdbox/diver/pkg/ucp/types"
)

//GetLDAPInfo - returns all information about the LDAP configuration
func (c *Client) GetLDAPInfo() (*ucptypes.LDAPConfig, error) {
	url := fmt.Sprintf("%s/enzi/v0/config/auth/ldap", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var l ucptypes.LDAPConfig

	err = json.Unmarshal(response, &l)
	if err != nil {
		return nil, err
	}

	return &l, nil
}

//SetLDAPInfo - returns all information about the LDAP configuration
func (c *Client) SetLDAPInfo(cfg *ucptypes.LDAPConfig) error {
	url := fmt.Sprintf("%s/enzi/v0/config/auth/ldap", c.UCPURL)

	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	response, err := c.putRequest(url, b)
	if err != nil {
		ParseUCPError(response)
		return err
	}

	return nil
}
