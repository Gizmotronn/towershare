package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		// --- Towers collection ---
		towers := core.NewBaseCollection("towers")
		towers.Fields.Add(
			&core.TextField{Name: "name", Required: true},
			&core.TextField{Name: "address", Required: true},
		)
		towers.Indexes = []string{
			"CREATE UNIQUE INDEX idx_towers_name ON towers (name)",
		}
		if err := app.Save(towers); err != nil {
			return err
		}

		// --- Apartments collection ---
		apartments := core.NewBaseCollection("apartments")
		apartments.Fields.Add(
			&core.TextField{Name: "unit_number", Required: true},
			&core.NumberField{Name: "floor", Required: false},
			&core.RelationField{
				Name:          "tower",
				Required:      true,
				CollectionId:  towers.Id,
				MaxSelect:     1,
				CascadeDelete: true,
			},
		)
		apartments.Indexes = []string{
			"CREATE UNIQUE INDEX idx_apartments_unit_tower ON apartments (unit_number, tower)",
		}
		if err := app.Save(apartments); err != nil {
			return err
		}

		// --- Extend users: apartment relation ---
		users, err := app.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		users.Fields.Add(
			&core.RelationField{
				Name:         "apartment",
				Required:     false,
				CollectionId: apartments.Id,
				MaxSelect:    1,
			},
		)
		if err := app.Save(users); err != nil {
			return err
		}

		// --- Tower managers (relation back to users) ---
		towersRefresh, err := app.FindCollectionByNameOrId("towers")
		if err != nil {
			return err
		}
		towersRefresh.Fields.Add(
			&core.RelationField{
				Name:         "managers",
				Required:     false,
				CollectionId: users.Id,
				MaxSelect:    100,
			},
		)
		if err := app.Save(towersRefresh); err != nil {
			return err
		}

		// --- Listings collection ---
		listings := core.NewBaseCollection("listings")
		listings.Fields.Add(
			&core.TextField{Name: "title", Required: true},
			&core.TextField{Name: "description", Required: false},
			&core.SelectField{
				Name:      "category",
				Required:  true,
				MaxSelect: 1,
				Values:    []string{"lend", "give", "request"},
			},
			&core.SelectField{
				Name:      "status",
				Required:  true,
				MaxSelect: 1,
				Values:    []string{"active", "claimed", "closed"},
			},
			&core.RelationField{
				Name:         "owner",
				Required:     true,
				CollectionId: users.Id,
				MaxSelect:    1,
			},
			&core.RelationField{
				Name:         "tower",
				Required:     true,
				CollectionId: towers.Id,
				MaxSelect:    1,
			},
		)
		listings.Indexes = []string{
			"CREATE INDEX idx_listings_tower_status ON listings (tower, status)",
			"CREATE INDEX idx_listings_owner ON listings (owner)",
		}
		if err := app.Save(listings); err != nil {
			return err
		}

		return nil
	}, func(app core.App) error {
		// Down migration: drop in reverse order.
		if listings, err := app.FindCollectionByNameOrId("listings"); err == nil {
			if err := app.Delete(listings); err != nil {
				return err
			}
		}

		if users, err := app.FindCollectionByNameOrId("users"); err == nil {
			if f := users.Fields.GetByName("apartment"); f != nil {
				users.Fields.RemoveById(f.GetId())
				if err := app.Save(users); err != nil {
					return err
				}
			}
		}

		if apartments, err := app.FindCollectionByNameOrId("apartments"); err == nil {
			if err := app.Delete(apartments); err != nil {
				return err
			}
		}

		if towers, err := app.FindCollectionByNameOrId("towers"); err == nil {
			if err := app.Delete(towers); err != nil {
				return err
			}
		}

		return nil
	})
}
