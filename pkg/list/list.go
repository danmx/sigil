package list

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	remoteSession "github.com/danmx/sigil/pkg/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

// StartInput struct contains all input data
type StartInput struct {
	// Define output format
	OutputFormat *string
	AWSSession   *session.Session
	Filters      *string
	Interactive  *bool
	Type         *string
}

// Filters contain all types of filters used to limit results
type Filters struct {
	InstanceIDs []*string    `json:"instance_ids"`
	Tags        []*TagValues `json:"tags"`
	After       *string      `json:"after"`
	Before      *string      `json:"before"`
	Target      *string      `json:"target"`
	Owner       *string      `json:"owner"`
}

// TagValues contain list of values for specific key
type TagValues struct {
	Key    string    `json:"key"`
	Values []*string `json:"values"`
}

// StartOutput struct will contain all output data
type StartOutput struct {
	Instances []*Instance
	Sessions  []*Session
	// Define output format
	format     *string
	outputType *string
}

// Instance contain information about the EC2 instance
type Instance struct {
	Hostname       *string `json:"hostname" yaml:"hostname"`
	IPAddress      *string `json:"ip_address" yaml:"ip_address"`
	InstanceID     *string `json:"instance_id" yaml:"instance_id"`
	PrivateDNSName *string `json:"private_dns_name" yaml:"private_dns_name"`
	Name           *string `json:"instance_name" yaml:"instance_name"`
	OSName         *string `json:"os_name" yaml:"os_name"`
	OSType         *string `json:"os_type" yaml:"os_type"`
	OSVersion      *string `json:"os_version" yaml:"os_version"`
}

// Session contains information about the SSM sessions
type Session struct {
	SessionID *string `json:"session_id" yaml:"session_id"`
	Target    *string `json:"target_instance" yaml:"target_instance"`
	Status    *string `json:"status" yaml:"status"`
	StartDate *string `json:"start_date" yaml:"start_date"`
	Owner     *string `json:"owner" yaml:"owner"`
}

const (
	capacity      = 10
	capMultiplier = 2
)

// Start will output a ist of all available EC2 instances
func Start(input *StartInput) error {
	switch *input.Type {
	case "instances":
		err := input.listInstances()
		if err != nil {
			return err
		}
	case "sessions":
		err := input.listSessions()
		if err != nil {
			return err
		}
	default:
		err := fmt.Errorf("Unsupported list type: %s", *input.Type)
		log.WithField("type", *input.Type).Error(err)
		return err
	}
	return nil
}

// String will stringify StartOutput
func (o *StartOutput) String() (string, error) {
	switch *o.format {
	case "text":
		output := ""
		buf := bytes.NewBufferString(output)
		w := new(tabwriter.Writer)
		w.Init(buf, 0, 0, 2, ' ', 0)
		switch *o.outputType {
		case "instances":
			fmt.Fprintln(w, "Index\tName\tInstance ID\tIP Address\tPrivate DNS Name")
			for i, instance := range o.Instances {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\n",
					(i + 1), *instance.Name, *instance.InstanceID, *instance.IPAddress, *instance.PrivateDNSName)
			}
			err := w.Flush()
			if err != nil {
				return "", err
			}
		case "sessions":
			fmt.Fprintln(w, "Index\tSession ID\tTarget\tStart Date")
			for i, session := range o.Sessions {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\n",
					(i + 1), *session.SessionID, *session.Target, *session.StartDate)
			}
			err := w.Flush()
			if err != nil {
				return "", err
			}
		default:
			err := fmt.Errorf("Unsupported type: %s", *o.outputType)
			log.WithField("outputType", *o.outputType).Error(err)
			return "", err
		}
		return buf.String(), nil
	case "json":
		var data []byte
		var err error
		switch *o.outputType {
		case "instances":
			data, err = json.Marshal(o.Instances)
			if err != nil {
				return "", err
			}
		case "sessions":
			data, err = json.Marshal(o.Sessions)
			if err != nil {
				return "", err
			}
		default:
			err = fmt.Errorf("Unsupported type: %s", *o.outputType)
			log.WithField("outputType", *o.outputType).Error(err)
			return "", err
		}
		// JSON output was missing new line
		return string(data) + "\n", nil
	case "yaml":
		var data []byte
		var err error
		switch *o.outputType {
		case "instances":
			data, err = yaml.Marshal(o.Instances)
			if err != nil {
				return "", err
			}
		case "sessions":
			data, err = yaml.Marshal(o.Sessions)
			if err != nil {
				return "", err
			}
		default:
			err = fmt.Errorf("Unsupported type: %s", *o.outputType)
			log.WithField("outputType", *o.outputType).Error(err)
			return "", err
		}
		return string(data), nil
	case "wide":
		buf := new(bytes.Buffer)
		w := new(tabwriter.Writer)
		w.Init(buf, 0, 0, 2, ' ', 0)
		switch *o.outputType {
		case "instances":
			fmt.Fprintln(w, "Index\tName\tInstance ID\tIP Address\tPrivate DNS Name\tHostname\tOS Name\tOS Version\tOS Type")
			for i, instance := range o.Instances {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
					(i + 1), *instance.Name, *instance.InstanceID, *instance.IPAddress, *instance.PrivateDNSName,
					*instance.Hostname, *instance.OSName, *instance.OSVersion, *instance.OSType)
			}
			err := w.Flush()
			if err != nil {
				return "", err
			}
		case "sessions":
			fmt.Fprintln(w, "Index\tSession ID\tTarget\tStart Date\tOwner\tStatus")
			for i, session := range o.Sessions {
				fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
					(i + 1), *session.SessionID, *session.Target, *session.StartDate,
					*session.Owner, *session.Status)
			}
			err := w.Flush()
			if err != nil {
				return "", err
			}
		default:
			err := fmt.Errorf("Unsupported type: %s", *o.outputType)
			log.WithField("outputType", *o.outputType).Error(err)
			return "", err
		}
		return buf.String(), nil
	default:
		return "", fmt.Errorf("Unsupported output format: %s", *o.format)
	}
}

func (input *StartInput) listInstances() error {
	instanceList := make([]*Instance, 0, capacity)
	output := &StartOutput{
		format:     input.OutputFormat,
		outputType: input.Type,
		Instances:  instanceList,
	}
	// Get the list of instances
	ssmDescribeInstancesInput := &ssm.DescribeInstanceInformationInput{}
	if *input.Filters != "" {
		filters := Filters{}
		log.WithField("filetrs", *input.Filters).Debug("Filters")
		err := json.Unmarshal([]byte(*input.Filters), &filters)
		if err != nil {
			return err
		}
		filterList := []*ssm.InstanceInformationStringFilter{}
		for _, tag := range filters.Tags {
			log.WithFields(log.Fields{
				"key":    tag.Key,
				"values": tag.Values,
			}).Debug("Tags Filter")
			filterList = append(filterList, &ssm.InstanceInformationStringFilter{
				Key:    aws.String("tag:" + tag.Key),
				Values: tag.Values,
			})
		}
		if len(filters.InstanceIDs) > 0 {
			log.WithFields(log.Fields{
				"IDs": filters.InstanceIDs,
			}).Debug("Instance IDs Filter")
			key := "InstanceIds"
			filterList = append(filterList, &ssm.InstanceInformationStringFilter{
				Key:    &key,
				Values: filters.InstanceIDs,
			})
		}
		log.WithFields(log.Fields{
			"filterList": filterList,
		}).Debug("Describe Instance Filters")
		ssmDescribeInstancesInput.Filters = filterList
	}
	ssmClient := ssm.New(input.AWSSession)
	err := ssmClient.DescribeInstanceInformationPages(ssmDescribeInstancesInput,
		func(page *ssm.DescribeInstanceInformationOutput, lastPage bool) bool {
			for _, instance := range page.InstanceInformationList {
				if len(output.Instances)+1 > cap(output.Instances) {
					newSlice := make([]*Instance, len(output.Instances), (cap(output.Instances))*capMultiplier)
					n := copy(newSlice, output.Instances)
					log.WithField("no. copied elements", n).Debug("Expand Instances slice")
					output.Instances = newSlice
				}
				log.WithFields(log.Fields{
					"InstanceId":      *instance.InstanceId,
					"ComputerName":    *instance.ComputerName,
					"IPAddress":       *instance.IPAddress,
					"PlatformName":    *instance.PlatformName,
					"PlatformType":    *instance.PlatformType,
					"PlatformVersion": *instance.PlatformVersion,
				}).Debug("Describe Instance")
				output.Instances = append(
					output.Instances,
					&Instance{
						Hostname:   instance.ComputerName,
						IPAddress:  instance.IPAddress,
						InstanceID: instance.InstanceId,
						OSName:     instance.PlatformName,
						OSType:     instance.PlatformType,
						OSVersion:  instance.PlatformVersion,
					},
				)
			}
			return !lastPage
		})
	if err != nil {
		return err
	}
	if len(output.Instances) < 1 {
		log.Debug("No instances found")
		outString, err := output.String()
		if err != nil {
			return err
		}
		fmt.Print(outString)
		return nil
	}
	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: make([]*string, 0, cap(output.Instances)),
	}
	// Adding instances private DNS name
	for _, instance := range output.Instances {
		describeInstancesInput.InstanceIds = append(describeInstancesInput.InstanceIds, instance.InstanceID)
	}
	// 0 for PrivateDNSName, 1 for Name Tag
	describeInstance := make(map[string][2]*string)
	ec2Client := ec2.New(input.AWSSession)
	err = ec2Client.DescribeInstancesPages(describeInstancesInput,
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, reservation := range page.Reservations {
				for _, instance := range reservation.Instances {
					nameTag := ""
					for _, tag := range instance.Tags {
						if *tag.Key == "Name" {
							nameTag = *tag.Value
							break
						}
					}
					describeInstance[*instance.InstanceId] = [2]*string{
						instance.PrivateDnsName, &nameTag}
				}
			}
			return !lastPage
		})
	if err != nil {
		return err
	}
	for _, instance := range output.Instances {
		instance.PrivateDNSName = describeInstance[*instance.InstanceID][0]
		instance.Name = describeInstance[*instance.InstanceID][1]
		log.WithFields(log.Fields{
			"IPAddress":      *instance.IPAddress,
			"Hostname":       *instance.Hostname,
			"OSName":         *instance.OSName,
			"OSType":         *instance.OSType,
			"OSVersion":      *instance.OSVersion,
			"PrivateDNSName": *instance.PrivateDNSName,
			"Name":           *instance.Name,
		}).Debug("Instance")
	}
	outString, err := output.String()
	if err != nil {
		return err
	}
	fmt.Print(outString)
	if *input.Interactive && len(output.Instances) > 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Fprintf(os.Stderr, "Choose an instance to connect to [1 - %d]: ", len(output.Instances))
		textInput, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		i, err := strconv.Atoi(strings.ReplaceAll(textInput, "\n", ""))
		if err != nil {
			return err
		}
		log.WithField("index", i).Debug("Picked EC2 Instance")
		if i < 1 || i > len(output.Instances) {
			return fmt.Errorf("instance index out of range: %d", i)
		}
		instance := output.Instances[i-1]
		targetType := "instance-id"
		remoteInput := &remoteSession.StartInput{
			Target:     instance.InstanceID,
			TargetType: &targetType,
			AWSSession: input.AWSSession,
		}
		err = remoteSession.Start(remoteInput)
		if err != nil {
			return err
		}
	}
	return nil
}

func (input *StartInput) listSessions() error {
	sessionList := make([]*Session, 0, capacity)
	output := &StartOutput{
		format:     input.OutputFormat,
		outputType: input.Type,
		Sessions:   sessionList,
	}
	// By default show only connected sessions
	ssmDescribeSessionsInput := &ssm.DescribeSessionsInput{
		State: aws.String("Active"),
	}
	// Parse filters
	if *input.Filters != "" {
		data := Filters{}
		filters := []*ssm.SessionFilter{}
		err := json.Unmarshal([]byte(*input.Filters), &data)
		if err != nil {
			return err
		}
		if data.After != nil {
			filters = append(filters, &ssm.SessionFilter{
				Key:   aws.String("InvokedAfter"),
				Value: data.After,
			})
		}
		if data.Before != nil {
			filters = append(filters, &ssm.SessionFilter{
				Key:   aws.String("InvokedBefore"),
				Value: data.Before,
			})
		}
		if data.Target != nil {
			filters = append(filters, &ssm.SessionFilter{
				Key:   aws.String("Target"),
				Value: data.Target,
			})
		}
		if data.Owner != nil {
			filters = append(filters, &ssm.SessionFilter{
				Key:   aws.String("Owner"),
				Value: data.Owner,
			})
		}
		ssmDescribeSessionsInput.Filters = filters
	}
	ssmClient := ssm.New(input.AWSSession)
	for out, err := ssmClient.DescribeSessions(ssmDescribeSessionsInput); ; {
		if err != nil {
			return err
		}
		log.WithField("sessions array len", len(out.Sessions)).Debug("Sessions Output")
		for i, sess := range out.Sessions {
			log.WithField("session", sess).Debugf("Single session #%d", i)
			startDate, err := sess.StartDate.MarshalText()
			if err != nil {
				return err
			}
			startDateString := string(startDate[:])
			output.Sessions = append(output.Sessions, &Session{
				SessionID: sess.SessionId,
				Target:    sess.Target,
				Status:    sess.Status,
				StartDate: &startDateString,
				Owner:     sess.Owner,
			})
		}
		if out.NextToken == nil {
			break
		}
		ssmDescribeSessionsInput.NextToken = out.NextToken
	}
	outString, err := output.String()
	if err != nil {
		return err
	}
	fmt.Print(outString)
	if *input.Interactive && len(output.Sessions) > 0 {
		reader := bufio.NewReader(os.Stdin)
		fmt.Fprintf(os.Stderr, "Terminate session [1 - %d]: ", len(output.Sessions))
		textInput, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		i, err := strconv.Atoi(strings.ReplaceAll(textInput, "\n", ""))
		if err != nil {
			return err
		}
		log.WithField("index", i).Debug("Picked session")
		if i < 1 || i > len(output.Sessions) {
			return fmt.Errorf("session index out of range: %d", i)
		}
		chosenSession := output.Sessions[i-1]
		err = remoteSession.TerminateSession(ssmClient, chosenSession.SessionID)
		if err != nil {
			return err
		}
	}
	return nil
}
