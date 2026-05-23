package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	// Ensure migration init() funcs run.
	_ "github.com/gizmotronn/towershare/pb_migrations"

	"github.com/gizmotronn/towershare/pb_modules/example_module"
	"github.com/gizmotronn/towershare/pb_modules/messaging"
	"github.com/gizmotronn/towershare/pb_modules/simulator"
)

func main() {
	app := pocketbase.New()

	// Auto-migrate in development; in prod migrations run explicitly.
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isDevMode(),
	})

	// Register modules.
	example_module.Register(app)
	messaging.Register(app)
	if isDevMode() {
		simulator.Register(app)
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func isDevMode() bool {
	return os.Getenv("APP_ENV") != "production"
}
