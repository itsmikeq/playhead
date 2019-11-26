package queues

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
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

func (q *Queue) deleteQMessage(message *sqs.Message, qUrl string) {
	if _, err := q.getSQSSession().DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      aws.String(qUrl),
		ReceiptHandle: message.ReceiptHandle,
	}); err != nil {
		ErrorHandler(err)
		logrus.Errorf("Error removing message from queue: %v\n", err)
	} else {
		logrus.Debug(fmt.Sprintf("Deleted +%v\n", message))
	}
}

func (q *Queue) getSession() (*session.Session, error) {
	// sess = session.Must(session.NewSessionWithOptions(session.Options{
	// 	AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	// 	SharedConfigState:       session.SharedConfigEnable,
	// 	Config: aws.Config{
	// 		Region: aws.String(getAwsRegion),
	// 		CredentialsChainVerboseErrors: aws.Bool(true),
	// 	},
	// }))
	if len(q.Config.CmsQueueUrl) > 0 {
		fmt.Println("Using cms q ", q.Config.CmsQueueUrl)
		if sess, err := session.NewSession(&aws.Config{LogLevel: aws.LogLevel(3), DisableSSL: aws.Bool(true), Region: aws.String(q.Config.AwsRegion)}); err == nil {
			return sess, nil
		} else {
			logrus.Errorf(fmt.Sprintf("Error +%v\n", err))
			return nil, err
		}
	} else {
		if sess, err := session.NewSession(&aws.Config{Region: aws.String(q.Config.AwsRegion)}); !ErrorHandler(err) {
			return sess, nil
		} else {
			return nil, err
		}
	}
}

func (q *Queue) getSQSSession() *sqs.SQS {
	if sess, err := q.getSession(); err != nil {
		panic(err)
	} else {
		return sqs.New(sess)
	}
}
