package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var addIndexMigration_0004 = &Migration{
	Number: 4,
	Name:   "Add indexes",
	Forwards: func(db *gorm.DB) error {
		const addIndexes = `
			drop index user_uuid_series_idx;
			create unique index user_series_unq_idx on user_playheads (user_uuid, series_uuid);
		`

		err := db.Exec(addIndexes).Error
		return errors.Wrap(err, "unable to create series_info indexes")
	},
}

func init() {
	Migrations = append(Migrations, addIndexMigration_0004)
}
