package handler

import (
	"bytes"
	"net/http"

	"github.com/benpate/derp"
	"github.com/labstack/echo/v4"
	"github.com/whisperverse/whisperverse/server"
)

// GetSignIn generates an echo.HandlerFunc that handles GET /signin requests
func GetSignIn(factoryManager *server.Factory) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		var buffer bytes.Buffer

		factory, err := factoryManager.ByContext(ctx)

		if err != nil {
			return derp.Report(derp.Wrap(err, "whisper.handler.GetSignin", "Error getting factory"))
		}

		template := factory.Layout().Global().HTMLTemplate

		if err := template.ExecuteTemplate(&buffer, "signin", "error message goes here."); err != nil {
			return derp.Report(derp.Wrap(err, "whisper.handler.GetSignin", "Error executing template"))
		}

		return ctx.HTML(200, buffer.String())
	}
}

// PostSignIn generates an echo.HandlerFunc that handles POST /signin requests
func PostSignIn(factoryManager *server.Factory) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		factory, err := factoryManager.ByContext(ctx)

		if err != nil {
			return derp.New(500, "whisper.handler.PostSignIn", "Invalid Request.  Please try again later.")
		}

		s := factory.Steranko()

		if err := s.SignIn(ctx); err != nil {
			ctx.Response().Header().Add("HX-Trigger", "SigninError")
			return ctx.HTML(derp.ErrorCode(err), derp.Message(err))
		}

		ctx.Response().Header().Add("HX-Trigger", "SigninSuccess")

		return ctx.NoContent(200)
	}
}

// PostSignOut generates an echo.HandlerFunc that handles POST /signout requests
func PostSignOut(factoryManager *server.Factory) echo.HandlerFunc {

	return func(ctx echo.Context) error {

		factory, err := factoryManager.ByContext(ctx)

		if err != nil {
			return derp.New(500, "whisper.handler.PostSignOut", "Invalid Request.  Please try again later.")
		}

		s := factory.Steranko()

		if err := s.SignOut(ctx); err != nil {
			return derp.Wrap(err, "whisper.handler.PostSignOut", "Error Signing Out")
		}

		// Forward the user back to the home page of the website.
		ctx.Response().Header().Add("HX-Redirect", "/")
		return ctx.NoContent(http.StatusNoContent)
	}
}
