package model

import (
	"github.com/benpate/data/journal"
	"github.com/benpate/schema"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Template represents an HTML template to be used for generating an HTML page.
type Template struct {
	TemplateID  primitive.ObjectID    `json:"templateId"  bson:"_id"`         // Unique Identifier for this Template.
	URL         string                `json:"url"         bson:"url"`         // URL where this template is published
	Category    string                `json:"category"    bson:"category"`    // Human-readable category (grouping) used in management UI.
	Label       string                `json:"label"       bson:"label"`       // Human-readable label used in management UI.
	IconURL     string                `json:"iconUrl"     bson:"iconUrl"`     // Icon image used in management UI.
	Format      string                `json:"format"      bson:"format"`      // TOKEN that other templates will use to reference this template.
	Schema      schema.Schema         `json:"schema"      bson:"schema"`      // JSON Schema that describes the data required to populate this Template.
	Views       map[string]View       `json:"views"       bson:"views"`       // Map of Views (by view.ID) that are available to Streams of this Template.
	States      map[string]State      `json:"states"      bson:"states"`      // Map of States (by state.ID) that Streams of this Template can be in.
	Transitions map[string]Transition `json:"transitions" bson:"transitions"` // Map of Transitions (by transition.ID) between States of this Template.

	journal.Journal
}

// ID returns the primary key of this object
func (t *Template) ID() string {
	return t.TemplateID.Hex()
}
