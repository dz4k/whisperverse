package model

import (
	"html/template"

	"github.com/benpate/compare"
	"github.com/benpate/data/option"
	"github.com/benpate/schema"
)

// Template represents an HTML template used for rendering Streams
type Template struct {
	TemplateID         string            `path:"templateId"         json:"templateId"         bson:"templateId"`         // Internal name/token other objects (like streams) will use to reference this Template.
	Label              string            `path:"label"              json:"label"              bson:"label"`              // Human-readable label used in management UI.
	Description        string            `path:"description"        json:"description"        bson:"description"`        // Human-readable long-description text used in management UI.
	Category           string            `path:"category"           json:"category"           bson:"category"`           // Human-readable category (grouping) used in management UI.
	Icon               string            `path:"icon"               json:"icon"               bson:"icon"`               // Icon image used in management UI.
	ContainedBy        []string          `path:"containedBy"        json:"containedBy"        bson:"containedBy"`        // Slice of Templates that can contain Streams that use this Template.
	ChildSortType      string            `path:"childSortType"      json:"childSortType"      bson:"childSortType"`      // SortType used to display children
	ChildSortDirection string            `path:"childSortDirection" json:"childSortDirection" bson:"childSortDirection"` // Sort direction "asc" or "desc" (Default is ascending)
	URL                string            `path:"url"                json:"url"                bson:"url"`                // URL where this template is published
	Schema             schema.Schema     `path:"schema"             json:"schema"             bson:"schema"`             // JSON Schema that describes the data required to populate this Template.
	States             map[string]State  `path:"states"             json:"states"             bson:"states"`             // Map of States (by state.ID) that Streams of this Template can be in.
	Roles              map[string]Role   `path:"roles"              json:"roles"              bson:"roles"`              // Map of custom roles defined by this Template.
	Actions            map[string]Action `path:"actions"            json:"actions"            bson:"actions"`            // Map of actions that can be performed on streams of this Template
	HTMLTemplate       *template.Template
}

// NewTemplate creates a new, fully initialized Template object
func NewTemplate(templateID string, funcMap template.FuncMap) Template {

	return Template{
		TemplateID:         templateID,
		ContainedBy:        make([]string, 0),
		ChildSortType:      "rank",
		ChildSortDirection: option.SortDirectionAscending,
		States:             make(map[string]State),
		Roles:              make(map[string]Role),
		Actions:            make(map[string]Action),
		HTMLTemplate:       template.New("").Funcs(funcMap),
	}
}

// CanBeContainedBy returns TRUE if this Streams using this Template can be nested inside of
// Streams using the Template named in the parameters
func (template *Template) CanBeContainedBy(templateName string) bool {
	return compare.Contains(template.ContainedBy, templateName)
}

// State searches for the State in this Template that matches the provided StateID
// If found, it is returned along with a TRUE
// If not found, an empty state is returned along with a FALSE
func (template *Template) State(stateID string) (State, bool) {
	state, ok := template.States[stateID]
	return state, ok
}

// Action returns the action object for a specified name
func (template *Template) Action(actionID string) *Action {

	if action, ok := template.Actions[actionID]; ok {
		return &action
	}

	return nil
}

// Validate runs any post-processing required after a Template is parsed by the TemplateService
func (template *Template) Validate() {

	for actionID, action := range template.Actions {
		action.ActionID = actionID
		action.Validate()
		template.Actions[actionID] = action
	}
}
