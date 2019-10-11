package api

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
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

type UserResponse struct {
	SeriesUUID  string `json:"series_uuid"`
	EpisodeUUID string `json:"episode_uuid"`
}

//  Here we will snag the userID from the user's secret key
// and build the input from either the sessionID, or the user's uuid
// Once we get a session ID with a userID, we'll update all of the matching session IDs with the
// corresponding userID

func (a *API) CreateUserPlayhead(ctx *app.Context, w http.ResponseWriter, r *http.Request) error {
	var input UserInput
	// fmt.Printf("Got user %+v\n", ctx.User)
	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &input); err != nil {
		return err
	}

	playhead := &model.UserPlayhead{UserUUID: ctx.User.UserID, SeriesUUID: input.SeriesUUID, EpisodeUUID: input.EpisodeUUID}

	if err := ctx.CreateUserPlayhead(playhead); err != nil {
		return err
	}

	data, err := json.Marshal(&UserResponse{SeriesUUID: playhead.SeriesUUID, EpisodeUUID: playhead.EpisodeUUID})
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

	playhead := &model.UserPlayhead{SeriesUUID: input.SeriesUUID, EpisodeUUID: input.EpisodeUUID}

	if err := ctx.UpdatePlayhead(playhead); err != nil {
		if data, errm := json.Marshal(&app.ValidationError{Message: fmt.Sprintf("%s", err)}); errm != nil {
			logrus.Error(errm)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(data)
		}
		return nil
	}

	data, err := json.Marshal(&UserResponse{SeriesUUID: playhead.SeriesUUID, EpisodeUUID: playhead.EpisodeUUID})
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

	playhead := &model.UserPlayhead{UserUUID: ctx.User.UserID, SeriesUUID: input.SeriesUUID, EpisodeUUID: input.EpisodeUUID}

	if err := ctx.DeletePlayheadBySeriesUUID(playhead.UserUUID, playhead.SeriesUUID); err != nil {
		if data, errm := json.Marshal(&app.ValidationError{Message: fmt.Sprintf("%s", err)}); errm != nil {
			logrus.Error(errm)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(data)
		}
		return nil
	}

	data, err := json.Marshal(&UserResponse{})
	if err != nil {
		return err
	}

	_, err = w.Write(data)
	return err
}

// Get a single playhead by user UUID and Series UUID
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

	if input.SeriesUUID == "" {
		// not found, carry on
		return nil
	}
	var playhead *model.UserPlayhead
	if playhead, err = ctx.GetPlayhead(input.SeriesUUID); err != nil {
		return err
	}
	data, err := json.Marshal(&UserResponse{EpisodeUUID: playhead.EpisodeUUID})
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

	if playheads, err := ctx.GetUserPlayheads(); err != nil {
		logrus.Error(err)
		return err
	} else if len(playheads) < 1 {
		data, err := json.Marshal(&UserResponse{})
		_, err = w.Write(data)
		return err
	} else {
		datas := make([]byte, 0)
		for playhead := range playheads {
			data, err := json.Marshal(&UserResponse{SeriesUUID: playheads[playhead].SeriesUUID, EpisodeUUID: playheads[playhead].EpisodeUUID})
			if err == nil {
				datas = append(datas[:], data...)
			}
		}
		_, err = w.Write(datas)
	}

	return err
}
