package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"strings"

	"github.com/danmx/sigil/pkg/aws"

	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

// SSH wraps methods used from the pkg/ssh package
type SSH interface {
	Start(input *StartInput) error
}

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
	Trace      *bool
}

// Start will start ssh session
func Start(input *StartInput) error {
	return input.start(new(aws.Provider))
}

func (input *StartInput) start(provider aws.CloudSSH) error {
	err := provider.NewWithConfig(&aws.Config{
		Region:   *input.Region,
		Profile:  *input.Profile,
		MFAToken: *input.MFAToken,
		Trace:    *input.Trace,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	pubKey := *input.PublicKey
	if *input.GenKeyPair {
		const rsaKeySize = 4092
		privKeyBlob, errKey := rsa.GenerateKey(rand.Reader, rsaKeySize)
		if errKey != nil {
			return errKey
		}
		pubKeyBlob := privKeyBlob.PublicKey
		if errPubPEM := savePublicPEMKey(pubKey, &pubKeyBlob); errPubPEM != nil {
			return errPubPEM
		}
		defer func() {
			if err = deleteTempKey(pubKey); err != nil {
				log.Error(err)
			}
		}()
		privKey := strings.TrimSuffix(pubKey, ".pub")
		if errPrivPEM := savePrivPEMKey(privKey, privKeyBlob); errPrivPEM != nil {
			return errPrivPEM
		}
		defer func() {
			if err = deleteTempKey(privKey); err != nil {
				log.Error(err)
			}
		}()
	}

	pubKeyData := []byte{}
	if pubKey != "" {
		pubKeyData, err = os.ReadFile(pubKey)
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

	// returns err
	return provider.StartSSH(*input.TargetType, *input.Target, *input.OSUser, *input.PortNumber, pubKeyData)
}

// Helper functions

func savePrivPEMKey(fileName string, key *rsa.PrivateKey) error {
	privateKey := &pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   x509.MarshalPKCS1PrivateKey(key),
	}

	// returns err
	return os.WriteFile(fileName, pem.EncodeToMemory(privateKey), 0o600) //nolint:gomnd // Linux file permissions
}

func savePublicPEMKey(fileName string, pubkey *rsa.PublicKey) error {
	pub, err := ssh.NewPublicKey(pubkey)
	if err != nil {
		return err
	}
	// returns err
	return os.WriteFile(fileName, ssh.MarshalAuthorizedKey(pub), 0o600) //nolint:gomnd // Linux file permissions
}

func deleteTempKey(keyPath string) error {
	stat, err := os.Stat(keyPath)
	log.WithFields(log.Fields{
		"stat": stat,
		"err":  err,
	}).Debug("Checking if key exist")
	if err == nil {
		if errRm := os.Remove(keyPath); errRm != nil {
			return errRm
		}
	}
	return err
}
