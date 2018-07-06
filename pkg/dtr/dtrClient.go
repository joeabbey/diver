package dtr

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

// Client - Is the basic Client struct
type Client struct {
	Username   string
	Password   string
	DTRURL     string
	IgnoreCert bool
	Token      string
}

// NewBasicAuthClient - Creates a basic client to connecto the DTR
func NewBasicAuthClient(username, password, url string, ignoreCert bool) *Client {
	return &Client{
		Username:   username,
		Password:   password,
		DTRURL:     url,
		IgnoreCert: ignoreCert,
	}
}

// Connect - Will attempt to connect to DTR
func (c *Client) Connect() error {
	if c.Username == "" {
		return fmt.Errorf("DTR Username hasn't been entered")
	}

	if c.Password == "" {
		return fmt.Errorf("DTR Password is blank")
	}

	if c.DTRURL == "" {
		return fmt.Errorf("DTR URL hasn't been entered")
	}
	// Add the /auth/log to the URL
	url := fmt.Sprintf("%s/api/v0/api_tokens?username=%s", c.DTRURL, c.Username)
	log.Debugf("Created URL [%s]", url)

	// Create a default token label for DTR
	data := map[string]string{
		"tokenLabel": "diver",
	}
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	response, err := c.postRequest(url, b)
	if err != nil {
		log.Debugf("%v", response)
		return err
	}

	var responseData map[string]interface{}

	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return err
	}

	if responseData["token"] != nil {
		c.Token = responseData["token"].(string)
	} else {
		log.Debugf("%s", string(response))
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
	req.SetBasicAuth(c.Username, c.Password)
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
	}

	req.Header.Add("Content-Type", "application/json")

	bytes, err := c.doRequest(req)
	if err != nil {
		return bytes, err
	}
	return bytes, nil
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
	DTRAddress string `json:"address"`
	Token      string `json:"token"`
	IgnoreCert bool   `json:"ignoreCert"`
}

// WriteToken - Writes a copy of the token to the
func (c *Client) WriteToken() error {

	if c.Token == "" {
		return fmt.Errorf("Not logged in, or no DTR token present")
	}

	// build path
	path := fmt.Sprintf("%s/.dtrtoken", os.Getenv("HOME"))
	log.Debugf("Writing Token to [%s]", path)

	clientToken := internal{
		DTRAddress: c.DTRURL,
		Token:      c.Token,
		IgnoreCert: c.IgnoreCert,
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
	path := fmt.Sprintf("%s/.dtrtoken", os.Getenv("HOME"))
	log.Debugf("Reading Token from [%s]", path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("No Session Token could be found, please login")
	}

	clientToken := internal{}

	err = json.Unmarshal(data, &clientToken)
	if err != nil {
		return nil, fmt.Errorf("Corrupted Session Token, please login")
	}

	client := &Client{
		DTRURL:     clientToken.DTRAddress,
		Token:      clientToken.Token,
		IgnoreCert: clientToken.IgnoreCert,
	}

	return client, nil
}
