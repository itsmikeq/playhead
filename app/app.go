package app

import (
	"github.com/sirupsen/logrus"
	"playhead/model"

	"playhead/db"
)

type App struct {
	Config   *Config
	Database *db.Database
}

func (a *App) NewContext() *Context {
	return &Context{
		Logger:   logrus.StandardLogger(),
		Database: a.Database,
	}
}

func New() (app *App, err error) {
	app = &App{}
	app.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := db.InitConfig()
	if err != nil {
		return nil, err
	}

	app.Database, err = db.New(dbConfig)
	if err != nil {
		return nil, err
	}

	return app, err
}

func (a *App) Close() error {
	return a.Database.Close()
}

type ValidationError struct {
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}

type UserPlayheadError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func (e *UserPlayheadError) Error() string {
	return e.Message
}

func (a *App) ValidateBearerToken(token string) (user *model.User, failureReason string, err error){
	return user, failureReason, err
}