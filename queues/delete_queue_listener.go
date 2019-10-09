package queues

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"path/filepath"
	"wtc_go/wtc"
)

var sqsMaxMessages = int64(1)
var timeWaitSeconds = int64(1)

type QMessageBody struct {
	UserUUID    string      `json:"user_uuid" binding:"required"`
	RequestID   string 		`json:"request_id" binding:"required"`
	RequestType string      `json:"request_type" binding:"required"`
	ServiceName string      `json:"service_name" binding:"required"`
}

type QMessage struct {
	Subject     string `json:"Subject"`
	Message     string `json:"Message"`
	MessageBody QMessageBody
}

func StartListener() {
	chnMessages := make(chan *sqs.Message, sqsMaxMessages)
	go pollSqs(chnMessages)

	fmt.Printf("Listening on stack queue: %s\n", getListenQueueUrl())

	go func() {
		for message := range chnMessages {
			if err := handleMessage(message); ErrorHandler(err) {
				fmt.Printf("Error with handling message %v\n", err)
			} else {
				deleteMessage(message)
			}
		}
	}()
}

// Messages:
// UserDataDownloadRequest | UserDataDeleteRequest
// {"request_id":"uuid1","request_type":"UserDataDownloadRequest","user_uuid":"bb70da7e-a5c1-455e-9f3f-74208fdee1f5","created_at":"2019-04-23T17:54:36.000Z"}
// {"request_id":"uuid2","request_type":"UserDataDownloadRequest","user_uuid":"0e16e2bb-ac83-4cd6-b320-77abcbbc820e","created_at":"2019-04-23T17:54:36.000Z"}
// {"request_id":"uuid3","request_type":"UserDataDeleteRequest","user_uuid":"0e16e2bb-ac83-4cd6-b320-77abcbbc820e","created_at":"2019-04-23T17:54:36.000Z"}

func handleMessage(message *sqs.Message) error {
	var qMessage QMessage
	// fmt.Printf("Got Message:\n%v\n", message)
	qMessage.MessageBody.ServiceName = "wtc"
	if err := json.Unmarshal([]byte(aws.StringValue(message.Body)), &qMessage); ErrorHandler(err) {
		fmt.Printf("Error decoding json: '%v'\n", err)
	}
	// fmt.Printf("Subject as %v\n", qMessage.Subject)
	msg := qMessage.Message
	if err := json.Unmarshal([]byte(msg), &qMessage.MessageBody); ErrorHandler(err) {
		fmt.Println(err)
	}
	if qMessage.Subject == "UserDataDownloadRequest" || qMessage.MessageBody.RequestType == "UserDataDownloadRequest" {
		if downErr := handleDownload(qMessage); ErrorHandler(downErr) {
			// Error is handled and logged
			return nil
		}
	} else if qMessage.Subject == "UserDataDeleteRequest" || qMessage.MessageBody.RequestType == "UserDataDeleteRequest" {
		if delerr := handleDelete(qMessage); ErrorHandler(delerr) {
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

func handleDelete(qMessage QMessage) error {
	pm := PublishMessage{
		RequestID:    qMessage.MessageBody.RequestID,
		RequestType:  qMessage.MessageBody.RequestType,
		UserUUID:     qMessage.MessageBody.UserUUID,
		ServiceName:  "wtc",
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

	if err := wtc.DeleteLastPlayed(qMessage.MessageBody.UserUUID); ErrorHandler(err) {
		pm.ErrorMessage = fmt.Sprintf("error deleting last played %v\n", err.Error())
		pm.Success = false
		if perr := Publish(pm); ErrorHandler(perr) {
			return perr
		}
		return err
	}
	if err := wtc.DeletePlayedPositions(qMessage.MessageBody.UserUUID); ErrorHandler(err) {
		pm.ErrorMessage = fmt.Sprintf("error deleting play positions %v\n", err.Error())
		pm.Success = false
		if perr := Publish(pm); ErrorHandler(perr) {
			return perr
		}
		return err
	}
	if perr := Publish(pm); ErrorHandler(perr) {
		return perr
	}
	return nil
}

func handleDownload(qMessage QMessage) error {
	pm := PublishMessage{
		RequestID:    qMessage.MessageBody.RequestID,
		RequestType:  qMessage.MessageBody.RequestType,
		UserUUID:     qMessage.MessageBody.UserUUID,
		ServiceName:  "wtc",
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

	if allInfo, err := GetAllItems(qMessage.MessageBody.UserUUID); ErrorHandler(err) {
		return err
	} else {
		filePath := filepath.Join(getGdprBasePath(), qMessage.MessageBody.RequestID, "wtc.json")
		if err := AddFileToS3(filePath, JsonifyAllInfo(allInfo)); ErrorHandler(err) {
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
	return nil
}

func deleteMessage(message *sqs.Message) {
	if _, err := getSQSSession().DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(getListenQueueUrl()),
		ReceiptHandle: message.ReceiptHandle,
	}); err != nil {
		ErrorHandler(err)
		fmt.Printf("Error removing message from queue: %v\n", err)
	}
}

func pollSqs(chn chan<- *sqs.Message) {
	for {
		output, err := getSQSSession().ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(getListenQueueUrl()),
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
