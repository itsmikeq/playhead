package queues

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func (q *Queue) StartGdprListener(ctx *Context) {
	chnMessages := make(chan *sqs.Message, int64(q.Config.SqsMaxMessages))
	go q.pollGdprSqs(chnMessages)

	fmt.Printf("Listening on stack queue: %s\n", string(q.Config.GdprQueueUrl))

	go func() {
		for message := range chnMessages {
			if err := q.handleGdprMessage(message); ErrorHandler(err) {
				fmt.Printf("Error with handling message %v\n", err)
			} else {
				// error handled in subroutine
				q.deleteQMessage(message, string(q.Config.GdprQueueUrl))
			}
		}
	}()
}

func (q *Queue) deleteGdprQMessage(message *sqs.Message) {
	if _, err := q.getSQSSession().DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(string(q.Config.GdprQueueUrl)),
		ReceiptHandle: message.ReceiptHandle,
	}); err != nil {
		ErrorHandler(err)
		fmt.Printf("Error removing message from queue: %v\n", err)
	}
}

func (q *Queue) pollGdprSqs(chn chan<- *sqs.Message) {
	for {
		output, err := q.getSQSSession().ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(string(q.Config.GdprQueueUrl)),
			MaxNumberOfMessages: aws.Int64(int64(q.Config.SqsMaxMessages)),
			WaitTimeSeconds:     aws.Int64(int64(q.Config.TimeWaitSeconds)),
		})

		if err != nil {
			logrus.Errorf("failed to fetch sqs message %v", err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}

}

func (q *Queue) handleGdprMessage(message *sqs.Message) error {
	var qMessage QMessage
	// fmt.Printf("Got Message:\n%v\n", message)
	qMessage.MessageBody.ServiceName = "playhead"
	if err := json.Unmarshal([]byte(aws.StringValue(message.Body)), &qMessage); ErrorHandler(err) {
		ErrorLogger(fmt.Errorf("Error decoding json: '%v'\n", err))
	}
	// fmt.Printf("Subject as %v\n", qMessage.Subject)
	msg := qMessage.Message
	unmerr := json.Unmarshal([]byte(msg), &qMessage.MessageBody)
	if unmerr != nil {
		ErrorLogger(unmerr)
		return unmerr
	}
	if qMessage.Subject == "UserDataDownloadRequest" || qMessage.MessageBody.RequestType == "UserDataDownloadRequest" {
		downErr := q.handleGdprDownload(qMessage)
		if downErr != nil {
			// Error is handled and logged
			ErrorLogger(downErr)
			return downErr
		}
	} else if qMessage.Subject == "UserDataDeleteRequest" || qMessage.MessageBody.RequestType == "UserDataDeleteRequest" {
		delerr := q.handleGdprDeleteReq(qMessage)
		if delerr != nil {
			ErrorLogger(delerr)
		}
	}
	return nil
}

func (q *Queue) handleGdprDeleteReq(qMessage QMessage) error {
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
	if err := CheckForEmpty(qMessage); err != nil {
		pm.ErrorMessage = err.Error()
		pm.Success = false
		if perr := q.PublishGdpr(pm); ErrorHandler(perr) {
			return perr
		}
		return err
	}

	if playheads, err := q.Database.GetUserPlayheads(pm.UserUUID); err != nil {
		logrus.Error(err.Error())
		pm.ErrorMessage = fmt.Sprintf("error deleting last played %v\n", err.Error())
		pm.Success = false
	} else {
		del := q.Database.Delete(&playheads)
		if del.Error != nil {
			logrus.Error(del.Error)
			return del.Error
		}
		if err := q.PublishGdpr(pm); err != nil {
			logrus.Error(err)
		}
		return nil
	}
	return nil
}

func (q *Queue) handleGdprDownload(qMessage QMessage) error {
	pm := PublishMessage{
		RequestID:    qMessage.MessageBody.RequestID,
		RequestType:  qMessage.MessageBody.RequestType,
		UserUUID:     qMessage.MessageBody.UserUUID,
		ServiceName:  "playhead",
		ErrorMessage: "",
		Success:      true,
	}
	if err := CheckForEmpty(qMessage); err != nil {
		pm.ErrorMessage = err.Error()
		pm.Success = false
		if perr := q.PublishGdpr(pm); ErrorHandler(perr) {
			return perr
		}
		return err
	}

	if playheads, err := q.Database.GetUserPlayheads(pm.UserUUID); err != nil {
		return err
	} else {
		filePath := filepath.Join(string(q.Config.GdprBasePath), qMessage.MessageBody.RequestID, "playheads.json")
		if data, err := json.Marshal(&playheads); err == nil {
			if err := q.AddFileToS3(filePath, string(data)); ErrorHandler(err) {
				pm.Success = false
				pm.ErrorMessage = err.Error()
				if perr := q.PublishGdpr(pm); ErrorHandler(perr) {
					return perr
				}
				return err
			} else {
				pm.S3Location = filePath
				if err := q.PublishGdpr(pm); ErrorHandler(err) {
					return err
				}
			}

		}
	}
	return nil
}

// Private
