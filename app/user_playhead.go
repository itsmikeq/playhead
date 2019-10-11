package app

import (
	"fmt"
	"playhead/model"
)

func (a *App) GetUserPlayhead(userUUID string, seriesUUID string) (*model.UserPlayhead, error) {
	return a.Database.GetPlayheadByUserUUIDAndSeriesUUID(userUUID, seriesUUID)
}

func (a *App) UpdatePlayhead(userUUID string, seriesUUID string, episodeUUID string) (*model.UserPlayhead, error) {
	ph, err := a.Database.GetPlayheadByUserUUIDAndSeriesUUID(userUUID, seriesUUID)
	fmt.Printf("PH: %+v\n", ph)
	if err != nil {
		return nil, err
	}
	ph.EpisodeUUID = episodeUUID
	fmt.Printf("PH: %+v\n", ph)
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
	return ctx.Database.GetUserPlayheads(ctx.User.UserID)
}

func (ctx *Context) GetPlayheadBySeriesId(seriesUUID string) (*model.UserPlayhead, error) {
	if ctx.User == nil {
		return nil, ctx.AuthorizationError()
	}
	return ctx.Database.GetPlayheadByUserUUIDAndSeriesUUID(ctx.User.UserID, seriesUUID)
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
	newEpisode := playhead.EpisodeUUID
	if playhead, err := ctx.Database.GetPlayheadByUserUUIDAndSeriesUUID(playhead.UserUUID, playhead.SeriesUUID); err == nil {
		playhead.EpisodeUUID = newEpisode
		return ctx.Database.Save(playhead).Error
	}
	return nil
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
