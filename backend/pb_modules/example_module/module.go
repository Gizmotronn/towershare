// Package example_module shows the pattern for extending PocketBase.
// Copy this package, rename it, and wire it into main.go's import block.
package example_module

import (
	"net/http"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

// Register wires the module into the PocketBase app.
// Call this from main.go, or use a blank import + init() pattern.
func Register(app *pocketbase.PocketBase) {
	// Hook: run logic after every user record is created.
	app.OnRecordAfterCreateSuccess("users").BindFunc(func(e *core.RecordEvent) error {
		app.Logger().Info("new user created", "id", e.Record.Id)
		return e.Next()
	})

	// Custom REST route: GET /api/towershare/ping
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		se.Router.GET("/api/towershare/ping", func(re *core.RequestEvent) error {
			return re.JSON(http.StatusOK, map[string]string{"message": "pong from example_module"})
		})
		return se.Next()
	})
}
