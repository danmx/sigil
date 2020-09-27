package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	log "github.com/sirupsen/logrus"
)

// StartSession will start a session for a chosed node
func (p *Provider) StartSession(targetType, target string) error {
	instance, err := p.getInstance(targetType, target)
	if err != nil {
		log.WithFields(log.Fields{
			"targetType": targetType,
			"target":     target,
		}).Error("failed getting the instance")
		return err
	}
	log.WithField("target instance id", *instance.InstanceId).Debug("Checking the target instance ID")
	startSessionInput := &ssm.StartSessionInput{
		Target: instance.InstanceId,
	}
	ssmClient := ssm.New(p.session)
	output, err := ssmClient.StartSession(startSessionInput)
	if err != nil {
		return err
	}

	defer func() {
		if err = terminateSession(ssmClient, *output.SessionId); err != nil {
			err = fmt.Errorf("failed terminating the session (it could be already terminated): %e", err)
			log.Warn(err)
		}
	}()

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

	return runSessionPluginManager(string(payload), *p.session.Config.Region, p.awsProfile, string(startSessionInputJSON), ssmClient.Client.Endpoint)
}

// TerminateSession will close chosed active session
func (p *Provider) TerminateSession(sessionID string) error {
	return terminateSession(ssm.New(p.session), sessionID)
}

func terminateSession(ssmClient ssmiface.SSMAPI, sessionID string) error {
	_, err := ssmClient.TerminateSession(&ssm.TerminateSessionInput{
		SessionId: &sessionID,
	})
	if err != nil {
		log.WithFields(log.Fields{"sessionID": sessionID}).Warn(err)
		return err
	}
	return nil
}
