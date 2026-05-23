// Package ui serves the TowerShare web app at /app.
package ui

import (
	_ "embed"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

//go:embed app.html
var appHTML []byte

func Register(app *pocketbase.PocketBase) {
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/app", func(c echo.Context) error {
			return c.HTMLBlob(http.StatusOK, appHTML)
		})
		return nil
	})
}
