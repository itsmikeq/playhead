package queues

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

func (q *Queue) StartCmsListener(ctx *Context) {
	chnMessages := make(chan *sqs.Message, int64(q.Config.SqsMaxMessages))
	go q.pollCmsSqs(chnMessages)

	fmt.Printf("Listening on stack queue: %s\n", q.Config.CmsQueueUrl)

	go func() {
		for message := range chnMessages {
			if err := q.handleCMSMessage(message); err == nil {
				q.deleteQMessage(message, q.Config.CmsQueueUrl)
			} else {
				ErrorLogger(err)
			}
		}
	}()
}

func (q *Queue) pollCmsSqs(chn chan<- *sqs.Message) {
	for {
		output, err := q.getSQSSession().ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(q.Config.CmsQueueUrl),
			MaxNumberOfMessages: aws.Int64(int64(q.Config.SqsMaxMessages)),
			WaitTimeSeconds:     aws.Int64(int64(q.Config.TimeWaitSeconds)),
		})

		if err != nil {
			logrus.Error("failed to fetch sqs message: ", err)
			panic(err)
		}

		for _, message := range output.Messages {
			chn <- message
		}

	}

}

func (q *Queue) handleCMSMessage(message *sqs.Message) error {
	// TODO: handle the messages/update the series and episode info
	fmt.Printf("Message: %+v\n", message)
	return nil
}
