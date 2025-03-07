package route

import (
	"net/http"

	"github.com/benpate/compare"
	"github.com/benpate/derp"
	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/whisperverse/whisperverse/handler"
	"github.com/whisperverse/whisperverse/middleware"
	"github.com/whisperverse/whisperverse/server"
)

// New returns all of the routes required for this application
func New(factory *server.Factory) *echo.Echo {

	e := echo.New()

	// echo Configuration
	e.HideBanner = true

	// Middleware
	domain := middleware.Domain(factory)

	// Well-Known API calls
	// https://en.wikipedia.org/wiki/List_of_/.well-known/_services_offered_by_webservers

	e.GET("/favicon.ico", handler.GetFavicon(factory))
	e.GET("/.well-known/nodeinfo", handler.GetNodeInfo(factory), domain)
	e.GET("/.well-known/oembed", handler.GetOEmbed(factory), domain)
	e.GET("/.well-known/webfinger", handler.GetWebfinger(factory), domain)
	e.GET("/.well-known/webmention", handler.PostWebMention(factory), domain)

	// Local links for static resources
	e.Static("/static", factory.StaticPath())

	// RSS Feed
	e.GET("/feed.json", handler.GetRSS(factory), domain)

	// Authentication Pages
	e.GET("/signin", handler.GetSignIn(factory), domain)
	e.POST("/signin", handler.PostSignIn(factory), domain)
	e.POST("/signout", handler.PostSignOut(factory), domain)

	// STREAM PAGES
	e.GET("/", handler.GetStream(factory), domain)
	e.GET("/:stream", handler.GetStream(factory), domain)
	e.GET("/:stream/:action", handler.GetStream(factory), domain)
	e.POST("/:stream/:action", handler.PostStream(factory), domain)
	e.DELETE("/:stream", handler.PostStream(factory), domain)

	e.GET("/subscriptions", handler.ListSubscriptions(factory))
	e.GET("/subscriptions/:subscriptionId", handler.GetSubscription(factory))
	e.POST("/subscriptions/:subscriptionId", handler.PostSubscription(factory))
	e.DELETE("/subscriptions/:subscriptionId", handler.DeleteSubscription(factory))

	// TODO: Can Attachments and SSE be moved into a custom render step?

	// DOMAIN ADMIN PAGES
	e.GET("/admin", handler.GetAdmin(factory), domain)
	e.GET("/admin/:param1", handler.GetAdmin(factory), domain)
	e.POST("/admin/:param1", handler.PostAdmin(factory), domain)
	e.GET("/admin/:param1/:param2", handler.GetAdmin(factory), domain)
	e.POST("/admin/:param1/:param2", handler.PostAdmin(factory), domain)
	e.GET("/admin/:param1/:param2/:param3", handler.GetAdmin(factory), domain)
	e.POST("/admin/:param1/:param2/:param3", handler.PostAdmin(factory), domain)

	// Hard-coded routes for additional stream services
	e.GET("/:stream/attachments/:attachment", handler.GetAttachment(factory), domain)
	e.GET("/:stream/sse", handler.ServerSentEvent(factory), domain)

	// Profile Pages / ActivityPub
	e.GET("/inbox", handler.GetProfileInbox(factory), domain)

	/*
		e.GET("/outbox", handler.GetProfileOutbox(factory), domain)
		e.GET("/people/:userId", handler.GetProfile(factory), domain)
		e.GET("/people/:userId/inbox", handler.GetSocialInbox(factory), domain)
		e.POST("/people/:userId/inbox", handler.PostSocialInbox(factory), domain)
		e.GET("/people/:userId/outbox", handler.GetSocialOutbox(factory), domain)
		e.POST("/people/:userId/outbox", handler.PostSocialOutbox(factory), domain)
		e.GET("/people/:userId/followers", handler.GetSocialFollowers(factory), domain)
		e.GET("/people/:userId/following", handler.GetSocialFollowing(factory), domain)
		e.GET("/people/:userId/liked", handler.GetSocialLiked(factory), domain)

		// PROFILE PAGES
		// e.GET("/me/", handler.GetProfile(factory))
		// e.POST("/me", handler.PostProfile(factory))
		// e.GET("/me/:action", handler.PostProfile(factory))
		// e.POST("/me/:action", handler.PostProfile(factory))
	*/

	// SERVER ADMIN PAGES (dynamic URLs help discourage 4337 H4XX0RZ)
	if serverAdminURL := factory.AdminURL(); serverAdminURL != "" {
		server := e.Group(serverAdminURL, middleware.ServerAdmin(factory))
		server.GET("", handler.GetServerIndex(factory))
		server.POST("", handler.GetServerIndex(factory))
		server.GET("/:domain", handler.GetServerDomain(factory))
		server.POST("/:domain", handler.PostServerDomain(factory))
		server.DELETE("/:domain", handler.DeleteServerDomain(factory))
	}

	// SITE-WIDE ERROR HANDLER
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {

		// Special handling of permisssion errors
		code := derp.ErrorCode(err)
		switch code {
		case http.StatusUnauthorized, http.StatusForbidden:

			// If the server admin console is active for this server
			if serverAdminURL := factory.AdminURL(); serverAdminURL != "" {
				// ... and the requested page is in the admin section ...
				if compare.BeginsWith(ctx.Request().URL.Path, serverAdminURL) {
					// ... then redirect to the root of the admin section.
					ctx.Redirect(http.StatusTemporaryRedirect, serverAdminURL)
					return
				}
			}

			if ctx.Request().URL.Path != "/signin" {
				ctx.Redirect(http.StatusTemporaryRedirect, "/signin")
				return
			}
			ctx.String(code, derp.Message(err))
			return
		}

		// On localhost, allow developers to see full error dump.
		if ctx.Request().Host == "localhost" {
			ctx.String(derp.ErrorCode(err), spew.Sdump(err))
		}

		// Fall through to general error handler
		ctx.String(derp.ErrorCode(err), derp.Message(err))
	}

	return e
}
