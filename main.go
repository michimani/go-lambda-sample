package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Response struct {
	Message    string   `json:"message"`
	BucketList []string `json:"bucket_list"`
}

func listBuckets(c *s3.S3) ([]string, error) {
	out, err := c.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	var list []string
	for _, b := range out.Buckets {
		list = append(list, aws.StringValue(b.Name))
	}

	return list, nil
}

func handleRequest() (Response, error) {
	if os.Getenv("AWS_DEFAULT_REGION") == "" {
		return Response{}, errors.New("'AWS_DEFAULT_REGION' is required.")
	}
	region := os.Getenv("AWS_DEFAULT_REGION")
	s3sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	s3client := s3.New(s3sess)

	buckets, err := listBuckets(s3client)
	if err != nil {
		return Response{Message: fmt.Sprintf("An error occurred: %s", err.Error())}, nil
	}

	return Response{
		Message:    "Success!",
		BucketList: buckets,
	}, nil
}

func main() {
	runtime.Start(handleRequest)
}
