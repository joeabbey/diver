package ucp

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
)

// Client - Is the basic Client struct
type Client struct {
	Username   string
	Password   string
	UCPURL     string
	IgnoreCert bool
	Token      string
}

// NewBasicAuthClient - Creates a basic client to connecto the UCP
func NewBasicAuthClient(username, password, url string, ignoreCert bool) *Client {
	return &Client{
		Username:   username,
		Password:   password,
		UCPURL:     url,
		IgnoreCert: ignoreCert,
	}
}

// Connect - Will attempt to connect to UCP
func (c *Client) Connect() error {
	if c.Username == "" {
		return fmt.Errorf("UCP Username hasn't been entered")
	}

	if c.Password == "" {
		return fmt.Errorf("UCP Password is blank")
	}

	if c.UCPURL == "" {
		return fmt.Errorf("UCP URL hasn't been entered")
	}
	// Add the /auth/log to the URL
	url := fmt.Sprintf("%s/auth/login", c.UCPURL)

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
	if responseData["auth_token"] != nil {
		c.Token = responseData["auth_token"].(string)
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

	switch resp.StatusCode {
	case 200:
		log.Debug("[success] HTTP Status code 200")
	case 201:
		log.Debug("[success] HTTP Status code 201")
	case 204:
		log.Debug("[success] HTTP Status code 204")
	case 400:
		log.Debugf("HTTP Error code: %d for URL: %s", resp.StatusCode, req.URL.String())
		// Return Body, can be processed with ucp.ParseURL elsewhere
		return nil, fmt.Errorf("%s", body)
	case 401:
		log.Debugf("HTTP Error code: %d for URL: %s", resp.StatusCode, req.URL.String())
		// Return Body, can be processed with ucp.ParseURL elsewhere
		return nil, fmt.Errorf("%s", body)
	case 404:
		log.Debugf("HTTP Error code: %d for URL: %s", resp.StatusCode, req.URL.String())
		// Return Body, can be processed with ucp.ParseURL elsewhere
		return nil, fmt.Errorf("%s", body)
	default:
		log.Debugf("[Untrapped return code] %d", resp.StatusCode)
	}
	return body, nil
}

type internal struct {
	UCPAddress string `json:"address"`
	Token      string `json:"token"`
	IgnoreCert bool   `json:"ignoreCert"`
}

// WriteToken - Writes a copy of the token to the
func (c *Client) WriteToken() error {

	if c.Token == "" {
		return fmt.Errorf("Not logged in, or no UCP token present")
	}

	// build path
	path := fmt.Sprintf("%s/.ucptoken", os.Getenv("HOME"))
	log.Debugf("Writing Token to [%s]", path)

	clientToken := internal{
		UCPAddress: c.UCPURL,
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
	path := fmt.Sprintf("%s/.ucptoken", os.Getenv("HOME"))
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
		UCPURL:     clientToken.UCPAddress,
		Token:      clientToken.Token,
		IgnoreCert: clientToken.IgnoreCert,
	}

	return client, nil
}
