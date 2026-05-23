// Package messaging handles inbound messages from WhatsApp (Twilio) and the
// local chat simulator. The same endpoint accepts both — Twilio sends
// application/x-www-form-urlencoded, the simulator sends JSON.
package messaging

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func Register(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		// POST /api/towershare/message
		// Accepts Twilio form payload OR {"from":"...","body":"..."} JSON.
		e.Router.POST("/api/towershare/message", func(c echo.Context) error {
			msg, err := parseIncoming(c)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}

			reply := handleMessage(app, msg)

			// Twilio expects TwiML; the simulator expects JSON.
			// Detect by Content-Type of the request.
			if c.Request().Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
				return c.XML(http.StatusOK, TwiMLResponse{Message: reply})
			}
			return c.JSON(http.StatusOK, map[string]string{
				"from":  "TowerShare Bot",
				"reply": reply,
			})
		})
		return nil
	})
}
