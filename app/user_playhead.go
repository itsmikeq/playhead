package app

import (
	"github.com/sirupsen/logrus"
	"playhead/model"
)

// Get a single playhead
func (ctx *Context) GetPlayhead(seriesUUID string) (*model.UserPlayhead, error) {
	if ctx.User == nil {
		return nil, ctx.AuthorizationError()
	}
	return ctx.Database.GetPlayheadByUserUUIDAndSeriesUUID(ctx.User.UserID, seriesUUID)
}

func (ctx *Context) GetUserPlayheads() ([]*model.UserPlayhead, error) {
	if ctx.User == nil {
		return nil, ctx.AuthorizationError()
	}
	return ctx.Database.GetUserPlayheads(ctx.User.UserID)
}

func (ctx *Context) GetPlayheadBySeriesId(seriesUUID string) (*model.UserPlayhead, error) {
	if ctx.User == nil {
		return nil, ctx.AuthorizationError()
	}
	return ctx.Database.GetPlayheadByUserUUIDAndSeriesUUID(ctx.User.UserID, seriesUUID)
}

func (ctx *Context) CreateUserPlayhead(playhead *model.UserPlayhead) error {

	if ctx.User == nil {
		return ctx.AuthorizationError()
	}

	if err := ctx.validatePlayhead(playhead); err != nil {
		return err
	}

	return ctx.Database.CreateUserPlayhead(playhead)
}

func (ctx *Context) UpdatePlayhead(playhead *model.UserPlayhead) error {
	if ctx.User == nil {
		return ctx.AuthorizationError()
	}

	if err := ctx.validatePlayhead(playhead); err != nil {
		return err
	}

	newEpisode := playhead.EpisodeUUID
	if playhead, err := ctx.Database.GetPlayheadByUserUUIDAndSeriesUUID(playhead.UserUUID, playhead.SeriesUUID); err == nil {
		playhead.EpisodeUUID = newEpisode
		if playhead.UserUUID != ctx.User.UserID {
			return ctx.AuthorizationError()
		}
		return ctx.Database.Save(playhead).Error
	} else {
		logrus.Error(err)
		return err
	}
}

func (ctx *Context) DeletePlayheadBySeriesUUID(userUUID string, seriesUUID string) error {
	if ctx.User == nil {
		return ctx.AuthorizationError()
	}

	playhead, err := ctx.GetPlayheadBySeriesId(seriesUUID)
	if err != nil {
		return err
	}

	if playhead.UserUUID != ctx.User.UserID {
		return ctx.AuthorizationError()
	}
	ctx.Database.Delete(&playhead)

	return err
}

func (ctx *Context) validatePlayhead(user *model.UserPlayhead) *ValidationError {
	// naive email validation
	if (len(user.EpisodeUUID) < 1) {
		return &ValidationError{"Missing episode UUID"}
	}
	if (len(user.SeriesUUID) < 1) {
		return &ValidationError{"Missing series UUID"}
	}
	return nil
}
