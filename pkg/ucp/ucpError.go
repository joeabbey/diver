package ucp

import (
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

// erro
type errorResponse struct {
	Errors []errorEntry `json:"errors"`
}

type errorEntry struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ParseUCPError - This will read through the return from UCP and report the error
func parseUCPError(response string) error {
	log.Debugf("%v", response)
	e := errorResponse{}

	// TODO
	// Assuming that all main UCP calls will respond using the above JSON structure (could fail silently if that isn't the case)
	err := json.Unmarshal([]byte(response), &e)
	if err != nil {
		return err
	}

	log.Errorf("%d reported error(s)", len(e.Errors))
	for i := range e.Errors {
		log.Errorf("Error: %s [%s]", e.Errors[i].Code, e.Errors[i].Message)
	}
	return nil
}
