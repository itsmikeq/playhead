package queues

import (
	"errors"
	"github.com/sirupsen/logrus"
	"playhead/db"
	"playhead/model"
)

type QMessageBody struct {
	UserUUID    string `json:"user_uuid" binding:"required"`
	RequestID   string `json:"request_id" binding:"required"`
	RequestType string `json:"request_type" binding:"required"`
	ServiceName string `json:"service_name" binding:"required"`
}

type QMessage struct {
	Subject     string `json:"Subject"`
	Message     string `json:"Message"`
	MessageBody QMessageBody
}

type Context struct {
	Logger       logrus.FieldLogger
	Database     *db.Database
	UserPlayhead *model.UserPlayhead
	User         *model.User
}

// Messages:
// UserDataDownloadRequest | UserDataDeleteRequest
// {"request_id":"uuid1","request_type":"UserDataDownloadRequest","user_uuid":"bb70da7e-a5c1-455e-9f3f-74208fdee1f5","created_at":"2019-04-23T17:54:36.000Z"}
// {"request_id":"uuid2","request_type":"UserDataDownloadRequest","user_uuid":"0e16e2bb-ac83-4cd6-b320-77abcbbc820e","created_at":"2019-04-23T17:54:36.000Z"}
// {"request_id":"uuid3","request_type":"UserDataDeleteRequest","user_uuid":"0e16e2bb-ac83-4cd6-b320-77abcbbc820e","created_at":"2019-04-23T17:54:36.000Z"}

func CheckForEmpty(qMessage QMessage) error {
	if len(qMessage.MessageBody.RequestID) < 1 {
		return errors.New("missing RequestID")
	}
	return nil
}
