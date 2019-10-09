package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var addSeriesInfoMigration_0002 = &Migration{
	Number: 2,
	Name:   "Add series info",
	Forwards: func(db *gorm.DB) error {
		const addUserSQL = `
			CREATE TABLE series_info(
 				id serial PRIMARY KEY,
 				series_uuid varchar (40) NOT NULL,
 				episode_uuid varchar (40) NOT NULL,
 				episode_number int not null,
 				created_at TIMESTAMP NOT NULL,
 				updated_at TIMESTAMP NOT NULL,
 				deleted_at TIMESTAMP);
		`

		err := db.Exec(addUserSQL).Error
		return errors.Wrap(err, "unable to create series_info table")
	},
}

func init() {
	Migrations = append(Migrations, addSeriesInfoMigration_0002)
}
