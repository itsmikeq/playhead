package queues

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"log"
	"os"
	"runtime"
)

func getSession() *session.Session {
	// sess = session.Must(session.NewSessionWithOptions(session.Options{
	// 	AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	// 	SharedConfigState:       session.SharedConfigEnable,
	// 	Config: aws.Config{
	// 		Region: aws.String(getAwsRegion),
	// 		CredentialsChainVerboseErrors: aws.Bool(true),
	// 	},
	// }))
	if sess, err := session.NewSession(&aws.Config{Region: aws.String(getAwsRegion())}); !ErrorHandler(err) {
		return sess
	} else {
		return nil
	}
}

func getSQSSession() *sqs.SQS {
	sqsSession := sqs.New(getSession())
	return sqsSession
}

func getListenQueueUrl() string {
	return os.Getenv("GDPR_QUEUE_URL")
}

func getCallbackUrl() string {
	return os.Getenv("CALLBACK_QUEUE_URL")
}

// func getCallbackTopicARN() string {
// 	return os.Getenv("CALLBACK_TOPIC_ARN")
// }

func getAwsRegion() string {
	return os.Getenv("AWS_REGION")
}

func getGdprBucket() string {
	return os.Getenv("GDPR_BUCKET")
}

func getGdprBasePath() string {
	return os.Getenv("GDPR_BASE_PATH")
}

func ErrorHandler(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, fn, line, _ := runtime.Caller(1)

		log.Printf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
		fmt.Printf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, err)
		b = true
	}
	return
}

func exitErrorf(msg string, args ...interface{}) {
	pc, fn, line, _ := runtime.Caller(1)
	log.Printf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, msg)
	fmt.Printf("[error] in %s[%s:%d] %v\n", runtime.FuncForPC(pc).Name(), fn, line, msg)
	// fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
