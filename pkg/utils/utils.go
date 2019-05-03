package utils

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

// StartAWSSession will return AWS Session
func StartAWSSession(region, profile, mfa string) *session.Session {
	options := session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: awsMFATokenProvider(mfa),
	}
	if profile != "" {
		options.Profile = profile
	}
	awsConfig := aws.NewConfig()
	if region != "" {
		awsConfig.Region = &region
	}
	options.Config = *awsConfig
	sess := session.Must(session.NewSessionWithOptions(options))
	return sess
}

// Helper functions

func awsMFATokenProvider(token string) func() (string, error) {
	log.WithFields(log.Fields{
		"token": token,
	}).Debug("Get MFA Token Provider")
	if token == "" {
		return stscreds.StdinTokenProvider
	}
	return func() (string, error) {
		return token, nil
	}
}
