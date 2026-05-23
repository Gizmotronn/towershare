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
		listings.Schema.AddField(&schema.SchemaField{
			Name: "images",
			Type: schema.FieldTypeFile,
			Options: &schema.FileOptions{
				MaxSelect: 5,
				MaxSize:   5242880, // 5 MB
				MimeTypes: []string{"image/jpeg", "image/png", "image/webp", "image/gif"},
				Thumbs:    []string{"100x100", "600x0"},
			},
		})
		return dao.SaveCollection(listings)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)
		listings, err := dao.FindCollectionByNameOrId("listings")
		if err != nil {
			return err
		}
		if f := listings.Schema.GetFieldByName("images"); f != nil {
			listings.Schema.RemoveField(f.Id)
			return dao.SaveCollection(listings)
		}
		return nil
	})
}
