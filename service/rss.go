package service

import (
	"time"

	"github.com/benpate/data/option"
	"github.com/benpate/derp"
	"github.com/benpate/exp"
	"github.com/gorilla/feeds"
	"github.com/whisperverse/whisperverse/model"
)

// RSS service generates RSS feeds of the available streams in the database
type RSS struct {
	streamService *Stream
}

// NewRSS returns a fully initialized RSS service
func NewRSS(streamService *Stream) *RSS {
	return &RSS{
		streamService: streamService,
	}
}

// Feed generates an RSS data feed based on the provided query criteria.  This feed
// has a lot of incomplete data at the top level, so we're expecting the handler
// that calls this to fill in the rest of the gaps before it passes the values back
// to the requester.
func (rss RSS) Feed(criteria ...exp.Expression) (*feeds.JSONFeed, error) {

	filter := exp.And(criteria...)

	streams, err := rss.streamService.List(filter, option.SortDesc("publishDate"))
	stream := model.NewStream()

	if err != nil {
		return nil, derp.Wrap(err, "service.rss.Feed", "Error loading streams")
	}

	result := feeds.JSONFeed{
		Items: []*feeds.JSONItem{},
	}

	for streams.Next(&stream) {
		result.Items = append(result.Items, rss.Item(stream))
	}

	return &result, nil
}

// Item converts a single model.Stream into a feeds.JSONItem
func (rss RSS) Item(stream model.Stream) *feeds.JSONItem {

	publishDate := time.Unix(stream.PublishDate, 0)
	modifiedDate := time.Unix(stream.Journal.UpdateDate, 0)

	return &feeds.JSONItem{
		Id: "",
		// Url:           stream.URL,
		ExternalUrl:   stream.SourceURL,
		Title:         stream.Label,
		Summary:       stream.Description,
		Image:         stream.ThumbnailImage,
		PublishedDate: &publishDate,
		ModifiedDate:  &modifiedDate,

		Author: &feeds.JSONAuthor{
			Name: stream.AuthorName,
			Url:  stream.AuthorURL,
		},
	}
}
