package ucp

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

type roles struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	ServiceRole bool            `json:"servicerole"`
	Operations  json.RawMessage // Captures the raw output of the remaining json object
}

// TODO fix the rest of the struct
type collection struct {
	// "name": "Private",
	// "path": "/Shared/Private",
	// "id": "private",
	// "parent_ids": [
	//   "root",
	//   "swarm",
	//   "shared"
	// ],
	// "label_constraints": [],
	// "legacylabelkey": "",
	// "legacylabelvalue": "",
	// "created_at": "2018-06-11T17:16:14.124Z",
	// "updated_at": "2018-06-11T17:16:14.124Z"
	Name string `json:"name"`
	Path string `json:"path"`
	ID   string `json:"id"`
}

// A grant is based upon three keys:
// -- ObjectID == Collection
// -- RoleID == Links the role that is applied (rights)
// -- SubjectID == User that has is linked to the collection with the appropriate rights

type grant struct {
	ObjectID  string `json:"objectID"`
	RoleID    string `json:"roleID"`
	SubjectID string `json:"subjectID"`
}

//collection’, 'namespace’, or 'grantobject

const (
	// GrantCollection - (default) specifies a grant is created against a collection
	GrantCollection uint = 1 << iota

	// GrantNamespace - A grant is made against a namespace (kubernetes)
	GrantNamespace

	// GrantObject - kubernetesnamespaces target, which is used to give grants against all Kubernetes namespaces.
	GrantObject
)

//GetOrg - TODO
func (c *Client) GetOrg(orgName string) error {
	log.Debugf("Searching for Org [%s]", orgName)
	return nil
}

func (c *Client) returnAllRoles() ([]roles, error) {

	url := fmt.Sprintf("%s/roles", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var r []roles

	log.Debugf("Parsing all roles")
	err = json.Unmarshal(response, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

//GetRoles - This will return a list of services
func (c *Client) GetRoles() error {

	r, err := c.returnAllRoles()
	if err != nil {
		return err
	}

	fmt.Printf("ID\t\tService Account\tName\n")

	for i := range r {
		fmt.Printf("%s\t%t\t%s\n", r[i].ID, r[i].ServiceRole, r[i].Name)
	}
	return nil
}

//GetRoleRuleset - This will return a list of rules attached to a role
func (c *Client) GetRoleRuleset(role string) (string, error) {

	r, err := c.returnAllRoles()
	if err != nil {
		return "", err
	}

	for i := range r {
		if role == r[i].Name {
			return string(r[i].Operations), nil
		}

	}
	return "", fmt.Errorf("Unable to find role [%s]", role)
}

//CreateRole - This set the role of a user in an organisation
func (c *Client) CreateRole(name, id, ruleset string, serviceAccount bool) error {

	url := fmt.Sprintf("%s/roles", c.UCPURL)

	newrole := roles{
		ID:          id,
		Name:        name,
		ServiceRole: serviceAccount,
		Operations:  json.RawMessage(ruleset),
	}

	b, err := json.Marshal(newrole)

	if err != nil {
		return err
	}

	response, err := c.postRequest(url, b)
	if err != nil {
		return nil
	}

	log.Debugf("%v", string(response))

	return nil
}

func (c *Client) returnAllCollections() ([]collection, error) {

	url := fmt.Sprintf("%s/collections?limit=1000", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var collections []collection

	log.Debugf("Parsing all collections")
	err = json.Unmarshal(response, &collections)
	if err != nil {
		return nil, err
	}
	return collections, nil
}

//GetGrants - This will return a list of all grants, it can also resolve the UUIDs to names
func (c *Client) GetGrants(resolve bool) error {

	url := fmt.Sprintf("%s/collectionGrants?limit=1000", c.UCPURL)
	log.Debugf("built URL [%s]", url)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return err
	}

	type subjects struct {
		ID             string  `json:"id"`
		Type           string  `json:"subject_type"`
		SubjectAccount Account `json:"account"`
	}

	var grants struct {
		Grants   []grant    `json:"grants"`
		Subjects []subjects `json:"subjects"`
	}

	var r []roles
	var collections []collection
	// If resolving cache the roles, and collections before hand (speed up the resolution process)
	if resolve {
		r, err = c.returnAllRoles()
		if err != nil {
			return err
		}
		collections, err = c.returnAllCollections()
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(response, &grants)
	for i := range grants.Grants {

		subject := grants.Grants[i].SubjectID
		role := grants.Grants[i].RoleID
		object := grants.Grants[i].ObjectID
		if resolve {
			for x := range grants.Subjects {
				// Don't replace the UUID with a blank name (as that is worse than useless)
				if subject == grants.Subjects[x].ID && grants.Subjects[x].SubjectAccount.Name != "" {
					subject = grants.Subjects[x].SubjectAccount.Name
				}
			}
			for x := range r {
				if role == r[x].ID {
					role = r[x].Name
				}
			}
			for x := range collections {
				if object == collections[x].ID {
					object = collections[x].Path
				}
			}
		}

		fmt.Printf("%s\t\t%s\t\t%s\n", subject, role, object)

	}

	log.Debugf("%v", grants.Subjects)
	return nil
}

//SetGrant - This takes a subject and a role (ruleset) and applies it to a collection
func (c *Client) SetGrant(collection, role, subject string, flags uint) error {

	// Parser flags
	var grantType string
	switch flags {
	case (GrantCollection):
		grantType = "collection"
	case (GrantNamespace):
		grantType = "namespace"
	case (GrantObject):
		grantType = "grantobject"
	default:
		return fmt.Errorf("Unknown Grant Type")
	}

	url := fmt.Sprintf("%s/collectionGrants/%s/%s/%s?type=%s", c.UCPURL, subject, collection, role, grantType)
	log.Debugf("built URL [%s]", url)

	_, err := c.putRequest(url, nil)
	if err != nil {

		return err
	}

	return nil
}
