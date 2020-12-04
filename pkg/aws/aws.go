package aws

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"github.com/danmx/sigil/pkg/aws/helpers"
	logger "github.com/danmx/sigil/pkg/aws/log"
	sigilOS "github.com/danmx/sigil/pkg/os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
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
	Trace    bool
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
	execEnvVar       = "AWS_EXECUTION_ENV"
	maxResults int64 = 50
	pluginName       = "session-manager-plugin"
	// API calls retry configuration
	numMaxRetries    = 100
	minRetryDelay    = 10 * time.Millisecond
	minThrottleDelay = 500 * time.Millisecond
	maxRetryDelay    = 5 * time.Second
	maxThrottleDelay = 30 * time.Second
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
	if c.Trace {
		awsConfig.LogLevel = aws.LogLevel(aws.LogDebugWithRequestRetries)
		awsConfig.Logger = logger.NewTraceLogger()
	}
	awsConfig.Retryer = client.DefaultRetryer{
		NumMaxRetries:    numMaxRetries,
		MinRetryDelay:    minRetryDelay,
		MinThrottleDelay: minThrottleDelay,
		MaxRetryDelay:    maxRetryDelay,
		MaxThrottleDelay: maxThrottleDelay,
	}
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
	return verifyDependencies(new(helpers.Helpers))
}

func verifyDependencies(exechelpers helpers.OSExecIface) error {
	o, err := exechelpers.LookPath(pluginName)
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

// AppendUserAgent will add given suffix to HTTP client's user agent
func AppendUserAgent(suffix string) error {
	return appendUserAgent(new(helpers.Helpers), suffix)
}

func appendUserAgent(oshelpers helpers.OSIface, suffix string) error {
	value, set := oshelpers.LookupEnv(execEnvVar)
	if set {
		value += "/"
	}
	value += suffix
	return oshelpers.Setenv(execEnvVar, value)
}

func (p *Provider) getInstance(targetType, target string) (*ec2.Instance, error) {
	return getInstance(ec2.New(p.session), targetType, target)
}

func getInstance(ec2Client ec2iface.EC2API, targetType, target string) (*ec2.Instance, error) {
	if target == "" {
		err := errors.New("no target")
		log.WithFields(log.Fields{
			"target": target,
		}).Error(err)
		return nil, err
	}
	filters, err := getFilters(targetType, target)
	if err != nil {
		return nil, err
	}
	instance, err := getFirstInstance(ec2Client, filters)
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

func getFilters(targetType, target string) ([]*ec2.Filter, error) {
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
	return filters, nil
}

func getFirstInstance(ec2Client ec2iface.EC2API, filters []*ec2.Filter) (*ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		Filters:    filters,
		MaxResults: aws.Int64(maxResults),
	}
	var target *ec2.Instance
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
