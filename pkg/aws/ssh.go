package aws

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface"
	"github.com/aws/aws-sdk-go/service/ssm"
	log "github.com/sirupsen/logrus"
)

// StartSSH will start a ssh proxy session for a chosed node
func (p *Provider) StartSSH(targetType, target, osUser string, portNumber uint64, publicKey []byte) error {
	instance, err := p.getInstance(targetType, target)
	if err != nil {
		return err
	}

	if len(publicKey) > 0 {
		svc := ec2instanceconnect.New(p.session)
		err = uploadPublicKey(svc, publicKey, osUser, *instance.InstanceId, *instance.Placement.AvailabilityZone)
		if err != nil {
			return err
		}
	}

	ssmClient := ssm.New(p.session)
	parameters := map[string][]*string{
		"portNumber": {aws.String(strconv.FormatUint(portNumber, 10))},
	}
	startSessionInput := &ssm.StartSessionInput{
		Parameters:   parameters,
		Target:       instance.InstanceId,
		DocumentName: aws.String("AWS-StartSSHSession"),
	}
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

	endpoint := ssmClient.Client.Endpoint

	return runSessionPluginManager(string(payload), *p.session.Config.Region, p.awsProfile, string(startSessionInputJSON), endpoint)
}

func uploadPublicKey(client ec2instanceconnectiface.EC2InstanceConnectAPI, publicKey []byte, osUser, instanceID, availabilityZone string) error {
	pubKey := string(publicKey)
	log.WithFields(log.Fields{
		"SSHPublicKey":     pubKey,
		"InstanceOSUser":   osUser,
		"InstanceId":       instanceID,
		"AvailabilityZone": availabilityZone,
	}).Debug("SendSSHPublicKey")

	out, err := client.SendSSHPublicKey(&ec2instanceconnect.SendSSHPublicKeyInput{
		AvailabilityZone: &availabilityZone,
		InstanceId:       &instanceID,
		InstanceOSUser:   &osUser,
		SSHPublicKey:     &pubKey,
	})
	if err != nil {
		log.WithFields(log.Fields{
			"AvailabilityZone": availabilityZone,
			"InstanceID":       instanceID,
			"InstanceOSUser":   osUser,
			"SSHPublicKey":     pubKey,
			"error":            err,
		}).Error("failed SendSSHPublicKey")
		return err
	}
	if !*out.Success {
		return fmt.Errorf("failed SendSSHPublicKey. RequestID: %s", *out.RequestId)
	}
	return nil
}
