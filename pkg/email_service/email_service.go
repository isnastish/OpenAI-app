package emailservice

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	_ "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

type AWSEmailService struct {
	client *ses.Client
}

type Recipient struct {
	toEmails []string
	ccEmails []string
}

func NewAWSEmailService() (*AWSEmailService, error) {
	// Load the shared AWS configuration (~/.aws/config)
	awsConfig, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to load aws shared config, %v", err)
	}

	return &AWSEmailService{
		client: ses.NewFromConfig(awsConfig),
	}, nil
}

func (s *AWSEmailService) SendEmail(messageBody string, subject string, fromEmail string, recipient Recipient) error {
	// s.client.SendEmail()
	return nil
}
