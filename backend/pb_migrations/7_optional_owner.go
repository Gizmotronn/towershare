package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)
		col, err := dao.FindCollectionByNameOrId("listings")
		if err != nil {
			return err
		}
		f := col.Schema.GetFieldByName("owner")
		if f == nil {
			return nil
		}
		f.Required = false
		return dao.SaveCollection(col)
	}, func(db dbx.Builder) error {
		dao := daos.New(db)
		col, err := dao.FindCollectionByNameOrId("listings")
		if err != nil {
			return err
		}
		f := col.Schema.GetFieldByName("owner")
		if f == nil {
			return nil
		}
		f.Required = true
		return dao.SaveCollection(col)
	})
}
