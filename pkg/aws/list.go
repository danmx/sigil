package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	log "github.com/sirupsen/logrus"
)

const capMultiplier = 2

// ListInstances provides a list of active instances with SSM agent
func (p *Provider) ListInstances() ([]*Instance, error) {
	instances, err := fetchInstances(ssm.New(p.session), p.filters.Instance.IDs)
	if err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		log.Info("No matching instances")
		return []*Instance{}, nil
	}
	filteredInstances, err := filterInstances(ec2.New(p.session), p.filters.Instance.Tags, instances)
	if err != nil {
		return nil, err
	}
	return filteredInstances, nil
}

func fetchInstances(ssmClient ssmiface.SSMAPI, ids []string) (map[string]*Instance, error) {
	// Show only instances that have active agents
	ssmDescribeInstancesInput := &ssm.DescribeInstanceInformationInput{
		MaxResults: aws.Int64(maxResults),
		Filters: []*ssm.InstanceInformationStringFilter{
			{
				Key:    aws.String("PingStatus"),
				Values: []*string{aws.String("Online")},
			},
		},
	}
	if len(ids) > 0 {
		log.WithFields(log.Fields{
			"ids": ids,
		}).Debug("Instance IDs Filter")
		ssmDescribeInstancesInput.Filters = append(ssmDescribeInstancesInput.Filters, &ssm.InstanceInformationStringFilter{
			Key:    aws.String("InstanceIds"),
			Values: aws.StringSlice(ids),
		})
	}

	instances := make(map[string]*Instance)

	err := ssmClient.DescribeInstanceInformationPages(ssmDescribeInstancesInput,
		func(page *ssm.DescribeInstanceInformationOutput, lastPage bool) bool {
			for _, instance := range page.InstanceInformationList {
				log.WithFields(log.Fields{
					"InstanceId":      *instance.InstanceId,
					"ComputerName":    *instance.ComputerName,
					"IPAddress":       *instance.IPAddress,
					"PlatformName":    *instance.PlatformName,
					"PlatformType":    *instance.PlatformType,
					"PlatformVersion": *instance.PlatformVersion,
				}).Debug("Describe Instance")
				instances[*instance.InstanceId] = &Instance{
					Hostname:  *instance.ComputerName,
					IPAddress: *instance.IPAddress,
					ID:        *instance.InstanceId,
					OSName:    *instance.PlatformName,
					OSType:    *instance.PlatformType,
					OSVersion: *instance.PlatformVersion,
				}
			}
			return !lastPage
		})
	if err != nil {
		log.WithFields(log.Fields{
			"input": *ssmDescribeInstancesInput,
			"error": err,
		}).Error("DescribeInstanceInformationPages")
		return nil, err
	}
	return instances, nil
}

func filterInstances(ec2Client ec2iface.EC2API, tags []TagValues, instances map[string]*Instance) ([]*Instance, error) {
	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: make([]*string, 0, len(instances)),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}

	for _, tag := range tags {
		log.WithFields(log.Fields{
			"key":    tag.Key,
			"values": tag.Values,
		}).Debug("Tags Filter")
		describeInstancesInput.Filters = append(describeInstancesInput.Filters, &ec2.Filter{
			Name:   aws.String("tag:" + tag.Key),
			Values: aws.StringSlice(tag.Values),
		})
	}

	outputInstances := make([]*Instance, 0, len(instances))
	// Discovering instances private DNS names
	for _, instance := range instances {
		describeInstancesInput.InstanceIds = append(describeInstancesInput.InstanceIds, &instance.ID)
	}
	err := ec2Client.DescribeInstancesPages(describeInstancesInput,
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
					instances[*instance.InstanceId].PrivateDNSName = *instance.PrivateDnsName
					instances[*instance.InstanceId].Name = nameTag
					outputInstances = append(outputInstances, instances[*instance.InstanceId])
				}
			}
			return !lastPage
		})
	if err != nil {
		log.WithFields(log.Fields{
			"input": *describeInstancesInput,
			"error": err,
		}).Error("DescribeInstancesPages")
		return nil, err
	}
	return outputInstances, nil
}

// ListSessions provides a list of active SSM sessions
func (p *Provider) ListSessions() ([]*Session, error) {
	return listSessions(ssm.New(p.session), p.filters.Session)
}

func listSessions(ssmClient ssmiface.SSMAPI, sessionFilters SessionFilters) ([]*Session, error) {
	// Show only connected sessions
	ssmDescribeSessionsInput := &ssm.DescribeSessionsInput{
		State: aws.String("Active"),
	}
	// Parse filters
	filters := []*ssm.SessionFilter{}
	if sessionFilters.After != "" {
		filters = append(filters, &ssm.SessionFilter{
			Key:   aws.String("InvokedAfter"),
			Value: &sessionFilters.After,
		})
	}
	if sessionFilters.Before != "" {
		filters = append(filters, &ssm.SessionFilter{
			Key:   aws.String("InvokedBefore"),
			Value: &sessionFilters.Before,
		})
	}
	if sessionFilters.Target != "" {
		filters = append(filters, &ssm.SessionFilter{
			Key:   aws.String("Target"),
			Value: &sessionFilters.Target,
		})
	}
	if sessionFilters.Owner != "" {
		filters = append(filters, &ssm.SessionFilter{
			Key:   aws.String("Owner"),
			Value: &sessionFilters.Owner,
		})
	}
	if len(filters) > 0 {
		ssmDescribeSessionsInput.Filters = filters
	}
	sessions := []*Session{}
	for out, err := ssmClient.DescribeSessions(ssmDescribeSessionsInput); ; {
		if err != nil {
			log.WithField("error", err).Error("DescribeSessions")
			return nil, err
		}
		log.WithField("sessions array len", len(out.Sessions)).Debug("Sessions Output")
		if len(sessions)+1 > cap(sessions) {
			newSlice := make([]*Session, len(sessions), (cap(sessions))*capMultiplier)
			n := copy(newSlice, sessions)
			log.WithField("no. copied elements", n).Debug("Expand Sessions slice")
			sessions = newSlice
		}
		for i, sess := range out.Sessions {
			log.WithField("session", sess).Debugf("Single session #%d", i)
			startDate, err := sess.StartDate.MarshalText()
			if err != nil {
				log.WithField("error", err).Error("StartDate MarshalText")
				return nil, err
			}
			startDateString := string(startDate)
			sessions = append(sessions, &Session{
				SessionID: *sess.SessionId,
				Target:    *sess.Target,
				Status:    *sess.Status,
				StartDate: startDateString,
				Owner:     *sess.Owner,
			})
		}
		if out.NextToken == nil {
			break
		}
		ssmDescribeSessionsInput.NextToken = out.NextToken
	}
	return sessions, nil
}
