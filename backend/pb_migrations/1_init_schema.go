package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/types"
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
			Options:  &schema.TextOptions{},
		})

		return dao.SaveCollection(users)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)

		users, err := dao.FindCollectionByNameOrId("users")
		if err != nil {
			return err
		}

		if f := users.Schema.GetFieldByName("display_name"); f != nil {
			users.Schema.RemoveField(f.Id)
			return dao.SaveCollection(users)
		}

		return nil
	})
}
