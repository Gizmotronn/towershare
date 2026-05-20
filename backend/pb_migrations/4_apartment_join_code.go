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
		apartments, err := dao.FindCollectionByNameOrId("apartments")
		if err != nil {
			return err
		}
		apartments.Schema.AddField(&schema.SchemaField{
			Name:     "join_code",
			Type:     schema.FieldTypeText,
			Required: false,
		})
		return dao.SaveCollection(apartments)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)
		apartments, err := dao.FindCollectionByNameOrId("apartments")
		if err != nil {
			return err
		}
		if f := apartments.Schema.GetFieldByName("join_code"); f != nil {
			apartments.Schema.RemoveField(f.Id)
			return dao.SaveCollection(apartments)
		}
		return nil
	})
}
