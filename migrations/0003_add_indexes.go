package migrations

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

var addIndexMigration_0003 = &Migration{
	Number: 3,
	Name:   "Add indexes",
	Forwards: func(db *gorm.DB) error {
		const addIndexes = `
			CREATE INDEX series_info_series_idx
			ON series_info (series_uuid);
			CREATE INDEX series_info_episodes_idx
			ON series_info (episode_uuid);
			create index user_uuid_idx on user_playheads (user_uuid);
			create index user_uuid_series_idx on user_playheads (user_uuid, series_uuid);
			create index user_uuid_episode_idx on user_playheads (user_uuid, episode_uuid);
		`

		err := db.Exec(addIndexes).Error
		return errors.Wrap(err, "unable to create series_info indexes")
	},
}

func init() {
	Migrations = append(Migrations, addIndexMigration_0003)
}
