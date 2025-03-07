package service

import (
	"fmt"
	"time"

	"github.com/benpate/data"
	"github.com/benpate/data/option"
	"github.com/benpate/datatype"
	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/benpate/list"
	"github.com/benpate/nebula"
	"github.com/mmcdole/gofeed"
	"github.com/whisperverse/whisperverse/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Subscription manages all interactions with the Subscription collection
type Subscription struct {
	collection     data.Collection
	streamService  *Stream
	contentLibrary *nebula.Library
}

// NewSubscription returns a fully populated Subscription service.
func NewSubscription(collection data.Collection, streamService *Stream, contentLibrary *nebula.Library) *Subscription {

	result := Subscription{
		collection:     collection,
		streamService:  streamService,
		contentLibrary: contentLibrary,
	}

	go result.start()

	return &result
}

// New creates a newly initialized Subscription that is ready to use
func (service *Subscription) New() model.Subscription {
	return model.NewSubscription()
}

func (service *Subscription) start() {

	ticker := time.NewTicker(20 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C
		fmt.Println(".. Polling Subscriptions")
		it, err := service.ListPollable()

		if err != nil {
			derp.Report(derp.Wrap(err, "whisper.service.Subscription.Run", "Error listing pollable subscriptions"))
			continue
		}

		subscription := model.Subscription{}

		for it.Next(&subscription) {
			service.pollSubscription(&subscription)
			subscription = model.Subscription{}
		}
	}
}

func (service *Subscription) pollSubscription(sub *model.Subscription) {
	// TODO: Check if subscription is past its polling window

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(sub.URL)

	if err != nil {
		derp.Report(derp.Wrap(err, "whisper.service.Subscription.Poll", "Error Parsing Feed URL"))
		return
	}

	for _, item := range feed.Items {
		if err := service.updateStream(sub, item); err != nil {
			derp.Report(derp.Wrap(err, "whisper.service.Subscription.Poll", "Error updating local stream"))
		}
	}
}

func (service *Subscription) updateStream(sub *model.Subscription, item *gofeed.Item) error {

	stream := model.NewStream()

	err := service.streamService.LoadBySource(sub.ParentStreamID, item.Link, &stream)

	if err != nil {

		// Anything but a "not found" error is a real error
		if !derp.NotFound(err) {
			return derp.Wrap(err, "whisper.service.Subscription.Poll", "Error loading local stream")
		}

		// Fall through means "not found" which means "make a new stream"
		stream = model.NewStream()
		stream.TemplateID = "rss-article"
		stream.ParentID = sub.ParentStreamID
		stream.SourceURL = item.Link
		stream.StateID = "unread"
	}

	updateDate := item.PublishedParsed.Unix()

	if item.UpdatedParsed != nil {
		updateDate = item.UpdatedParsed.Unix()
	}

	// If stream has been updated since previous save, then set new values
	if stream.SourceUpdated > updateDate {

		// Populate header information into the stream
		stream.Label = item.Title
		stream.Description = item.Description
		stream.PublishDate = item.PublishedParsed.Unix()
		stream.SourceUpdated = updateDate

		// Populate content into a nebula container
		stream.Content = nebula.NewContainer()
		stream.Content.NewItemWithInit(service.contentLibrary, nebula.ItemTypeHTML, datatype.Map{
			"html": item.Content,
		})

		if item.Author == nil {
			stream.AuthorName = ""
			// stream.AuthorEmail = ""
		} else {
			stream.AuthorName = item.Author.Name
			// stream.AuthorEmail = item.Author.Email
		}

		if item.Image != nil {
			stream.ThumbnailImage = item.Image.URL
		} else {
			stream.ThumbnailImage = ""

			// Search for an image in the enclosures
			for _, enclosure := range item.Enclosures {
				if list.Head(enclosure.Type, "/") == "image" {
					stream.ThumbnailImage = enclosure.URL
					break
				}
			}
		}

		if err := service.streamService.Save(&stream, "Imported from RSS feed"); err != nil {
			return derp.Wrap(err, "whisper.service.Subscription.Poll", "Error saving stream")
		}
	}

	return nil
}

/*******************************************
 * COMMON DATA FUNCTIONS
 *******************************************/

// List returns an iterator containing all of the Subscriptions who match the provided criteria
func (service *Subscription) List(criteria exp.Expression, options ...option.Option) (data.Iterator, error) {
	return service.collection.List(notDeleted(criteria), options...)
}

// Load retrieves an Subscription from the database
func (service *Subscription) Load(criteria exp.Expression, result *model.Subscription) error {

	if err := service.collection.Load(notDeleted(criteria), result); err != nil {
		return derp.Wrap(err, "whisper.service.Subscription", "Error loading Subscription", criteria)
	}

	return nil
}

// Save adds/updates an Subscription in the database
func (service *Subscription) Save(subscription *model.Subscription, note string) error {

	if err := service.collection.Save(subscription, note); err != nil {
		return derp.Wrap(err, "whisper.service.Subscription", "Error saving Subscription", subscription, note)
	}

	return nil
}

// Delete removes an Subscription from the database (virtual delete)
func (service *Subscription) Delete(subscription *model.Subscription, note string) error {

	if err := service.collection.Delete(subscription, note); err != nil {
		return derp.Wrap(err, "whisper.service.Subscription", "Error deleting Subscription", subscription, note)
	}

	return nil
}

/*******************************************
 * CUSTOM QUERIES
 *******************************************/

func (service *Subscription) ListPollable() (data.Iterator, error) {
	pollDuration := time.Now().Add(-1 * time.Hour).Unix()
	criteria := exp.LessThan("lastPolled", pollDuration)

	return service.List(criteria, option.SortAsc("lastPolled"))
}

func (service *Subscription) ListByUserID(userID primitive.ObjectID) (data.Iterator, error) {
	criteria := exp.Equal("userId", userID)
	return service.List(criteria, option.SortAsc("lastPolled"))
}

func (service *Subscription) LoadByID(subscriptionID primitive.ObjectID, result *model.Subscription) error {

	criteria := exp.Equal("_id", subscriptionID)

	if err := service.Load(criteria, result); err != nil {
		return derp.Wrap(err, "service.Subscription.LoadByID", "Error loading Subscription", criteria)
	}

	return nil
}
