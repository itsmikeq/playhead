package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"playhead/app"
	"playhead/model"
)

type UserInput struct {
	UserUUID    string `json:"user_uuid"`
	SessionID   string `json:"session_id"`
	SeriesUUID  string `json:"series_uuid"`
	EpisodeUUID string `json:"episode_uuid"`
}

type BigUserInput struct {
	UserUUID    string `json:"user_uuid"`
	SessionID   string `json:"session_id"`
	SeriesUUID  string `json:"series_uuid"`
	EpisodeUUID string `json:"episode_uuid"`
}

type UserResponse struct {
	Id uint `json:"id"`
}

//  Here we will snag the userID from the user's secret key
// and build the input from either the sessionID, or the user's uuid
// Once we get a session ID with a userID, we'll update all of the matching session IDs with the
// corresponding userID

func (a *API) CreateUserPlayhead(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input UserInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	playhead := &model.UserPlayhead{UserUUID: input.UserUUID, SeriesUUID: input.SeriesUUID, EpisodeUUID: input.EpisodeUUID}

	if err := ctx.CreateUserPlayhead(playhead); err != nil {
		return err
	}

	data, err := json.Marshal(&UserResponse{Id: playhead.ID})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) UpdateUserPlayhead(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input UserInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	playhead := &model.UserPlayhead{UserUUID: input.UserUUID, SeriesUUID: input.SeriesUUID, EpisodeUUID: input.EpisodeUUID}

	if err := ctx.UpdatePlayhead(playhead); err != nil {
		return err
	}

	data, err := json.Marshal(&UserResponse{Id: playhead.ID})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) DeleteUserPlayhead(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input UserInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	playhead := &model.UserPlayhead{UserUUID: input.UserUUID, SeriesUUID: input.SeriesUUID, EpisodeUUID: input.EpisodeUUID}

	if err := ctx.DeletePlayheadByUserIdAndSeriesId(playhead.UserUUID, playhead.SeriesUUID); err != nil {
		return err
	}

	data, err := json.Marshal(&UserResponse{Id: playhead.ID})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) GetUserPlayhead(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input UserInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	playhead := &model.UserPlayhead{UserUUID: input.UserUUID, SeriesUUID: input.SeriesUUID, EpisodeUUID: input.EpisodeUUID}

	if r, err := ctx.GetUserPlayheads(); err != nil {
		fmt.Println(r)
		return err
	}

	data, err := json.Marshal(&UserResponse{Id: playhead.ID})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

func (a *API) GetUserPlayheads(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input UserInput

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	// playheads := []model.UserPlayhead{}

	if r, err := ctx.GetUserPlayheads(); err != nil {
		fmt.Println(r)
		return err
	}

	data, err := json.Marshal(&UserResponse{Id: 1})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}
