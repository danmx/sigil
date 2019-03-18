package list

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/tabwriter"

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
	AWSRegion    *string
	TagFilter    *map[string]string
}

// StartOutput struct will contain all output data
type StartOutput struct {
	Instances *[]*Instance
	// Define output format
	format *string
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

const (
	instanceCapacity      = 10
	instanceCapMultiplier = 2
)

// Start will output a ist of all available EC2 instances
func Start(input *StartInput) (*StartOutput, error) {
	instanceList := make([]*Instance, 0, instanceCapacity)
	output := &StartOutput{
		format:    input.OutputFormat,
		Instances: &instanceList,
	}
	awsConfig := aws.NewConfig()
	if *input.AWSRegion != "" {
		awsConfig.Region = input.AWSRegion
	}
	sess := session.Must(session.NewSession(awsConfig))
	// Get the list of instances
	ssmDescribeInstancesInput := &ssm.DescribeInstanceInformationInput{}
	if len(*input.TagFilter) > 0 {
		filterList := []*ssm.InstanceInformationStringFilter{}
		for key, value := range *input.TagFilter {
			filterList = append(filterList, &ssm.InstanceInformationStringFilter{
				Key:    aws.String("tag:" + key),
				Values: []*string{&value},
			})
		}
		ssmDescribeInstancesInput.Filters = filterList
	}
	ssmClient := ssm.New(sess)
	err := ssmClient.DescribeInstanceInformationPages(ssmDescribeInstancesInput,
		func(page *ssm.DescribeInstanceInformationOutput, lastPage bool) bool {
			for _, instance := range page.InstanceInformationList {
				if len(*output.Instances)+1 > cap(*output.Instances) {
					newSlice := make([]*Instance, len(*output.Instances), (cap(*output.Instances))*instanceCapMultiplier)
					n := copy(newSlice, *output.Instances)
					log.WithField("no. copied elements", n).Debug("Expand Instances slice")
					*output.Instances = newSlice
				}
				log.WithFields(log.Fields{
					"InstanceId":      *instance.InstanceId,
					"ComputerName":    *instance.ComputerName,
					"IPAddress":       *instance.IPAddress,
					"PlatformName":    *instance.PlatformName,
					"PlatformType":    *instance.PlatformType,
					"PlatformVersion": *instance.PlatformVersion,
				}).Debug("Describe Instance")
				*output.Instances = append(
					*output.Instances,
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
		return nil, err
	}
	if len(*output.Instances) < 1 {
		log.Debug("No instances found")
		return output, nil
	}
	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: make([]*string, 0, cap(*output.Instances)),
	}
	// Adding instances private DNS name
	for _, instance := range *output.Instances {
		describeInstancesInput.InstanceIds = append(describeInstancesInput.InstanceIds, instance.InstanceID)
	}
	// 0 for PrivateDNSName, 1 for Name Tag
	describeInstance := make(map[string][2]*string)
	ec2Client := ec2.New(sess)
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
		return nil, err
	}
	for _, instance := range *output.Instances {
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
	return output, nil
}

// String will stringify StartOutput
func (o *StartOutput) String() (string, error) {
	switch *o.format {
	case "text":
		output := ""
		buf := bytes.NewBufferString(output)
		w := new(tabwriter.Writer)
		w.Init(buf, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Name\tInstance ID\tIP Address\tPrivate DNS Name")
		for _, instance := range *o.Instances {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n",
				*instance.Name, *instance.InstanceID, *instance.IPAddress, *instance.PrivateDNSName)
		}
		err := w.Flush()
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	case "json":
		data, err := json.Marshal(*o.Instances)
		if err != nil {
			return "", err
		}
		// JSON output was missing new line
		return string(data) + "\n", nil
	case "yaml":
		data, err := yaml.Marshal(*o.Instances)
		if err != nil {
			return "", err
		}
		return string(data), nil
	case "wide":
		buf := new(bytes.Buffer)
		w := new(tabwriter.Writer)
		w.Init(buf, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "Name\tInstance ID\tIP Address\tPrivate DNS Name\tHostname\tOS Name\tOS Version\tOS Type")
		for _, instance := range *o.Instances {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
				*instance.Name, *instance.InstanceID, *instance.IPAddress, *instance.PrivateDNSName,
				*instance.Hostname, *instance.OSName, *instance.OSVersion, *instance.OSType)
		}
		err := w.Flush()
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	default:
		return "", fmt.Errorf("Unsupported output format: %s", *o.format)
	}
}
