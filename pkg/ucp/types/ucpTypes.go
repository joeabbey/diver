package ucptypes

import "encoding/json"

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
	Description string `json:"description"`
	Name        string `json:"name"`
}

type Roles struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	ServiceRole bool            `json:"servicerole"`
	Operations  json.RawMessage // Captures the raw output of the remaining json object
}

// TODO fix the rest of the struct
type Collection struct {
	// "name": "Private",
	// "path": "/Shared/Private",
	// "id": "private",
	// "parent_ids": [
	//   "root",
	//   "swarm",
	//   "shared"
	// ],
	// "label_constraints": [],
	// "legacylabelkey": "",
	// "legacylabelvalue": "",
	// "created_at": "2018-06-11T17:16:14.124Z",
	// "updated_at": "2018-06-11T17:16:14.124Z"
	Name string `json:"name"`
	Path string `json:"path"`
	ID   string `json:"id"`
}

// A grant is based upon three keys:
// -- ObjectID == Collection
// -- RoleID == Links the role that is applied (rights)
// -- SubjectID == User that has is linked to the collection with the appropriate rights

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
