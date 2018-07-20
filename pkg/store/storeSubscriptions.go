package store

import (
	"encoding/json"
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
)

// Subscription - this is the returned struct from a store query
type Subscription struct {
	Name               string    `json:"name"`
	SubscriptionID     string    `json:"subscription_id"`
	DockerID           string    `json:"docker_id"`
	ProductID          string    `json:"product_id"`
	CreatedByDockerID  string    `json:"created_by_docker_id"`
	ProductRatePlan    string    `json:"product_rate_plan"`
	ProductRatePlanID  string    `json:"product_rate_plan_id"`
	InitialPeriodStart time.Time `json:"initial_period_start"`
	CurrentPeriodStart time.Time `json:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end"`
	State              string    `json:"state"`
	Eusa               struct {
		Accepted   bool      `json:"accepted"`
		AcceptedBy string    `json:"accepted_by"`
		AcceptedOn time.Time `json:"accepted_on"`
	} `json:"eusa"`
	PricingComponents []struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	} `json:"pricing_components"`
	MarketingOptIn bool `json:"marketing_opt_in"`
}

//GetAllSubscriptions - Retrieves all subscriptions
func (c *Client) GetAllSubscriptions(id string) ([]Subscription, error) {
	log.Debugf("Retrieving all subscriptions")

	if id == "" {
		log.Debugf("Attempting to read ID from ~/.dockerstore")
		id = c.ID
	}

	url := fmt.Sprintf("%s/?docker_id=%s", c.HUBURL, id)
	log.Debugf("Url = %s", url)
	response, err := c.getRequest(url, nil)
	if err != nil {
		return nil, err
	}

	var returnedSubs []Subscription

	err = json.Unmarshal(response, &returnedSubs)
	if err != nil {
		return nil, err
	}

	return returnedSubs, nil
}
