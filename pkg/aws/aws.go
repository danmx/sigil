package aws

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	sigilOS "github.com/danmx/sigil/pkg/os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

// Provider contains necessary components like the session
type Provider struct {
	filters    Filters
	session    *session.Session
	awsProfile string
}

// Config contains provider's configuration
type Config struct {
	Filters  Filters
	Region   string
	Profile  string
	MFAToken string
}

// Filters grouped per type
type Filters struct {
	Instance InstanceFilters `mapstructure:"instance"`
	Session  SessionFilters  `mapstructure:"session"`
}

// Instance contain information about the EC2 instance
type Instance struct {
	Hostname       string `json:"hostname" yaml:"hostname"`
	IPAddress      string `json:"ip_address" yaml:"ip_address"`
	ID             string `json:"id" yaml:"id"`
	PrivateDNSName string `json:"private_dns_name" yaml:"private_dns_name"`
	Name           string `json:"name" yaml:"name"`
	OSName         string `json:"os_name" yaml:"os_name"`
	OSType         string `json:"os_type" yaml:"os_type"`
	OSVersion      string `json:"os_version" yaml:"os_version"`
}

// InstanceFilters contain all types of filters used to limit results
type InstanceFilters struct {
	IDs  []string    `json:"ids" mapstructure:"ids,remain"`
	Tags []TagValues `json:"tags" mapstructure:"tags,remain"`
}

// TagValues contain list of values for specific key
type TagValues struct {
	Key    string   `json:"key" mapstructure:"key,remain"`
	Values []string `json:"values" mapstructure:"values,remain"`
}

// Session contains information about SSM sessions
type Session struct {
	SessionID string `json:"session_id" yaml:"session_id"`
	Target    string `json:"target" yaml:"target"`
	Status    string `json:"status" yaml:"status"`
	StartDate string `json:"start_date" yaml:"start_date"`
	Owner     string `json:"owner" yaml:"owner"`
}

// SessionFilters for SSM sessions
type SessionFilters struct {
	After  string `json:"after" mapstructure:"after,remain"`
	Before string `json:"before" mapstructure:"before,remain"`
	Target string `json:"target" mapstructure:"target,remain"`
	Owner  string `json:"owner" mapstructure:"owner,remain"`
}

const (
	maxResults int64  = 50
	pluginName string = "session-manager-plugin"
	// TargetTypeInstanceID points to an instance ID type
	TargetTypeInstanceID = "instance-id"
	// TargetTypePrivateDNS points to a private DNS type
	TargetTypePrivateDNS = "private-dns"
	// TargetTypeName points to a name type
	TargetTypeName = "name"
)

// NewWithConfig will generate an AWS Provider with given configuration
func (p *Provider) NewWithConfig(c *Config) error {
	options := session.Options{
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: mfaTokenProvider(c.MFAToken),
		Profile:                 c.Profile,
	}
	awsConfig := aws.NewConfig()
	awsConfig.Region = aws.String(c.Region)
	options.Config = *awsConfig
	sess, err := session.NewSessionWithOptions(options)
	if err != nil {
		log.WithFields(log.Fields{
			"err":     err,
			"options": options,
		}).Error("Failed starting a session")
		return err
	}
	p.session = sess
	p.awsProfile = c.Profile
	p.filters = c.Filters
	return nil
}

// VerifyDependencies will check all necessary dependencies
func VerifyDependencies() error {
	o, err := exec.LookPath(pluginName)
	if err != nil {
		err = fmt.Errorf("required plugin not found: %s", err)
		log.Error(err)
		return err
	}
	log.WithFields(log.Fields{
		"plugin": pluginName,
		"path":   o,
	}).Debugf("%s is installed successfully in %s\n", pluginName, o)
	return nil
}

func (p *Provider) getInstance(targetType, target string) (*ec2.Instance, error) {
	if target == "" {
		err := errors.New("no target")
		log.WithFields(log.Fields{
			"target": target,
		}).Error(err)
		return nil, err
	}
	var filters []*ec2.Filter
	switch targetType {
	case TargetTypeInstanceID:
		filters = []*ec2.Filter{
			{
				Name:   aws.String("instance-id"),
				Values: []*string{&target},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		}
	case TargetTypePrivateDNS:
		filters = []*ec2.Filter{
			{
				Name:   aws.String("private-dns-name"),
				Values: []*string{&target},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		}
	case TargetTypeName:
		filters = []*ec2.Filter{
			{
				Name:   aws.String("tag:Name"),
				Values: []*string{&target},
			},
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		}
	default:
		log.WithFields(log.Fields{
			"target":     target,
			"targetType": targetType,
		}).Error("unsupported target type")
		return nil, fmt.Errorf("unsupported target type: %s", targetType)
	}
	instance, err := p.getFirstInstance(filters)
	if err != nil {
		log.WithFields(log.Fields{
			"filters": filters,
		}).Error("failed getting the first instance")
		return nil, err
	}
	if instance == nil {
		err := fmt.Errorf("no instance that matches the target (%s) and the type (%s)", target, targetType)
		log.WithFields(log.Fields{
			"targetType": targetType,
			"taget":      target,
		}).Info(err)
		return nil, err
	}
	return instance, nil
}

func (p *Provider) getFirstInstance(filters []*ec2.Filter) (*ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters:    filters,
		MaxResults: aws.Int64(maxResults),
	}
	var target *ec2.Instance
	ec2Client := ec2.New(p.session)
	err := ec2Client.DescribeInstancesPages(input,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, reservation := range page.Reservations {
				for _, instance := range reservation.Instances {
					target = instance
					// Escape the function
					return false
				}
			}
			return !lastPage
		})
	if err != nil {
		log.WithFields(log.Fields{
			"filters":                filters,
			"DescribeInstancesInput": input,
		}).Error("failed DescribeInstancesPages")
		return nil, err
	}
	return target, nil
}

func mfaTokenProvider(token string) func() (string, error) {
	log.WithFields(log.Fields{
		"token": token,
	}).Debug("Get MFA Token Provider")
	if token == "" {
		return stscreds.StdinTokenProvider
	}
	return func() (string, error) {
		return token, nil
	}
}

func runSessionPluginManager(payloadJSON, region, profile, inputJSON, endpoint string) error {
	log.WithFields(log.Fields{
		"payload":  payloadJSON,
		"region":   region,
		"profile":  profile,
		"input":    inputJSON,
		"endpoint": endpoint,
	}).Debug("Inspect session-manager-plugin args")
	// TODO allowing logging
	// https://docs.aws.amazon.com/systems-manager/latest/userguide/session-manager-working-with-install-plugin.html#configure-logs-linux
	// https://github.com/aws/aws-cli/blob/5f16b26/awscli/customizations/sessionmanager.py#L83-L89
	shell := exec.Command(pluginName, payloadJSON, region, "StartSession", profile, inputJSON, endpoint)
	shell.Stdout = os.Stdout
	shell.Stdin = os.Stdin
	shell.Stderr = os.Stderr
	sigilOS.IgnoreUserEnteredSignals()
	// This allows to gracefully close the process and execute all defers
	signal.Ignore(syscall.SIGHUP)
	defer signal.Reset()
	err := shell.Run()
	if err != nil {
		return err
	}
	return nil
}
