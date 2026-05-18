package migrations

import (
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db)

		open := ""
		for _, name := range []string{"listings", "towers", "apartments", "users"} {
			col, err := dao.FindCollectionByNameOrId(name)
			if err != nil {
				continue
			}
			col.ListRule = &open
			col.ViewRule = &open
			col.CreateRule = &open
			col.UpdateRule = &open
			col.DeleteRule = &open
			if err := dao.SaveCollection(col); err != nil {
				return err
			}
		}
		return nil
	}, func(db dbx.Builder) error {
		// Rules were not set before — nothing to restore.
		return nil
	})
}
