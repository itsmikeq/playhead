package queues

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
)

// used to publish messages to SNS
type PublishMessage struct {
	UserUUID     string `json:"user_uuid" binding:"required"`
	S3Location   string `json:"s3_url" binding:"required"`
	RequestID    string `json:"request_id" binding:"required"`
	RequestType  string `json:"request_type" binding:"required"`
	ServiceName  string `json:"service_name" binding:"required"`
	ErrorMessage string `json:"error_message"`
	Success      bool   `json:"success" binding:"required"` // "true" || "false"
}

func Publish(message PublishMessage) error {
	if messageJSON, err := json.Marshal(message); ErrorHandler(err) {
		log.Printf("Error marshaling message: %v\n", err)
	} else {
		sess := getSQSSession()
		// message = {
		//  request_id: "<request_id_from_original_message>",
		//  request_type: "UserDataDeleteRequest",
		//  user_uuid: "<user_uuid_from_original_message>",
		//  service_name: <client_service>,
		//  success: true,
		//  error_message: "<custom message>"
		// }
		input := &sqs.SendMessageInput{
			MessageBody: aws.String(string(messageJSON)),
			QueueUrl:    aws.String(getCallbackUrl()),
		}
		// fmt.Printf("Delivering message %v\n", input)
		if _, err1 := sess.SendMessage(input); ErrorHandler(err1) { // Call to puclish the message
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err1.Error())
			return err1
		}
	}
	return nil
}
