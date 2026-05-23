// Package example_module shows the pattern for extending PocketBase.
// Copy this package, rename it, and wire it into main.go's import block.
package example_module

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// Register wires the module into the PocketBase app.
func Register(app *pocketbase.PocketBase) {
	// Hook: run logic after every user record is created.
	app.OnRecordAfterCreateRequest("users").Add(func(e *core.RecordCreateEvent) error {
		app.Logger().Info("new user created", "id", e.Record.Id)
		return nil
	})

	// Custom REST route: GET /api/towershare/ping
	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.GET("/api/towershare/ping", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "pong from example_module"})
		})
		return nil
	})
}
