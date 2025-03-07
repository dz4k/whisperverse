package model

// Role is used in a map[sring]Role within each Template.  Role IDs are used to
// identify what actions a User can take on a Stream (given the user's Groups and the Stream's Template)
type Role struct {
	RoleID      string `path:"roleId"      json:"roleId"      bson:"roleId"`      // Unique ID for this role
	Label       string `path:"label"       json:"label"       bson:"label"`       // Short, human-friendly label used to select this role in UX
	Description string `path:"description" json:"description" bson:"description"` // Medium-length, human-friendly description that gives more details about this role
}
