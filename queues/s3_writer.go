package queues

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"net/http"
)

// func Write(fileName string, content string) error {
// 	if err := AddFileToS3(fileName, content); ErrorHandler(err) {
// 		log.Fatal(err)
// 		return err
// 	}
// 	return nil
// }

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func AddFileToS3(fileName string, content string) error {
	// Get file size and read the file content into a buffer
	var size = len(content)
	var b bytes.Buffer
	if _, err := fmt.Fprintf(&b, content); err != nil {
		return err
	}
	// fmt.Printf("Writing:\n%v\n", b.String())

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	if _, err := s3.New(getSession()).PutObject(&s3.PutObjectInput{
		Bucket:             aws.String(getGdprBucket()),
		Key:                aws.String(fileName),
		ACL:                aws.String("private"),
		Body:               bytes.NewReader(b.Bytes()),
		ContentLength:      aws.Int64(int64(size)),
		ContentType:        aws.String(http.DetectContentType(b.Bytes())),
		ContentDisposition: aws.String("application/octet-stream"),
	}); ErrorHandler(err) {
		return err
	}
	return nil
}
