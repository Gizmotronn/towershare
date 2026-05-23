package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// The `users` collection is created automatically by PocketBase.
		// Add any extra fields or custom collections here.
		// Example: extend users with a display_name field.
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		users.Fields.Add(&core.TextField{
			Name:     "display_name",
			Required: false,
		})

		return app.Save(users)
	}, func(app core.App) error {
		// Down migration: remove the display_name field.
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		users.Fields.RemoveById(users.Fields.GetByName("display_name").GetId())
		return app.Save(users)
	})
}
