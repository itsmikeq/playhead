package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	"playhead/model"
)

func (db *Database) GetPlayheadByUserUUIDAndSeriesUUID(userUUID string, seriesUUID string) (*model.UserPlayhead, error) {
	var userPlayhead model.UserPlayhead

	if err := db.First(&userPlayhead, model.UserPlayhead{SeriesUUID: seriesUUID, UserUUID: userUUID}).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "unable to get userPlayhead")
	}

	return &userPlayhead, nil
}

func (db *Database) CreateUserPlayhead(userPlayhead *model.UserPlayhead) error {
	return db.Create(userPlayhead).Error
}

func (db *Database) UpdateUserPlayhead(userPlayhead *model.UserPlayhead) error {
	return db.Update(userPlayhead).Error
}

func (db *Database) DeleteUserPlayhead(userPlayhead *model.UserPlayhead) error {
	return db.Delete(userPlayhead).Error
}

func (db *Database) GetUserPlayheads(userUUID string) ([]*model.UserPlayhead, error) {
	var userPlayheads []*model.UserPlayhead
	r := db.Find(userPlayheads, &model.UserPlayhead{UserUUID: userUUID});
	if r.Error != nil {
		if gorm.IsRecordNotFoundError(r.Error) {
			return nil, nil
		}
		return nil, errors.Wrap(r.Error, "unable to get userPlayhead")
	}
	return userPlayheads, r.Error
}
