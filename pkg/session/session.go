package session

import (
	"github.com/danmx/sigil/pkg/aws"

	log "github.com/sirupsen/logrus"
)

// Session wraps methods used from the pkg/session package
type Session interface {
	Start(input *StartInput) error
}

// StartInput struct contains all input data
type StartInput struct {
	Target     *string
	TargetType *string
	MFAToken   *string
	Region     *string
	Profile    *string
	Trace      *bool
}

// Start will start a session in chosen instance
func Start(input *StartInput) error {
	return input.start(new(aws.Provider))
}

func (input *StartInput) start(provider aws.CloudInstances) error {
	err := provider.NewWithConfig(&aws.Config{
		Region:   *input.Region,
		Profile:  *input.Profile,
		MFAToken: *input.MFAToken,
		Trace:    *input.Trace,
	})
	if err != nil {
		log.Error("Failed to generate new provider")
		return err
	}
	if err := provider.StartSession(*input.TargetType, *input.Target); err != nil {
		log.WithFields(log.Fields{
			"target":     *input.Target,
			"targetType": *input.TargetType,
		}).Error("Failed to start a session")
		return err
	}
	return nil
}
