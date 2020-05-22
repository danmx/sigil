package aws

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/ssm"
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
	ssmClient := ssm.New(p.session)
	startSessionInput := &ssm.StartSessionInput{
		Target: instance.InstanceId,
	}
	output, err := ssmClient.StartSession(startSessionInput)
	if err != nil {
		return err
	}

	defer func() {
		if err = p.TerminateSession(*output.SessionId); err != nil {
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

	endpoint := ssmClient.Client.Endpoint

	// returns err
	return runSessionPluginManager(string(payload), *p.session.Config.Region, p.awsProfile, string(startSessionInputJSON), endpoint)
}

// TerminateSession will close chosed active session
func (p *Provider) TerminateSession(sessionID string) error {
	client := ssm.New(p.session)
	_, err := client.TerminateSession(&ssm.TerminateSessionInput{
		SessionId: &sessionID,
	})
	if err != nil {
		log.WithFields(log.Fields{"sessionID": sessionID}).Warn(err)
		return err
	}
	return nil
}
