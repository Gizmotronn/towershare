// Package simulator serves a local chat UI at /simulator.
// It is only registered when APP_ENV != "production".
package simulator

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

//go:embed chat.html
var chatHTML []byte

func Register(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/simulator", func(c echo.Context) error {
			return c.HTMLBlob(http.StatusOK, chatHTML)
		})
		return nil
	})
}
