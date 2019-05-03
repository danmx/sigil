package session

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"

	"github.com/danmx/sigil/pkg/utils"

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
	var target string
	switch *input.TargetType {
	case "instance-id":
		target = *input.Target
	case "private-dns":
		id, err := getFirstInstanceID(input.AWSSession, &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("private-dns-name"),
					Values: []*string{input.Target},
				},
			},
		})
		if err != nil {
			return err
		}
		if id == "" {
			return fmt.Errorf("no instance with private dns name: %s", *input.Target)
		}
		target = id
	case "name-tag":
		id, err := getFirstInstanceID(input.AWSSession, &ec2.DescribeInstancesInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("tag:Name"),
					Values: []*string{input.Target},
				},
			},
		})
		if err != nil {
			return err
		}
		if id == "" {
			return fmt.Errorf("no instance with name tag: %s", *input.Target)
		}
		target = id
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
		Target: &target,
	}
	output, err := ssmClient.StartSession(startSessionInput)
	if err != nil {
		return err
	}
	defer TerminateSession(ssmClient, output.SessionId)
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
	utils.IgnoreUserEnteredSignals()
	defer signal.Reset()
	err = shell.Run()
	if err != nil {
		return err
	}

	return nil
}

// TerminateSession will close chosed active session
func TerminateSession(client *ssm.SSM, sessionID *string) error {
	_, err := client.TerminateSession(&ssm.TerminateSessionInput{
		SessionId: sessionID,
	})
	if err != nil {
		log.WithFields(log.Fields{"sessionID": *sessionID}).Error(err)
		return err
	}
	return nil
}

func getFirstInstanceID(sess *session.Session, input *ec2.DescribeInstancesInput) (string, error) {
	var target string
	ec2Client := ec2.New(sess)
	err := ec2Client.DescribeInstancesPages(input,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, reservation := range page.Reservations {
				for _, instance := range reservation.Instances {
					target = *instance.InstanceId
					// Escape the function
					return false
				}
			}
			return !lastPage
		})
	if err != nil {
		return "", err
	}
	return target, nil
}
