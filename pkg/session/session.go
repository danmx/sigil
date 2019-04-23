package session

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	log "github.com/sirupsen/logrus"
)

// StartInput struct contains all input data
type StartInput struct {
	Target     *string
	TargetType *string
	AWSSession *session.Session
}

// Start will start a session in chosen EC2 instance
func Start(input *StartInput) error {
	var instanceID *string
	switch *input.TargetType {
	case "instance-id":
		instanceID = input.Target
	case "priv-dns":
		id, err := getIDFromPrivDNS(input.AWSSession, input.Target)
		if err != nil {
			return err
		}
		if id == "" {
			return fmt.Errorf("no instance with private dns name: %s", *input.Target)
		}
		instanceID = &id
	case "name-tag":
		id, err := getIDFromName(input.AWSSession, input.Target)
		if err != nil {
			return err
		}
		if id == "" {
			return fmt.Errorf("no instance with name tag: %s", *input.Target)
		}
		instanceID = &id
	default:
		return fmt.Errorf("Unsupported target type: %s", *input.Target)
	}
	if *input.Target == "" {
		err := fmt.Errorf("Specify the target")
		log.WithFields(log.Fields{
			"target": *input.Target,
		}).Error(err)
		return err
	}

	ssmClient := ssm.New(input.AWSSession)
	startSessionInput := &ssm.StartSessionInput{
		Target: instanceID,
	}
	output, err := ssmClient.StartSession(startSessionInput)
	if err != nil {
		return err
	}
	defer terminateSession(ssmClient, output.SessionId)
	log.WithFields(log.Fields{
		"sessionID": *output.SessionId,
		"streamURL": *output.StreamUrl,
		"token":     *output.TokenValue,
	}).Debug("SSM Start Session Output")
	payload, err := json.Marshal(output)
	if err != nil {
		return err
	}
	shell := exec.Command("session-manager-plugin", string(payload), *input.AWSSession.Config.Region, "StartSession")
	shell.Stdout = os.Stdout
	shell.Stdin = os.Stdin
	shell.Stderr = os.Stderr
	err = shell.Run()
	if err != nil {
		return err
	}

	return nil
}

func terminateSession(client *ssm.SSM, sessionID *string) {
	_, err := client.TerminateSession(&ssm.TerminateSessionInput{
		SessionId: sessionID,
	})
	if err != nil {
		log.WithFields(log.Fields{"sessionID": *sessionID}).Error(err)
	}
}

func getIDFromPrivDNS(sess *session.Session, dnsName *string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{dnsName},
			},
		},
	}
	return getFirstInstanceID(sess, input)
}

func getIDFromName(sess *session.Session, name *string) (string, error) {
	input := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{name},
			},
		},
	}
	return getFirstInstanceID(sess, input)
}

func getFirstInstanceID(sess *session.Session, input *ec2.DescribeInstancesInput) (string, error) {
	var instanceID string
	ec2Client := ec2.New(sess)
	err := ec2Client.DescribeInstancesPages(input,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, reservation := range page.Reservations {
				for _, instance := range reservation.Instances {
					instanceID = *instance.InstanceId
					// Escape the function
					return false
				}
			}
			return !lastPage
		})
	if err != nil {
		return "", err
	}
	return instanceID, nil
}
