package ssh

import (
	"fmt"
	"io/ioutil"
	"crypto/rsa"
	"crypto/rand"
	"strings"
	"encoding/pem"
	"crypto/x509"
	"os"

	"github.com/danmx/sigil/pkg/utils"
	remoteSession "github.com/danmx/sigil/pkg/session"

	"golang.org/x/crypto/ssh"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	log "github.com/sirupsen/logrus"
)

// StartInput struct contains all input data
type StartInput struct {
	Target       *string
	TargetType   *string
	PortNumber   *int
	PublicKey    *string
	OSUser       *string
	GenKeyPair   *bool
	AWSSession   *session.Session
	AWSProfile   *string
}

// Start will start ssh session
func Start(input *StartInput) error {
	pubKey := *input.PublicKey
	if *input.GenKeyPair {
		privKeyBlob, err := rsa.GenerateKey(rand.Reader, 4092)
		if err != nil {
			return err
		}
		pubKeyBlob := privKeyBlob.PublicKey
		if err = savePublicPEMKey(pubKey, &pubKeyBlob); err != nil {
			return err
		}
		privKey := strings.TrimSuffix(pubKey, ".pub")
		if err = savePrivPEMKey(privKey, privKeyBlob); err != nil {
			return err
		}
		// Remove temporary keys
		defer deleteTempKey(pubKey)
		defer deleteTempKey(privKey)
	}
	instance, err := utils.GetInstance(input.AWSSession, *input.TargetType, *input.Target)
	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"pubKey": pubKey,
	}).Debug("Checking the path of a public key")

	if pubKey != "" {
		pubKeyString := ""
		dat, err := ioutil.ReadFile(pubKey)
		if err != nil {
			return err
		}
		pubKeyString = string(dat)

		log.WithFields(log.Fields{
			"SSHPublicKey": pubKeyString,
			"InstanceOSUser": *input.OSUser,
			"InstanceId": *instance.InstanceId,
			"AvailabilityZone": *instance.Placement.AvailabilityZone,
		}).Debug("SendSSHPublicKey")
		
		svc := ec2instanceconnect.New(input.AWSSession)
		out, err := svc.SendSSHPublicKey(&ec2instanceconnect.SendSSHPublicKeyInput{
			AvailabilityZone: instance.Placement.AvailabilityZone,
			InstanceId: instance.InstanceId,
			InstanceOSUser: input.OSUser,
			SSHPublicKey: &pubKeyString,
		})
		if err != nil {
			return err
		}
		if !*out.Success {
			return fmt.Errorf("SendSSHPublicKey has not succeeded. RequestID: %s", *out.RequestId)
		}
	}

	log.WithFields(log.Fields{
		"InstanceID": *instance.InstanceId,
		"PortNumber": *input.PortNumber,
		"AWSSession": *input.AWSSession,
		"AWSProfile": *input.AWSProfile,
		"RemoveTempKeyPair": *input.GenKeyPair,
		"PublicKeyPath": pubKey,
	}).Debug("StartSSHInput")

	err = remoteSession.StartSSH(&remoteSession.StartSSHInput{
		InstanceID:	instance.InstanceId,
		PortNumber:	input.PortNumber,
		AWSSession: input.AWSSession,
		AWSProfile: input.AWSProfile,
	})
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

func savePrivPEMKey(fileName string, key *rsa.PrivateKey) error {
	var privateKey = &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Headers: nil,
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	if err := ioutil.WriteFile(fileName, pem.EncodeToMemory(privateKey), 0600); err != nil {
		return err
	}
	return nil
}

func savePublicPEMKey(fileName string, pubkey *rsa.PublicKey) error {
	pub, err := ssh.NewPublicKey(pubkey)
    if err != nil {
        return err
    }
    if err := ioutil.WriteFile(fileName, ssh.MarshalAuthorizedKey(pub), 0655); err != nil {
		return err
	}
	return nil
}

func deleteTempKey(keyPath string) {
	stat, err := os.Stat(keyPath)
	log.WithFields(log.Fields{
		"stat": stat,
		"err": err,
	}).Debug("Checking if key exist")
	if err == nil {
		if err = os.Remove(keyPath); err != nil {
			log.Error(err)
		}
	}
}
