package messaging

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Register(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.POST("/api/towershare/message", func(c echo.Context) error {
			msg, err := parseIncoming(c)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}

			r := handleMessage(app, msg)

			// Twilio expects TwiML; everything else gets JSON with optional actions.
			if c.Request().Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
				return c.XML(http.StatusOK, TwiMLResponse{Message: r.Reply})
			}
			return c.JSON(http.StatusOK, r)
		})
		return nil
	})
}
