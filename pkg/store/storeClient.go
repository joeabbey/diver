package store

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

// Client - Is the basic Client struct
type Client struct {
	Username   string
	Password   string
	STOREURL   string
	HUBURL     string
	IgnoreCert bool
	Token      string
	ID         string
}

// NewBasicAuthClient - Creates a basic client to connecto the UCP
func NewBasicAuthClient(username, password, url string, ignoreCert bool) *Client {
	return &Client{
		Username:   username,
		Password:   password,
		STOREURL:   url,
		IgnoreCert: ignoreCert,
	}
}

// Connect - Will attempt to connect to UCP
func (c *Client) Connect() error {
	if c.Username == "" {
		return fmt.Errorf("Username hasn't been entered")
	}

	if c.Password == "" {
		return fmt.Errorf("Password is blank")
	}

	if c.STOREURL == "" {
		return fmt.Errorf("URL hasn't been entered")
	}
	// Add the /auth/log to the URL
	url := fmt.Sprintf("%s/users/login", c.STOREURL)

	data := map[string]string{
		"username": c.Username,
		"password": c.Password,
	}
	b, err := json.Marshal(data)

	if err != nil {
		return err
	}

	response, err := c.postRequest(url, b)
	if err != nil {
		return err
	}
	log.Debugf("%v", string(response))
	var responseData map[string]interface{}
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return err
	}
	if responseData["token"] != nil {
		c.Token = responseData["token"].(string)
		user, err := c.userInfo(c.Username)
		if err != nil {
			return err
		}
		c.ID = user.ID
	} else {
		return fmt.Errorf("No Authorisation token returned")
	}
	return nil
}

// Disconnect - TODO
func (c *Client) Disconnect() error {
	return nil
}

// POST data to the server and return the response as bytes
func (c *Client) postRequest(url string, d []byte) ([]byte, error) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}

	// Add authorisation token to HTTP header
	if len(c.Token) != 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}
	req.Header.Add("Content-Type", "application/json")

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GET data from the server and return the response as bytes
func (c *Client) getRequest(url string, d []byte) ([]byte, error) {

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}

	// Add authorisation token to HTTP header
	if len(c.Token) != 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
		log.Debugf("Adding AuthHeader = [%v]", req.Header)
	}

	req.Header.Add("Content-Type", "application/json")

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GET data from the server and stream the output
func (c *Client) getRequestStream(url string, d []byte) error {

	req, err := http.NewRequest("GET", url, bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	// Add authorisation token to HTTP header
	if len(c.Token) != 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	req.Header.Add("Content-Type", "application/json")

	err = c.doStreamRequest(req)
	if err != nil {
		return err
	}
	return nil
}

// postRequestStream, this will stream the output of a post request
func (c *Client) postRequestStream(url string, d []byte) error {
	log.Debugln("Creating a new Streaming POST request")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(d))
	if err != nil {
		return err
	}

	// Add authorisation token to HTTP header
	if len(c.Token) != 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	req.Header.Add("Content-Type", "application/json")
	log.Debugln("Starting stream...")

	err = c.doStreamRequest(req)
	if err != nil {
		return err
	}
	return nil
}

// PUT will update an existing element
func (c *Client) putRequest(url string, d []byte) ([]byte, error) {

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}

	// Add authorisation token to HTTP header
	if len(c.Token) != 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// DELETE data from the server and return the response as bytes
func (c *Client) delRequest(url string, d []byte) ([]byte, error) {

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}

	// Add authorisation token to HTTP header
	if len(c.Token) != 0 {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	bytes, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (c *Client) doStreamRequest(req *http.Request) error {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.IgnoreCert},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	dec := json.NewDecoder(resp.Body)

	for {
		var message struct {
			Status string `json:"status,omitempty"`

			Stream string `json:"stream,omitempty"`
			ID     string `json:"id,omitempty"`
		}

		err := dec.Decode(&message)

		if resp.StatusCode != http.StatusOK {
			log.Fatalf("Status code is not OK: %v (%s)", resp.StatusCode, resp.Status)
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("%s", err)
		}
		log.Debugf("%v", message)
		if message.Status != "" {
			log.Infof("Building on %s", message.ID)
		}
		if message.Stream != "\n" {
			fmt.Printf("%s", message.Stream)
		}
	}

	return nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: c.IgnoreCert},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 2xx Success / 3xx Redirection
	if resp.StatusCode < 400 {
		log.Debugf("[success] HTTP Status code %d", resp.StatusCode)
		return body, nil
	}
	// The error code is > 400
	log.Debugf("HTTP Error code: %d for URL: %s", resp.StatusCode, req.URL.String())

	// Catches the "Majority" of expected responses
	switch resp.StatusCode {
	case 400:
		return body, fmt.Errorf("Code %d, Bad Request", resp.StatusCode)
	case 401:
		return body, fmt.Errorf("Code %d, Unauthorised", resp.StatusCode)
	case 402:
		return body, fmt.Errorf("Code %d, Payment Required", resp.StatusCode) //unused
	case 403:
		return body, fmt.Errorf("Code %d, Forbidden", resp.StatusCode)
	case 404:
		return body, fmt.Errorf("Code %d, Not Found", resp.StatusCode)
	case 405:
		return body, fmt.Errorf("Code %d, Method Not Allowed", resp.StatusCode)
	case 500:
		return body, fmt.Errorf("Code %d, Internal Server Error", resp.StatusCode)
	case 501:
		return body, fmt.Errorf("Code %d, Not Implemented", resp.StatusCode)
	case 502:
		return body, fmt.Errorf("Code %d, Bad Gateway", resp.StatusCode)
	case 503:
		return body, fmt.Errorf("Code %d, Service Unavailable", resp.StatusCode)
	case 504:
		return body, fmt.Errorf("Code %d, Gateway Timeout", resp.StatusCode)
	default:
		log.Debugf("[Untrapped return code] %d", resp.StatusCode)
		return body, fmt.Errorf("Code %s", resp.Status)

	}
}

type internal struct {
	StoreAddress string `json:"storeAddress"`
	HubAddress   string `json:"hubAddress"`
	Token        string `json:"token"`
	IgnoreCert   bool   `json:"ignoreCert"`
	Username     string `json:"username"`
	ID           string `json:"id"`
}

// WriteToken - Writes a copy of the token to the
func (c *Client) WriteToken() error {

	if c.Token == "" {
		return fmt.Errorf("Not logged in, or no UCP token present")
	}

	// build path
	path := fmt.Sprintf("%s/.storetoken", os.Getenv("HOME"))
	log.Debugf("Writing Token to [%s]", path)

	clientToken := internal{
		StoreAddress: c.STOREURL,
		HubAddress:   c.HUBURL,
		Username:     c.Username,
		Token:        c.Token,
		IgnoreCert:   c.IgnoreCert,
		ID:           c.ID,
	}

	b, err := json.Marshal(clientToken)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadToken - Reads the token from a file
func ReadToken() (*Client, error) {
	// build path
	path := fmt.Sprintf("%s/.storetoken", os.Getenv("HOME"))
	log.Debugf("Reading Token from [%s]", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	clientToken := internal{}

	err = json.Unmarshal(data, &clientToken)
	if err != nil {
		return nil, err
	}

	client := &Client{
		STOREURL:   clientToken.StoreAddress,
		HUBURL:     clientToken.HubAddress,
		Username:   clientToken.Username,
		Token:      clientToken.Token,
		IgnoreCert: clientToken.IgnoreCert,
		ID:         clientToken.ID,
	}
	if client.ID == "" {
		return client, fmt.Errorf("No User ID found in token")
	}
	if client.Token == "" {
		return client, fmt.Errorf("No User Token found in token")
	}
	return client, nil
}
