package session

import (
	"encoding/json"
	"os"
	"os/exec"
	"os/signal"

	"github.com/danmx/sigil/pkg/utils"

	"github.com/aws/aws-sdk-go/aws/session"
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
