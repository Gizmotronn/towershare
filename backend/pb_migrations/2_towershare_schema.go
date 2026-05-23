package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
)

func intPtr(v int) *int { return &v }

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		// --- Towers ---
		towers := &models.Collection{
			Name:   "towers",
			Type:   models.CollectionTypeBase,
			Schema: schema.NewSchema(),
			Indexes: types.JsonArray[string]{
				"CREATE UNIQUE INDEX idx_towers_name ON towers (name)",
			},
		}
		towers.Schema.AddField(&schema.SchemaField{Name: "name", Type: schema.FieldTypeText, Required: true})
		towers.Schema.AddField(&schema.SchemaField{Name: "address", Type: schema.FieldTypeText, Required: true})
		if err := dao.SaveCollection(towers); err != nil {
			return err
		}

		// --- Apartments ---
		apartments := &models.Collection{
			Name:   "apartments",
			Type:   models.CollectionTypeBase,
			Schema: schema.NewSchema(),
			Indexes: types.JsonArray[string]{
				"CREATE UNIQUE INDEX idx_apartments_unit_tower ON apartments (unit_number, tower)",
			},
		}
		apartments.Schema.AddField(&schema.SchemaField{Name: "unit_number", Type: schema.FieldTypeText, Required: true})
		apartments.Schema.AddField(&schema.SchemaField{Name: "floor", Type: schema.FieldTypeNumber, Required: false})
		apartments.Schema.AddField(&schema.SchemaField{
			Name:     "tower",
			Type:     schema.FieldTypeRelation,
			Required: true,
			Options: &schema.RelationOptions{
				CollectionId:  towers.Id,
				MaxSelect:     intPtr(1),
				CascadeDelete: true,
			},
		})
		if err := dao.SaveCollection(apartments); err != nil {
			return err
		}

		// --- Extend users: apartment relation ---
		users, err := dao.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		users.Schema.AddField(&schema.SchemaField{
			Name:     "apartment",
			Type:     schema.FieldTypeRelation,
			Required: false,
			Options: &schema.RelationOptions{
				CollectionId: apartments.Id,
				MaxSelect:    intPtr(1),
			},
		})
		if err := dao.SaveCollection(users); err != nil {
			return err
		}

		// --- Add managers relation to towers ---
		towersRefresh, err := dao.FindCollectionByNameOrId("towers")
		if err != nil {
			return err
		}
		towersRefresh.Schema.AddField(&schema.SchemaField{
			Name:     "managers",
			Type:     schema.FieldTypeRelation,
			Required: false,
			Options: &schema.RelationOptions{
				CollectionId: users.Id,
				MaxSelect:    intPtr(100),
			},
		})
		if err := dao.SaveCollection(towersRefresh); err != nil {
			return err
		}

		// --- Listings ---
		listings := &models.Collection{
			Name:   "listings",
			Type:   models.CollectionTypeBase,
			Schema: schema.NewSchema(),
			Indexes: types.JsonArray[string]{
				"CREATE INDEX idx_listings_tower_status ON listings (tower, status)",
				"CREATE INDEX idx_listings_owner ON listings (owner)",
			},
		}
		listings.Schema.AddField(&schema.SchemaField{Name: "title", Type: schema.FieldTypeText, Required: true})
		listings.Schema.AddField(&schema.SchemaField{Name: "description", Type: schema.FieldTypeText, Required: false})
		listings.Schema.AddField(&schema.SchemaField{
			Name:     "category",
			Type:     schema.FieldTypeSelect,
			Required: true,
			Options:  &schema.SelectOptions{MaxSelect: 1, Values: []string{"lend", "give", "request"}},
		})
		listings.Schema.AddField(&schema.SchemaField{
			Name:     "status",
			Type:     schema.FieldTypeSelect,
			Required: true,
			Options:  &schema.SelectOptions{MaxSelect: 1, Values: []string{"active", "claimed", "closed"}},
		})
		listings.Schema.AddField(&schema.SchemaField{
			Name:     "owner",
			Type:     schema.FieldTypeRelation,
			Required: true,
			Options:  &schema.RelationOptions{CollectionId: users.Id, MaxSelect: intPtr(1)},
		})
		listings.Schema.AddField(&schema.SchemaField{
			Name:     "tower",
			Type:     schema.FieldTypeRelation,
			Required: true,
			Options:  &schema.RelationOptions{CollectionId: towers.Id, MaxSelect: intPtr(1)},
		})
		return dao.SaveCollection(listings)

	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		// Drop in reverse dependency order.
		for _, name := range []string{"listings", "apartments", "towers"} {
			col, err := dao.FindCollectionByNameOrId(name)
			if err != nil {
				continue // already gone
			}
			if err := dao.DeleteCollection(col); err != nil {
				return err
			}
		}

		// Remove apartment field from users.
		users, err := dao.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		if f := users.Schema.GetFieldByName("apartment"); f != nil {
			users.Schema.RemoveField(f.Id)
			return dao.SaveCollection(users)
		}
		return nil
	})
}
