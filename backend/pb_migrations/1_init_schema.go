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
		users, err := dao.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		users.Schema.AddField(&schema.SchemaField{
			Name:     "display_name",
			Type:     schema.FieldTypeText,
			Required: false,
		})
		return dao.SaveCollection(users)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)
		users, err := dao.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}
		field := users.Schema.GetFieldByName("display_name")
		if field != nil {
			users.Schema.RemoveField(field.Id)
			return dao.SaveCollection(users)
		}
		return nil
	})
}
