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
)

// Account - Is the basic Account struct
type Account struct {
	FullName   string `json:"fullName"`
	IsActive   bool   `json:"isActive"`
	IsAdmin    bool   `json:"isAdmin"`
	IsOrg      bool   `json:"isOrg"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	SearchLDAP bool   `json:"searchLDAP"`
}

// Team - is the structure for defining a team
type Team struct {
	Description string `json:"description"`
	Name        string `json:"name"`
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

	log.Debugf("Retrieving Client Bundle")
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
func NewAccount(fullname, username, password string, isActive, isOrg, isAdmin, searchLDAP bool) *Account {
	return &Account{
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
func NewUser(fullname, username, password string, isActive, isAdmin, searchLDAP bool) *Account {
	return &Account{
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
func NewOrg(fullname, username, password string, isActive, isAdmin, searchLDAP bool) *Account {
	return &Account{
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
func (c *Client) AddAccount(account *Account) error {
	log.Infof("Creating account for user [%s]", account.FullName)

	url := fmt.Sprintf("%s/accounts", c.UCPURL)

	b, err := json.Marshal(account)

	log.Debugf("%s", string(b))
	if err != nil {
		return err
	}
	response, err := c.postRequest(url, b)
	if err != nil {
		err = parseUCPError(err.Error())
		if err != nil {
			log.Errorf("Error parsing UCP error: %v", err)
		}
		return err
	}

	log.Debugf("%v", string(response))

	return nil
}

//DeleteAccount - deletes an account in UCP
func (c *Client) DeleteAccount(account string) error {
	log.Infof("Deleting account for user [%s]", account)

	url := fmt.Sprintf("%s/accounts/%s", c.UCPURL, account)

	_, err := c.delRequest(url, nil)
	if err != nil {
		err = parseUCPError(err.Error())
		if err != nil {
			log.Errorf("Error parsing UCP error: %v", err)
		}
		return err
	}
	return nil
}

//AddTeamToOrganisation - adds a team to an existing organisation
func (c *Client) AddTeamToOrganisation(team *Team, org string) error {
	log.Infof("Creating team [%s]", team.Name)

	url := fmt.Sprintf("%s/accounts/%s/teams", c.UCPURL, org)

	b, err := json.Marshal(team)

	log.Debugf("%s", string(b))
	if err != nil {
		return err
	}
	response, err := c.postRequest(url, b)
	if err != nil {
		err = parseUCPError(err.Error())
		if err != nil {
			log.Errorf("Error parsing UCP error: %v", err)
			log.Debugf("%v", response)
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

	_, err := c.delRequest(url, nil)
	if err != nil {
		err = parseUCPError(err.Error())
		if err != nil {
			log.Errorf("Error parsing UCP error: %v", err)
		}
		return err
	}
	return nil
}

//ImportAccountsFromCSV -
func (c *Client) ImportAccountsFromCSV(path string) error {
	csvFile, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer csvFile.Close()

	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return nil
	}

	var accounts []Account

	// Parsing from a CSV to a GO struct is a bit messy, using the ParseBool method
	log.Infof("Parsing CSV file [%s]", path)
	for _, line := range csvLines {

		var acct Account
		acct.FullName = line[0]
		// Is Active
		var b bool
		b, err := strconv.ParseBool(line[1])
		if err != nil {
			return err
		}
		acct.IsActive = b
		// Is Admin
		b, err = strconv.ParseBool(line[2])
		if err != nil {
			return err
		}
		acct.IsAdmin = b

		// Is Org
		b, err = strconv.ParseBool(line[3])
		if err != nil {
			return err
		}
		acct.IsOrg = b

		// Name       string `json:"name"`
		acct.Name = line[4]
		// Password   string `json:"password"`
		acct.Password = line[5]
		// SearchLDAP bool   `json:"searchLDAP"`
		b, err = strconv.ParseBool(line[6])
		if err != nil {
			return err
		}
		acct.SearchLDAP = b
		accounts = append(accounts, acct)
	}

	log.Debugf("About to add %d accounts", len(accounts))

	//TODO - loop through accounts array and add accounts
	for _, acct := range accounts {
		c.AddAccount(&acct)
	}

	return nil
}

//CreateExampleAccountCSV -
func CreateExampleAccountCSV() error {
	path := "example_accounts.csv"

	csvFile, err := os.Create(path)
	if err != nil {
		return nil
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	acct := Account{
		FullName:   "John Candy",
		IsActive:   true,
		IsAdmin:    true,
		IsOrg:      false,
		Name:       "jcandy",
		Password:   "Gr3At0utd00r5",
		SearchLDAP: false,
	}

	var csvString = []string{acct.FullName,
		strconv.FormatBool(acct.IsActive),
		strconv.FormatBool(acct.IsAdmin),
		strconv.FormatBool(acct.IsOrg),
		acct.Name,
		acct.Password,
		strconv.FormatBool(acct.SearchLDAP)}

	return writer.Write(csvString)
}
