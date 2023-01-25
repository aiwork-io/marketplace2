package internal

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Mailer interface {
	Send(to, subject, body string) error
}

type SESMailer struct {
	configs *Configs
	client  *ses.SES
}

func (mailer *SESMailer) Send(to, subject, body string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(body),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String("verification@aiworkmarketplace.kubeplusplus.com"),
	}

	_, err := mailer.client.SendEmail(input)
	return err
}

func NewMailer(configs *Configs) Mailer {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(configs.Region)},
	)
	if err != nil {
		panic(err)
	}
	client := ses.New(sess)

	return &SESMailer{configs: configs, client: client}
}
