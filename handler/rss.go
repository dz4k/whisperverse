package handler

import (
	"github.com/benpate/derp"
	"github.com/benpate/ghost/service"
	"github.com/labstack/echo/v4"
)

// GetRSS returns an RSS data feed for the requested URL
func GetRSS(fm service.FactoryMaker) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		factory := fm.Factory(ctx.Request().Context())

		service := factory.RSS()

		feed, err := service.Feed()

		if err != nil {
			return derp.Wrap(err, "handler.GetRSS", "Error generating RSS feed").Report()
		}

		// TODO: Replace these with real values from the server setup.
		feed.Title = "Title Goes Here"
		feed.Description = "Description Goes Here"
		feed.FeedUrl = "Feed URL goes here"

		result, errr := feed.ToJSON()

		if errr != nil {
			return derp.New(500, "handler.GetRSS", "Error writing JSON feed information", errr).Report()
		}

		response := ctx.Response()
		response.Header().Set("content-type", "application/json")

		return ctx.String(200, result)
	}
}
