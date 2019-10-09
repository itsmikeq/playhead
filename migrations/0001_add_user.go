package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var addUserMigration_0001 = &Migration{
	Number: 1,
	Name:   "user_playheads",
	Forwards: func(db *gorm.DB) error {
		const addUserPlayheadSQL = `
			CREATE TABLE user_playheads (
 				id serial PRIMARY KEY,
 				user_uuid varchar (40) NOT NULL,
 				series_uuid varchar (40) NOT NULL,
 				episode_uuid varchar (40) NOT NULL,
 				created_at TIMESTAMP NOT NULL,
 				updated_at TIMESTAMP NOT NULL,
 				deleted_at TIMESTAMP);
		`

		err := db.Exec(addUserPlayheadSQL).Error
		return errors.Wrap(err, "unable to create user user_playheads table")
	},
}

func init() {
	Migrations = append(Migrations, addUserMigration_0001)
}
