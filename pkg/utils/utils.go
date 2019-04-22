package utils

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	log "github.com/sirupsen/logrus"
)

// StringMapStringType implements cli.GenericType
type StringMapStringType struct {
	Map map[string]string
}

// Set Map value
func (m *StringMapStringType) Set(value string) error {
	tagsMap, err := stringTagsToMap(value)
	if err != nil {
		return err
	}
	m.Map = tagsMap
	return nil
}

// String return string representation of Map value
func (m *StringMapStringType) String() string {
	list := make([]string, 0, len(m.Map))
	for key, value := range m.Map {
		list = append(list, key+"="+value)
	}
	return strings.Join(list, ",")
}

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

func stringTagsToMap(value string) (map[string]string, error) {
	tagsMap := make(map[string]string)
	keyValuePairs := strings.Split(value, ",")
	for _, pair := range keyValuePairs {
		splittedPair := strings.Split(pair, "=")
		if len(splittedPair) != 2 {
			log.WithFields(log.Fields{
				"keyValuePairs": keyValuePairs,
				"pair":          pair,
				"splittedPair":  splittedPair,
			}).Error("wrong format of a key-value pair")
			return nil, fmt.Errorf("wrong format of a key-value pair: %s", pair)
		}
		tagsMap[splittedPair[0]] = splittedPair[1]
	}
	log.WithFields(log.Fields{
		"Tags": tagsMap,
	}).Debug("Tags Map")
	return tagsMap, nil
}

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
