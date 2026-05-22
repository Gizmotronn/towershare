package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)
		listings, err := dao.FindCollectionByNameOrId("listings")
		if err != nil {
			return err
		}

		// Extend category to include hard-waste types
		if f := listings.Schema.GetFieldByName("category"); f != nil {
			opts, ok := f.Options.(*schema.SelectOptions)
			if ok {
				opts.Values = []string{
					"lend", "give", "request",
					"furniture", "whitegoods", "ewaste", "mattress",
				}
			}
		}

		// Extend status to include hard-waste values
		if f := listings.Schema.GetFieldByName("status"); f != nil {
			opts, ok := f.Options.(*schema.SelectOptions)
			if ok {
				opts.Values = []string{"active", "claimed", "closed", "available", "collected"}
			}
		}

		// Make tower optional
		if f := listings.Schema.GetFieldByName("tower"); f != nil {
			f.Required = false
		}

		// Add estimatedM3
		listings.Schema.AddField(&schema.SchemaField{
			Name:     "estimated_m3",
			Type:     schema.FieldTypeNumber,
			Required: false,
		})

		// Add pickupBy
		listings.Schema.AddField(&schema.SchemaField{
			Name:     "pickup_by",
			Type:     schema.FieldTypeText,
			Required: false,
		})

		return dao.SaveCollection(listings)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)
		listings, err := dao.FindCollectionByNameOrId("listings")
		if err != nil {
			return err
		}

		// Restore category
		if f := listings.Schema.GetFieldByName("category"); f != nil {
			if opts, ok := f.Options.(*schema.SelectOptions); ok {
				opts.Values = []string{"lend", "give", "request"}
			}
		}

		// Restore status
		if f := listings.Schema.GetFieldByName("status"); f != nil {
			if opts, ok := f.Options.(*schema.SelectOptions); ok {
				opts.Values = []string{"active", "claimed", "closed"}
			}
		}

		// Make tower required again
		if f := listings.Schema.GetFieldByName("tower"); f != nil {
			f.Required = true
		}

		// Remove added fields
		for _, name := range []string{"estimated_m3", "pickup_by"} {
			if f := listings.Schema.GetFieldByName(name); f != nil {
				listings.Schema.RemoveField(f.Id)
			}
		}

		return dao.SaveCollection(listings)
	})
}
