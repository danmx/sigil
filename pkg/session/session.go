package session

import (
	"encoding/json"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/danmx/sigil/pkg/utils"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	log "github.com/sirupsen/logrus"
)

// StartInput struct contains all input data
type StartInput struct {
	Target     *string
	TargetType *string
	AWSSession *session.Session
	AWSProfile *string
}

// StartSSHInput struct contains all input data
type StartSSHInput struct {
	InstanceID *string
	PortNumber *int
	AWSSession *session.Session
	AWSProfile *string
}

// Start will start a session in chosen EC2 instance
func Start(input *StartInput) error {
	instance, err := utils.GetInstance(input.AWSSession, *input.TargetType, *input.Target)
	if err != nil {
		return err
	}
	target := *instance.InstanceId
	log.WithField("target instance id", target).Debug("Checking the target instance ID")
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

	startSessionInputJSON, err := json.Marshal(startSessionInput)
	if err != nil {
		return err
	}

	endpoint := ssmClient.Client.Endpoint

	if err = runSessionPluginManager(string(payload), *input.AWSSession.Config.Region, *input.AWSProfile, string(startSessionInputJSON), endpoint); err != nil {
		return err
	}

	return nil
}

// StartSSH will start a ssh proxy session in chosen EC2 instance
func StartSSH(input *StartSSHInput) error {
	ssmClient := ssm.New(input.AWSSession)
	parameters := map[string][]*string{
		"portNumber": []*string{aws.String(strconv.Itoa(*input.PortNumber))},
	}
	startSessionInput := &ssm.StartSessionInput{
		Parameters:   parameters,
		Target:       input.InstanceID,
		DocumentName: aws.String("AWS-StartSSHSession"),
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

	startSessionInputJSON, err := json.Marshal(startSessionInput)
	if err != nil {
		return err
	}

	endpoint := ssmClient.Client.Endpoint

	if err = runSessionPluginManager(string(payload), *input.AWSSession.Config.Region, *input.AWSProfile, string(startSessionInputJSON), endpoint); err != nil {
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

func runSessionPluginManager(payload, region, profile, input, endpoint string) error {
	log.WithFields(log.Fields{
		"payload":  payload,
		"region":   region,
		"profile":  profile,
		"input":    input,
		"endpoint": endpoint,
	}).Debug("Inspect session-manager-plugin args")
	// https://github.com/aws/aws-cli/blob/5f16b26/awscli/customizations/sessionmanager.py#L83-L89
	shell := exec.Command("session-manager-plugin", payload, region, "StartSession", profile, input, endpoint)
	shell.Stdout = os.Stdout
	shell.Stdin = os.Stdin
	shell.Stderr = os.Stderr
	utils.IgnoreUserEnteredSignals()
	// This allows to gracefully close the process and execute all defers
	signal.Ignore(syscall.SIGHUP)
	defer signal.Reset()
	err := shell.Run()
	if err != nil {
		return err
	}
	return nil
}
