package queues

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

func (q *Queue) StartListener(ctx *Context) {
	chnMessages := make(chan *sqs.Message, sqsMaxMessages)
	go q.pollSqs(chnMessages)

	fmt.Printf("Listening on stack queue: %s\n", string(q.Config.GdprQueueUrl))

	go func() {
		for message := range chnMessages {
			if err := q.handleMessage(message); ErrorHandler(err) {
				fmt.Printf("Error with handling message %v\n", err)
			} else {
				q.deleteMessage(message)
			}
		}
	}()
}


func (q *Queue) deleteMessage(message *sqs.Message) {
	if _, err := q.getSQSSession().DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(string(q.Config.GdprQueueUrl)),
		ReceiptHandle: message.ReceiptHandle,
	}); err != nil {
		ErrorHandler(err)
		fmt.Printf("Error removing message from queue: %v\n", err)
	}
}

func (q *Queue) pollSqs(chn chan<- *sqs.Message) {
	for {
		output, err := q.getSQSSession().ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(string(q.Config.GdprQueueUrl)),
			MaxNumberOfMessages: aws.Int64(sqsMaxMessages),
			WaitTimeSeconds:     aws.Int64(timeWaitSeconds),
		})

		if err != nil {
			logrus.Errorf("failed to fetch sqs message %v", err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}

}

func (q *Queue) handleMessage(message *sqs.Message) error {
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
		downErr := q.handleDownload(qMessage)
		if downErr != nil {
			// Error is handled and logged
			ErrorLogger(downErr)
			return downErr
		}
	} else if qMessage.Subject == "UserDataDeleteRequest" || qMessage.MessageBody.RequestType == "UserDataDeleteRequest" {
		delerr := q.handleDelete(qMessage)
		if delerr != nil {
			ErrorLogger(delerr)
		}
	}
	return nil
}

func (q *Queue) handleDelete(qMessage QMessage) error {
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
		if perr := q.Publish(pm); ErrorHandler(perr) {
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
		if err := q.Publish(pm); err != nil {
			logrus.Error(err)
		}
		return nil
	}
	return nil
}

func (q *Queue) handleDownload(qMessage QMessage) error {
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
		if perr := q.Publish(pm); ErrorHandler(perr) {
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
				if perr := q.Publish(pm); ErrorHandler(perr) {
					return perr
				}
				return err
			} else {
				pm.S3Location = filePath
				if err := q.Publish(pm); ErrorHandler(err) {
					return err
				}
			}

		}
	}
	return nil
}
