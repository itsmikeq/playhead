package queues

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"playhead/db"
	"playhead/model"
)

var sqsMaxMessages = int64(1)
var timeWaitSeconds = int64(1)

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

func (q *Queue) StartListener(ctx *Context) {
	chnMessages := make(chan *sqs.Message, sqsMaxMessages)
	go q.pollSqs(chnMessages)

	fmt.Printf("Listening on stack queue: %s\n", getListenQueueUrl())

	go func() {
		for message := range chnMessages {
			if err := ctx.handleMessage(message); ErrorHandler(err) {
				fmt.Printf("Error with handling message %v\n", err)
			} else {
				q.deleteMessage(message)
			}
		}
	}()
}

// Messages:
// UserDataDownloadRequest | UserDataDeleteRequest
// {"request_id":"uuid1","request_type":"UserDataDownloadRequest","user_uuid":"bb70da7e-a5c1-455e-9f3f-74208fdee1f5","created_at":"2019-04-23T17:54:36.000Z"}
// {"request_id":"uuid2","request_type":"UserDataDownloadRequest","user_uuid":"0e16e2bb-ac83-4cd6-b320-77abcbbc820e","created_at":"2019-04-23T17:54:36.000Z"}
// {"request_id":"uuid3","request_type":"UserDataDeleteRequest","user_uuid":"0e16e2bb-ac83-4cd6-b320-77abcbbc820e","created_at":"2019-04-23T17:54:36.000Z"}

func (ctx *Context) handleMessage(message *sqs.Message) error {
	var qMessage QMessage
	// fmt.Printf("Got Message:\n%v\n", message)
	qMessage.MessageBody.ServiceName = "playhead"
	if err := json.Unmarshal([]byte(aws.StringValue(message.Body)), &qMessage); ErrorHandler(err) {
		fmt.Printf("Error decoding json: '%v'\n", err)
	}
	// fmt.Printf("Subject as %v\n", qMessage.Subject)
	msg := qMessage.Message
	if err := json.Unmarshal([]byte(msg), &qMessage.MessageBody); ErrorHandler(err) {
		fmt.Println(err)
	}
	if qMessage.Subject == "UserDataDownloadRequest" || qMessage.MessageBody.RequestType == "UserDataDownloadRequest" {
		if downErr := ctx.handleDownload(qMessage); ErrorHandler(downErr) {
			// Error is handled and logged
			return nil
		}
	} else if qMessage.Subject == "UserDataDeleteRequest" || qMessage.MessageBody.RequestType == "UserDataDeleteRequest" {
		if delerr := ctx.handleDelete(qMessage); ErrorHandler(delerr) {
			// Error is handled and logged
			return nil
		}
	}
	return nil
}

func checkForEmpty(qMessage QMessage) error {
	if len(qMessage.MessageBody.RequestID) < 1 {
		return errors.New("missing RequestID")
	}
	return nil
}

func (ctx *Context) handleDelete(qMessage QMessage) error {
	pm := PublishMessage{
		RequestID:    qMessage.MessageBody.RequestID,
		RequestType:  qMessage.MessageBody.RequestType,
		UserUUID:     qMessage.MessageBody.UserUUID,
		ServiceName:  "playhead",
		S3Location:   "",
		ErrorMessage: "",
		Success:      true,
	}
	// fmt.Printf("Doing job %v\n", pm)
	if err := checkForEmpty(qMessage); err != nil {
		pm.ErrorMessage = err.Error()
		pm.Success = false
		if perr := Publish(pm); ErrorHandler(perr) {
			return perr
		}
		return err
	}

	if playheads, err := ctx.Database.GetUserPlayheads(pm.UserUUID); err != nil {
		logrus.Error(err.Error())
		pm.ErrorMessage = fmt.Sprintf("error deleting last played %v\n", err.Error())
		pm.Success = false
	} else {
		del := ctx.Database.Delete(&playheads)
		if del.Error != nil {
			logrus.Error(del.Error)
			return del.Error
		}
		if err := Publish(pm); err != nil {
			logrus.Error(err)
		}
		return nil
	}
	return nil
}

func (ctx *Context) handleDownload(qMessage QMessage) error {
	pm := PublishMessage{
		RequestID:    qMessage.MessageBody.RequestID,
		RequestType:  qMessage.MessageBody.RequestType,
		UserUUID:     qMessage.MessageBody.UserUUID,
		ServiceName:  "playhead",
		ErrorMessage: "",
		Success:      true,
	}
	if err := checkForEmpty(qMessage); err != nil {
		pm.ErrorMessage = err.Error()
		pm.Success = false
		if perr := Publish(pm); ErrorHandler(perr) {
			return perr
		}
		return err
	}

	if playheads, err := ctx.Database.GetUserPlayheads(pm.UserUUID); err != nil {
		return err
	} else {
		filePath := filepath.Join(getGdprBasePath(), qMessage.MessageBody.RequestID, "playheads.json")
		if data, err := json.Marshal(&playheads); err == nil {
			if err := AddFileToS3(filePath, string(data)); ErrorHandler(err) {
				pm.Success = false
				pm.ErrorMessage = err.Error()
				if perr := Publish(pm); ErrorHandler(perr) {
					return perr
				}
				return err
			} else {
				pm.S3Location = filePath
				if err := Publish(pm); ErrorHandler(err) {
					return err
				}
			}

		}
	}
	return nil
}

func (q *Queue) deleteMessage(message *sqs.Message) {
	if _, err := getSQSSession().DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(string(q.Config.GdprQueueUrl)),
		ReceiptHandle: message.ReceiptHandle,
	}); err != nil {
		ErrorHandler(err)
		fmt.Printf("Error removing message from queue: %v\n", err)
	}
}

func (q *Queue) pollSqs(chn chan<- *sqs.Message) {
	for {
		output, err := getSQSSession().ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(string(q.Config.GdprQueueUrl)),
			MaxNumberOfMessages: aws.Int64(sqsMaxMessages),
			WaitTimeSeconds:     aws.Int64(timeWaitSeconds),
		})

		if err != nil {
			exitErrorf("failed to fetch sqs message %v", err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}

}
