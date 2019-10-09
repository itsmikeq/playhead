package model

type SeriesInfo struct {
	Model

	SeriesUUID    string `json:"series_uuid"`
	EpisodeUUID   string `json:"episode_uuid"`
	EpisodeNumber int    `json:"episode_uuid"`
}
