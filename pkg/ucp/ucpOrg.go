package ucp

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"

	log "github.com/Sirupsen/logrus"
	"github.com/thebsdbox/diver/pkg/ucp/types"
)

func (c *Client) returnAllRoles() ([]ucptypes.Roles, error) {

	url := fmt.Sprintf("%s/roles", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var r []ucptypes.Roles

	log.Debugf("Parsing all roles")
	err = json.Unmarshal(response, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (c *Client) returnAllTeamsFromOrg(org string) (*ucptypes.Teams, error) {

	url := fmt.Sprintf("%s/accounts/%s/teams?limit=10", c.UCPURL, org)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var t ucptypes.Teams

	log.Debugf("Parsing team JSON")
	err = json.Unmarshal(response, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// TotalRole - This will return all roles available in the system
func (c *Client) TotalRole() ([]byte, error) {

	url := fmt.Sprintf("%s/totalRole", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}
	return response, nil
}

//GetTeams - This will print a list of teams
func (c *Client) GetTeams(org string) error {
	t, err := c.returnAllTeamsFromOrg(org)
	if err != nil {
		return err
	}

	log.Debugf("Found %d teams for organisation %s", len(t.Teams), org)

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "Name\tID\tOrganisation\tDescription\tMember Count")

	for i := range t.Teams {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n", t.Teams[i].Name, t.Teams[i].ID, t.Teams[i].OrgID, t.Teams[i].Description, t.Teams[i].MembersCount)
	}
	w.Flush()
	return nil
}

//GetRoles - This will print a list of services
func (c *Client) GetRoles() error {

	r, err := c.returnAllRoles()
	if err != nil {
		return err
	}

	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "Name\tID\tService Account")

	for i := range r {
		fmt.Fprintf(w, "%s\t%s\t%t\n", r[i].Name, r[i].ID, r[i].SystemRole)
	}
	w.Flush()

	return nil
}

//GetRoleRuleset - This will return a list of rules attached to a role
func (c *Client) GetRoleRuleset(role string, id string) (string, error) {

	r, err := c.returnAllRoles()
	if err != nil {
		return "", err
	}

	if role != "" {
		for i := range r {
			if role == r[i].Name {
				return string(r[i].Operations), nil
			}
		}
	}

	if id != "" {
		for i := range r {
			if id == r[i].ID {
				return string(r[i].Operations), nil
			}
		}
	}
	if role != "" {
		return "", fmt.Errorf("Unable to find role [%s]", role)
	}
	return "", fmt.Errorf("Unable to find ID [%s]", id)

}

//CreateRole - This set the role of a user in an organisation
func (c *Client) CreateRole(name, id, ruleset string, systemRole bool) error {

	url := fmt.Sprintf("%s/roles", c.UCPURL)

	newrole := ucptypes.Roles{
		ID:         id,
		Name:       name,
		SystemRole: systemRole,
		Operations: json.RawMessage(ruleset),
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

//DeleteRole - deletes a role in UCP
func (c *Client) DeleteRole(account string) error {
	log.Infof("Deleting Role [%s]", account)

	url := fmt.Sprintf("%s/roles/%s", c.UCPURL, account)

	response, err := c.delRequest(url, nil)
	if err != nil {
		parseerr := ParseUCPError(response)
		if parseerr != nil {
			log.Errorf("Error parsing UCP error: %v", err)
		}
		return err
	}
	return nil
}

//GetGrants - This will return a list of all grants, it can also resolve the UUIDs to names
func (c *Client) GetGrants(resolve bool) error {

	url := fmt.Sprintf("%s/collectionGrants?limit=1000", c.UCPURL)
	log.Debugf("built URL [%s]", url)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return fmt.Errorf("%s", response)
	}

	type subjects struct {
		ID             string           `json:"id"`
		Type           string           `json:"subject_type"`
		SubjectAccount ucptypes.Account `json:"account"`
	}

	var grants struct {
		Grants   []ucptypes.Grant `json:"grants"`
		Subjects []subjects       `json:"subjects"`
	}

	var r []ucptypes.Roles
	var collections []ucptypes.Collection
	// If resolving cache the roles, and collections before hand (speed up the resolution process)
	if resolve {
		r, err = c.returnAllRoles()
		if err != nil {
			return err
		}
		collections, err = c.GetCollections()
		if err != nil {
			return err
		}
	}

	err = json.Unmarshal(response, &grants)
	if err != nil {
		return err
	}
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', 0)
	fmt.Fprintln(w, "Subject\tRole\tCollection")
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
		fmt.Fprintf(w, "%s\t%s\t%s\n", subject, role, object)
	}
	w.Flush()
	return nil
}

//SetGrant - This takes a subject and a role (ruleset) and applies it to a collection
func (c *Client) SetGrant(collection, role, subject string, flags uint) error {

	// Parser flags
	var grantType string
	var roleID string
	switch flags {
	case (ucptypes.GrantCollection):
		grantType = "collection"
	case (ucptypes.GrantNamespace):
		grantType = "namespace"
	case (ucptypes.GrantObject):
		grantType = "grantobject"
	default:
		return fmt.Errorf("Unknown Grant Type")
	}

	matched, err := regexp.MatchString("[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", role)
	if !matched {
		//This is not a UUID, let's try to figure out what the UUID should be
		r, err := c.returnAllRoles()
		if err != nil {
			return err
		}

		//Search (potentially slowly) for the role
		for i := range r {
			if r[i].Name == role {
				roleID = r[i].ID
			}
		}
		if roleID == "" {
			log.Fatalf("Failed to lookup role: %s.", role)
		}
	} else {
		//Looks like a UUID let's use it
		roleID = role
	}

	url := fmt.Sprintf("%s/collectionGrants/%s/%s/%s?type=%s", c.UCPURL, subject, collection, roleID, grantType)
	log.Debugf("built URL [%s]", url)

	_, err = c.putRequest(url, nil)
	if err != nil {

		return err
	}

	return nil
}

//DeleteGrant - This function will remove an existing grant
func (c *Client) DeleteGrant(collection, role, subject string) error {

	url := fmt.Sprintf("%s/collectionGrants/%s/%s/%s", c.UCPURL, subject, collection, role)
	log.Debugf("built URL [%s]", url)

	_, err := c.delRequest(url, nil)
	if err != nil {

		return err
	}
	return nil
}
