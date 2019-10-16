package queues

import (
	"github.com/sirupsen/logrus"
	"playhead/db"
)

type Queue struct {
	Config   *Config
	Database *db.Database
	Context  *Context
}

func (q *Queue) NewContext() *Context {
	return &Context{
		Logger:   logrus.StandardLogger(),
		Database: q.Database,
	}
}

func New() (q *Queue, err error) {
	q = &Queue{}
	q.Config, err = InitConfig()
	if err != nil {
		return nil, err
	}

	dbConfig, err := db.InitConfig()
	if err != nil {
		return nil, err
	}

	q.Database, err = db.New(dbConfig)
	if err != nil {
		return nil, err
	}

	return q, err
}

func (q *Queue) Close() error {
	return q.Database.Close()
}
