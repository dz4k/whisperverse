package model

import (
	"github.com/benpate/data/journal"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubscriptionMethodRSS represents an RSS subscription
const SubscriptionMethodRSS = "RSS"

// SubscriptionMethodWebSub represents a WebSub subscription
const SubscriptionMethodWebSub = "WEBSUB"

type Subscription struct {
	SubscriptionID  primitive.ObjectID `path:"subscriptionId" json:"subscriptionId" bson:"_id"`            // Unique Identifier of this record
	ParentStreamID  primitive.ObjectID `path:"parentStreamId" json:"parentStreamId" bson:"parentStreamId"` // ID of the stream that owns this subscription
	Method          string             `path:"method"         json:"method"         bson:"method"`         // Method used to subscribe to remote streams (RSS, etc)
	URL             string             `path:"url"            json:"url"            bson:"url"`            // Connection URL for obtaining new sub-streams.
	LastPolled      int64              `path:"lastPolled"     json:"lastPolled"     bson:"lastPolled"`     // Unix Timestamp of the last date that this resource was retrieved.
	PollDuration    int                `path:"pollDuration"   json:"pollDuration"   bson:"pollDuration"`   // Time (in minutes) to wait between polling this resource.
	journal.Journal `json:"-" bson:"journal"`
}

func NewSubscription() Subscription {
	return Subscription{}
}

/*******************************************
 * DATA.OBJECT INTERFACE
 *******************************************/

// ID returns the primary key of this object
func (sub *Subscription) ID() string {
	return sub.SubscriptionID.Hex()
}
