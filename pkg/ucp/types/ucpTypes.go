package ucptypes

import (
	"encoding/json"
	"time"
)

// Account - Is the basic Account struct
type Account struct {
	FullName   string `json:"fullName"`
	ID         string `json:"id"`
	IsActive   bool   `json:"isActive"`
	IsAdmin    bool   `json:"isAdmin"`
	IsOrg      bool   `json:"isOrg"`
	Name       string `json:"name"`
	Password   string `json:"password"`
	SearchLDAP bool   `json:"searchLDAP"`
}

// AccountList - The format returned by a query of accounts
type AccountList struct {
	Accounts []Account `json:"accounts"`
}

// Team - is the structure for defining a team
type Team struct {
	Description  string `json:"description"`
	ID           string `json:"id"`
	MembersCount int    `json:"membersCount"`
	Name         string `json:"name"`
	OrgID        string `json:"orgID"`
}

// Teams is returned from querying an organisation
type Teams struct {
	NextPage      string `json:"nextPageStart"`
	ResourceCount int    `json:"resourceCount"`
	Teams         []Team `json:"teams"`
}

// Roles are the structure the defines a role
type Roles struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	SystemRole bool            `json:"system_role"`
	Operations json.RawMessage // Captures the raw output of the remaining json object
}

// A grant is based upon three keys:
// -- ObjectID == Collection
// -- RoleID == Links the role that is applied (rights)
// -- SubjectID == User that has is linked to the collection with the appropriate rights

// Grant - the the three elements needed for a grant
type Grant struct {
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

// Collection - An array of JSON Structs that detail the collections in Docker UCP
type Collection struct {
	CreatedAt        time.Time                    `json:"created_at"`
	ID               string                       `json:"id"`
	LabelConstraints []CollectionLabelConstraints `json:"label_constraints"`
	Legacylabelkey   string                       `json:"legacylabelkey"`
	Legacylabelvalue string                       `json:"legacylabelvalue"`
	Name             string                       `json:"name"`
	ParentIds        []string                     `json:"parent_ids"`
	Path             string                       `json:"path"`
	UpdatedAt        time.Time                    `json:"updated_at"`
}

// CollectionLabelConstraints - defines constraints to be applied to a collection
type CollectionLabelConstraints struct {
	Equality   bool   `json:"equality"`
	LabelKey   string `json:"label_key"`
	LabelValue string `json:"label_value"`
	Type       string `json:"type"`
}

// InterlockConfig - configuration for Interlock layer 7 routing
type InterlockConfig struct {
	HTTPPort         int    `json:"HTTPPort"`
	HTTPSPort        int    `json:"HTTPSPort"`
	Arch             string `json:"Arch"`
	InterlockEnabled bool   `json:"InterlockEnabled"`
}

// ClientBundles - defines a client bundle list response
type ClientBundles struct {
	AccountPublicKeys []struct {
		ID           string `json:"id"`
		AccountID    string `json:"accountID"`
		PublicKey    string `json:"publicKey"`
		Label        string `json:"label"`
		Certificates []struct {
			Label string `json:"label"`
			Cert  string `json:"cert"`
		} `json:"certificates"`
	} `json:"accountPublicKeys"`
}

// ContainerProcesses - Struct that defines all Processes in a container
type ContainerProcesses struct {
	Processes [][]string `json:"Processes"`
	Titles    []string   `json:"Titles"`
}

// ServiceConfig - Used to interact with configurations applied to services
type ServiceConfig struct {
	ID      string `json:"ID"`
	Version struct {
		Index int `json:"Index"`
	} `json:"Version"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
	Spec      struct {
		Name   string `json:"Name"`
		Labels struct {
			Label []map[string]string
		} `json:"Labels"`
		Data string `json:"Data"`
	} `json:"Spec"`
}
