package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession AWS接続セッションを返す
func NewSession() (*session.Session, error) {
	creds := credentials.NewStaticCredentials(awsAccessKey, awsSecretAccessKey, "")
	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String(region),
		Endpoint:    aws.String(endpoint),
	})
	if err != nil {
		return nil, err
	}
	return sess, nil
}
