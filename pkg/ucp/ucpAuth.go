package ucp

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/thebsdbox/diver/pkg/ucp/types"
)

// AuthStatus - will return the current logged in user
func (c *Client) AuthStatus() (*ucptypes.Account, error) {
	log.Debugln("Retrieving the current authorisation status")
	url := fmt.Sprintf("%s/id/", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {

		return nil, fmt.Errorf("Unable to authenticate with existing token")
	}

	var a ucptypes.Account

	err = json.Unmarshal(response, &a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

//GetClientBundle - will download the UCP Client Bundle
func (c *Client) GetClientBundle() error {

	log.Infoln("Downloading the UCP Client Bundle")
	// Create the file
	out, err := os.Create("ucp-bundle.zip")
	if err != nil {
		return err
	}
	defer out.Close()

	log.Debugln("Retrieving Client Bundle")
	url := fmt.Sprintf("%s/api/clientbundle", c.UCPURL)

	response, err := c.getRequest(url, nil)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, bytes.NewReader(response))
	if err != nil {
		return err
	}

	return nil
}

// NewAccount - Creates a new account within UCP
func NewAccount(fullname, username, password string, isActive, isOrg, isAdmin, searchLDAP bool) *ucptypes.Account {
	return &ucptypes.Account{
		FullName:   fullname,
		IsActive:   isActive,
		IsAdmin:    isAdmin,
		IsOrg:      isOrg,
		Name:       username,
		Password:   password,
		SearchLDAP: searchLDAP,
	}
}

// NewUser - Creates a new user accound
func NewUser(fullname, username, password string, isActive, isAdmin, searchLDAP bool) *ucptypes.Account {
	return &ucptypes.Account{
		FullName:   fullname,
		IsActive:   isActive,
		IsAdmin:    isAdmin,
		IsOrg:      false,
		Name:       username,
		Password:   password,
		SearchLDAP: searchLDAP,
	}
}

// NewOrg - Creates a new organisation
func NewOrg(fullname, username, password string, isActive, isAdmin, searchLDAP bool) *ucptypes.Account {
	return &ucptypes.Account{
		FullName:   fullname,
		IsActive:   isActive,
		IsAdmin:    isAdmin,
		IsOrg:      true,
		Name:       username,
		Password:   password,
		SearchLDAP: searchLDAP,
	}
}

//AddAccount - adds a new account to UCP
func (c *Client) AddAccount(account *ucptypes.Account) error {

	if len(account.Password) < 8 && !account.IsOrg {
		// UCP appears to have a minimal password lenght (not sure about complexity)
		return fmt.Errorf("Password doesn't meet complexity levels or length")
	}
	if account.IsOrg {
		log.Infof("Creating account for Organisation [%s]", account.Name)
	} else {
		if account.FullName != "" {
			log.Infof("Creating account for user [%s]", account.FullName)

		} else {
			log.Infof("Creating account for user [%s]", account.Name)
		}
	}
	url := fmt.Sprintf("%s/accounts", c.UCPURL)

	b, err := json.Marshal(account)

	log.Debugf("%s", string(b))
	if err != nil {
		return err
	}
	response, err := c.postRequest(url, b)
	if err != nil {
		log.Debugf("response: %v", err)
		parseerr := ParseUCPError([]byte(response))
		if parseerr != nil {
			log.Debugf("%s", response)
			log.Errorf("Error parsing UCP error: %v", err)
		}
		return err
	}

	log.Debugf("%v", string(response))

	return nil
}

//DeleteAccount - deletes an account in UCP
func (c *Client) DeleteAccount(account string) error {
	log.Infof("Deleting account [%s]", account)

	url := fmt.Sprintf("%s/accounts/%s", c.UCPURL, account)

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

//AddTeamToOrganisation - adds a team to an existing organisation
func (c *Client) AddTeamToOrganisation(team *ucptypes.Team, org string) error {
	log.Infof("Creating team [%s]", team.Name)

	url := fmt.Sprintf("%s/accounts/%s/teams", c.UCPURL, org)

	b, err := json.Marshal(team)

	log.Debugf("%s", string(b))
	if err != nil {
		return err
	}
	response, err := c.postRequest(url, b)
	if err != nil {
		parseerr := ParseUCPError(response)
		if err != nil {
			log.Errorf("Error parsing UCP error: %v", parseerr)
		}
		return err
	}

	log.Debugf("%v", string(response))

	return nil
}

//DeleteTeamFromOrganisation - deletes an account in UCP
func (c *Client) DeleteTeamFromOrganisation(team, org string) error {
	log.Infof("Deleting team [%s] from org [%s]", team, org)

	url := fmt.Sprintf("%s/accounts/%s/teams/%s", c.UCPURL, org, team)

	response, err := c.delRequest(url, nil)
	if err != nil {
		parseerr := ParseUCPError(response)
		if parseerr != nil {
			log.Errorf("Error parsing UCP error: %v", parseerr)
		}
		return err
	}
	return nil
}

//AddUserToTeam - adds a team to an existing organisation
func (c *Client) AddUserToTeam(user, org, team string) error {
	log.Infof("Adding user [%s] to team [%s] in organisation [%s]", user, team, org)

	url := fmt.Sprintf("%s/accounts/%s/teams/%s/members/%s", c.UCPURL, org, team, user)

	// TODO - user can be an admin, but does he need to be a team admin (not sure what it does)
	response, err := c.putRequest(url, []byte(`{}`))
	if err != nil {
		log.Debugf("%v", string(response))
		parserr := ParseUCPError(response)
		if parserr != nil {
			log.Errorf("Error parsing UCP error: %v", err)
		}
		return err
	}

	log.Debugf("%v", string(response))

	return nil
}

//DelUserFromTeam - deletes a user from a team in an existing organisation
func (c *Client) DelUserFromTeam(user, org, team string) error {
	log.Infof("Deleting user [%s] from team [%s] in organisation [%s]", user, team, org)

	url := fmt.Sprintf("%s/accounts/%s/teams/%s/members/%s", c.UCPURL, org, team, user)

	// TODO - user can be an admin, but does he need to be a team admin (not sure what it does)
	response, err := c.delRequest(url, nil)
	if err != nil {
		if len(response) > 0 {
			log.Debugf("%v", string(response))
		} else {
			return err
		}
		parseerr := ParseUCPError(response)
		if parseerr != nil {
			log.Errorf("Error parsing UCP error: %v", parseerr)
		}
		return err
	}

	if len(response) > 0 {
		log.Debugf("%v", string(response))
	}

	return nil
}

// AddUserToOrganisation - Adds a user to an existing organisation TODO
func (c *Client) AddUserToOrganisation(user, org string) error {
	return nil
}

//ImportAccountsFromCSV - This will import accounts from a CSV file and create/delete them as described in the CSV
func (c *Client) ImportAccountsFromCSV(path string) error {
	csvFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return err
	}

	var accounts []ucptypes.Account

	var action []int8 // 0 = create, 1 = delete, 2 = update, 3 = ignore

	// Parsing from a CSV to a GO struct is a bit messy, using the ParseBool method
	log.Infof("Parsing CSV file [%s]", path)
	for _, line := range csvLines {

		// Check for Header line
		if line[1] == "Full Name" && line[2] == "Active account" {
			log.Debugf("Ignoring Header")
			continue
		}

		// Parse the action
		switch line[0] {
		case "create":
			log.Debug("Creating a new account")
			action = append(action, 0)
		case "delete":
			log.Debug("Deleting an existing account")
			action = append(action, 1)
		case "update":
			log.Debug("Updating an existing account")
			action = append(action, 2)
		case "ignore":
			log.Debugf("Ignoring the account [%s]", line[5])
			action = append(action, 3)
		default:
			return fmt.Errorf("Unknown action [%s] on account", line[0])
		}

		var acct ucptypes.Account
		acct.FullName = line[1]
		// Is Active
		var b bool
		b, err := strconv.ParseBool(line[2])
		if err != nil {
			return err
		}
		acct.IsActive = b
		// Is Admin
		b, err = strconv.ParseBool(line[3])
		if err != nil {
			return err
		}
		acct.IsAdmin = b

		// Is Org
		b, err = strconv.ParseBool(line[4])
		if err != nil {
			return err
		}
		acct.IsOrg = b

		// Name       string `json:"name"`
		acct.Name = line[5]
		// Password   string `json:"password"`
		acct.Password = line[6]
		// SearchLDAP bool   `json:"searchLDAP"`
		b, err = strconv.ParseBool(line[7])
		if err != nil {
			return err
		}
		acct.SearchLDAP = b
		accounts = append(accounts, acct)
	}

	log.Debugf("About to add %d accounts", len(accounts))

	//TODO - loop through accounts array and add accounts

	if len(accounts) != len(action) {
		return fmt.Errorf("Actions doesn't match the number of accounts")
	}

	for i := range accounts {
		switch action[i] {
		case 0:
			c.AddAccount(&accounts[i])
		case 1:
			c.DeleteAccount(accounts[i].Name)
		case 2:
			log.Warnf("Not implemented yet") //TODO
		case 3:
			// Account is ignored
		default:
			return fmt.Errorf("Unknown action being performed on user [%s]", accounts[i].FullName)
		}
	}

	return nil
}

//ExportAccountsToCSV -
func (c *Client) ExportAccountsToCSV(path string) error {

	// Default Query
	var accountQuery ucptypes.Account

	a, err := c.GetAccounts(accountQuery, -1)
	if err != nil {
		return err
	}

	csvFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// Write a header to the top of the CSV making it easier to understand
	var csvHeader = []string{
		"Action",
		"Full Name",
		"Active account",
		"Administrative User",
		"Is Organisation",
		"ID",
		"Password",
		"LDAP User"}
	writer.Write(csvHeader)

	for _, acct := range a.Accounts {
		var csvRow = []string{
			"ignore", // this will step over the account when parsing the CSV
			acct.FullName,
			strconv.FormatBool(acct.IsActive),
			strconv.FormatBool(acct.IsAdmin),
			strconv.FormatBool(acct.IsOrg),
			acct.Name,
			acct.Password,
			strconv.FormatBool(acct.SearchLDAP)}
		writer.Write(csvRow)
	}
	return nil
}

//CreateExampleAccountCSV -
func CreateExampleAccountCSV() error {
	path := "example_accounts.csv"

	csvFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	action := "create"

	acct := ucptypes.Account{
		FullName:   "John Doe",
		IsActive:   true,
		IsAdmin:    true,
		IsOrg:      false,
		Name:       "jdoe",
		Password:   "lockedDown",
		SearchLDAP: false,
	}

	for i := 1; i <= 100; i++ {
		var csvString = []string{
			action,
			fmt.Sprintf("%s%d", acct.FullName, i),
			strconv.FormatBool(acct.IsActive),
			strconv.FormatBool(acct.IsAdmin),
			strconv.FormatBool(acct.IsOrg),
			fmt.Sprintf("%s%d", acct.Name, i),
			acct.Password,
			strconv.FormatBool(acct.SearchLDAP),
		}
		err = writer.Write(csvString)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAccounts - This will get all accounts
func (c *Client) GetAccounts(query ucptypes.Account, count int) (*ucptypes.AccountList, error) {

	// Build the URL (TODO set limit)
	var url string
	if count > 0 {
		url = fmt.Sprintf("%s/accounts/?limit=%d", c.UCPURL, count)
	} else {
		// DEFAULT to 1000 records
		url = fmt.Sprintf("%s/accounts/?limit=1000", c.UCPURL)
	}
	if query.IsOrg == true {
		log.Debugln("Retrieving Organisations")
		url = url + "&filter=orgs"
	}

	if query.IsAdmin == true {
		log.Debugln("Retrieving Admins")
		url = url + "&filter=admins"
	}

	if query.IsActive == true {
		log.Debugln("Retrieving InActive users")
		url = url + "&filter=inactive"
	}

	log.Debugf("Built URL [%s]", url)
	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, fmt.Errorf("%s", string(response))
	}

	var a ucptypes.AccountList

	err = json.Unmarshal(response, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil

}
