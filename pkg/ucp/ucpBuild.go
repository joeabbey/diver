package ucp

import (
	"fmt"
	"html"

	log "github.com/sirupsen/logrus"
)

// BuildPlan - This collection of bools builds the output from a service query
type BuildPlan struct {
	// The name of the service to query
	GHURL     string
	Tag       string
	BuildHost string
}

//BuildImage - This will return a list of services
func (c *Client) BuildImage(b *BuildPlan) error {

	//TODO - RM of images.

	// Add the /build?remote=%s to the URL

	url := fmt.Sprintf("%s/build?remote=%s", c.UCPURL, html.EscapeString(b.GHURL))

	if b.Tag != "" {
		url = fmt.Sprintf("%s&t=%s", url, html.EscapeString(b.Tag))
	}

	if b.BuildHost != "" {
		// TODO - this is a hack as html.escapestring() wont escape "{}:"
		beginEncode := "%7B%22constraint:node%22%3A%22"
		endEncode := "%22%7D"
		encodeString := beginEncode + b.BuildHost + endEncode
		url = fmt.Sprintf("%s&buildargs=%s", url, encodeString)
	}

	log.Debugf("Built url [%s]", url)
	err := c.postRequestStream(url, []byte("test"))
	if err != nil {
		return err
	}
	return nil
}
