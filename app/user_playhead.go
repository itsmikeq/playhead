package app

import (
	"playhead/model"
)

func (a *App) GetUserPlayhead(userUUID string, seriesUUID string) (*model.UserPlayhead, error) {
	return a.Database.GetPlayheadByUserUUIDAndSeriesUUID(userUUID, seriesUUID)
}

func (a *App) UpdatePlayhead(userUUID string, seriesUUID string, episodeUUID string) (*model.UserPlayhead, error) {
	ph, err := a.Database.GetPlayheadByUserUUIDAndSeriesUUID(userUUID, seriesUUID)
	if err != nil {
		return nil, err
	}
	ph.EpisodeUUID = episodeUUID
	if err := a.Database.UpdateUserPlayhead(ph); err != nil {
		return nil, err
	}
	return ph, nil
}

func (ctx *Context) CreateUserPlayhead(playhead *model.UserPlayhead) error {
	if err := ctx.validatePlayhead(playhead); err != nil {
		return err
	}

	return ctx.Database.CreateUserPlayhead(playhead)
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

func (ctx *Context) GetUserPlayheads() ([]*model.UserPlayhead, error) {
	if ctx.User == nil {
		return nil, ctx.AuthorizationError()
	}

	return ctx.getPlayheadsByUserId(ctx.User.UserID)
}

func (ctx *Context) getPlayheadsByUserId(userUUID string) ([]*model.UserPlayhead, error) {
	return ctx.Database.GetUserPlayheads(userUUID)
}

func (ctx *Context) getPlayheadsByUserIdAndSeriesId(userUUID string, seriesUUID string) (*model.UserPlayhead, error) {
	var userPlayhead model.UserPlayhead
	r := ctx.Database.First(&userPlayhead, &model.UserPlayhead{UserUUID: userUUID, SeriesUUID: seriesUUID})
	return &userPlayhead, r.Error
}

func (ctx *Context) UpdatePlayhead(playhead *model.UserPlayhead) error {
	if ctx.User == nil {
		return ctx.AuthorizationError()
	}

	if playhead.UserUUID != ctx.User.UserID {
		return ctx.AuthorizationError()
	}

	if err := ctx.validatePlayhead(playhead); err != nil {
		return nil
	}

	return ctx.Database.Update(playhead).Error
}

func (ctx *Context) DeletePlayheadByUserIdAndSeriesId(userUUID string, seriesUUID string) error {
	if ctx.User == nil {
		return ctx.AuthorizationError()
	}

	playhead, err := ctx.getPlayheadsByUserIdAndSeriesId(userUUID, seriesUUID)
	if err != nil {
		return err
	}

	if playhead.UserUUID != ctx.User.UserID {
		return ctx.AuthorizationError()
	}

	return err
}
