package model

type UserPlayhead struct {
	Model
	UserUUID    string `json:"uuid"`
	SeriesUUID  string `json:"series_uuid"`
	EpisodeUUID string `json:"episode_uuid"`
}
