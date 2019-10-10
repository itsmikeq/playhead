package model

type UserPlayhead struct {
	Model
	UserUUID    string `json:"uuid" gorm:"type:varchar(40)"`
	SeriesUUID  string `json:"series_uuid" gorm:"type:varchar(40)"`
	EpisodeUUID string `json:"episode_uuid" gorm:"type:varchar(40)"`
}
