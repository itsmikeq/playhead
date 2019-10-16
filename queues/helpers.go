package queues

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"runtime"
)

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

// true if error exists
func ErrorHandler(err error) (b bool) {
	b = false
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(1)

		logrus.Error(fmt.Sprintf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
		logrus.Error(fmt.Sprintf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
		b = true
	}
	return b
}

func ErrorLogger(err error) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(1)

		logrus.Error(fmt.Sprintf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
		logrus.Error(fmt.Sprintf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err))
	}
}

func exitErrorf(msg string, args ...interface{}) {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, msg)
	fmt.Printf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, msg)
	// fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
