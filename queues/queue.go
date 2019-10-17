package queues

import (
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

// Private

func (q *Queue) getSession() *session.Session {
	// sess = session.Must(session.NewSessionWithOptions(session.Options{
	// 	AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	// 	SharedConfigState:       session.SharedConfigEnable,
	// 	Config: aws.Config{
	// 		Region: aws.String(getAwsRegion),
	// 		CredentialsChainVerboseErrors: aws.Bool(true),
	// 	},
	// }))
	if sess, err := session.NewSession(&aws.Config{Region: aws.String(string(q.Config.AwsRegion))}); !ErrorHandler(err) {
		return sess
	} else {
		return nil
	}
}

func (q *Queue) getSQSSession() *sqs.SQS {
	sqsSession := sqs.New(q.getSession())
	return sqsSession
}
