package main

import (
	"log"
	"os"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"

	// Register your modules here:
	_ "github.com/gizmotronn/towershare/pb_modules/example_module"
)

func main() {
	app := pocketbase.New()

	// Auto-migrate in development; in prod migrations run explicitly.
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		Automigrate: isDevMode(),
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func isDevMode() bool {
	return os.Getenv("APP_ENV") != "production"
}
