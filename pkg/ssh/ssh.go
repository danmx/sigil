package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"strings"

	"github.com/danmx/sigil/pkg/aws"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// StartInput struct contains all input data
type StartInput struct {
	Target     *string
	TargetType *string
	PortNumber *uint64
	PublicKey  *string
	OSUser     *string
	GenKeyPair *bool
	MFAToken   *string
	Region     *string
	Profile    *string
}

// Start will start ssh session
func Start(input *StartInput) error {
	provider, err := aws.NewWithConfig(&aws.Config{
		Region:   *input.Region,
		Profile:  *input.Profile,
		MFAToken: *input.MFAToken,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	pubKey := *input.PublicKey
	if *input.GenKeyPair {
		privKeyBlob, errKey := rsa.GenerateKey(rand.Reader, 4092)
		if errKey != nil {
			return err
		}
		pubKeyBlob := privKeyBlob.PublicKey
		if errPubPEM := savePublicPEMKey(pubKey, &pubKeyBlob); errPubPEM != nil {
			return errPubPEM
		}
		privKey := strings.TrimSuffix(pubKey, ".pub")
		if errPrivPEM := savePrivPEMKey(privKey, privKeyBlob); errPrivPEM != nil {
			return errPrivPEM
		}
		// Remove temporary keys
		defer deleteTempKey(pubKey)
		defer deleteTempKey(privKey)
	}

	pubKeyData := []byte{}
	if pubKey != "" {
		pubKeyData, err = ioutil.ReadFile(pubKey)
		if err != nil {
			return err
		}
	}

	log.WithFields(log.Fields{
		"targetType":    *input.TargetType,
		"PortNumber":    *input.PortNumber,
		"target":        *input.Target,
		"OSUser":        *input.OSUser,
		"pubKeyData":    string(pubKeyData),
		"PublicKeyPath": pubKey,
	}).Debug("StartSSHInput")

	err = provider.StartSSH(*input.TargetType, *input.Target, *input.OSUser, *input.PortNumber, pubKeyData)
	if err != nil {
		return err
	}

	return nil
}

// Helper functions

func savePrivPEMKey(fileName string, key *rsa.PrivateKey) error {
	var privateKey = &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(key),
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
		"err":  err,
	}).Debug("Checking if key exist")
	if err == nil {
		if err = os.Remove(keyPath); err != nil {
			log.Error(err)
		}
	}
}
