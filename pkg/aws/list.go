package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	log "github.com/sirupsen/logrus"
)

const capMultiplier = 2

// ListInstances provides a list of active instances with SSM agent
func (p *Provider) ListInstances() ([]*Instance, error) {
	instances := make(map[string]*Instance)

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
	for _, tag := range p.filters.Instance.Tags {
		log.WithFields(log.Fields{
			"key":    tag.Key,
			"values": tag.Values,
		}).Debug("Tags Filter")
		ssmDescribeInstancesInput.Filters = append(ssmDescribeInstancesInput.Filters, &ssm.InstanceInformationStringFilter{
			Key:    aws.String("tag:" + tag.Key),
			Values: aws.StringSlice(tag.Values),
		})
	}
	if len(p.filters.Instance.IDs) > 0 {
		log.WithFields(log.Fields{
			"IDs": p.filters.Instance.IDs,
		}).Debug("Instance IDs Filter")
		ssmDescribeInstancesInput.Filters = append(ssmDescribeInstancesInput.Filters, &ssm.InstanceInformationStringFilter{
			Key:    aws.String("InstanceIds"),
			Values: aws.StringSlice(p.filters.Instance.IDs),
		})
	}

	ssmClient := ssm.New(p.session)
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
	if len(instances) == 0 {
		log.Info("No matching instances")
		return []*Instance{}, nil
	}
	describeInstancesInput := &ec2.DescribeInstancesInput{
		InstanceIds: make([]*string, 0, len(instances)),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}
	outputInstances := make([]*Instance, 0, len(instances))
	// Discovering instances private DNS names
	for _, instance := range instances {
		describeInstancesInput.InstanceIds = append(describeInstancesInput.InstanceIds, &instance.ID)
		outputInstances = append(outputInstances, instance)
	}
	ec2Client := ec2.New(p.session)
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
					instances[*instance.InstanceId].PrivateDNSName = *instance.PrivateDnsName
					instances[*instance.InstanceId].Name = nameTag
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
	// Show only connected sessions
	ssmDescribeSessionsInput := &ssm.DescribeSessionsInput{
		State: aws.String("Active"),
	}
	// Parse filters
	filters := []*ssm.SessionFilter{}
	if p.filters.Session.After != "" {
		filters = append(filters, &ssm.SessionFilter{
			Key:   aws.String("InvokedAfter"),
			Value: &p.filters.Session.After,
		})
	}
	if p.filters.Session.Before != "" {
		filters = append(filters, &ssm.SessionFilter{
			Key:   aws.String("InvokedBefore"),
			Value: &p.filters.Session.Before,
		})
	}
	if p.filters.Session.Target != "" {
		filters = append(filters, &ssm.SessionFilter{
			Key:   aws.String("Target"),
			Value: &p.filters.Session.Target,
		})
	}
	if p.filters.Session.Owner != "" {
		filters = append(filters, &ssm.SessionFilter{
			Key:   aws.String("Owner"),
			Value: &p.filters.Session.Owner,
		})
	}
	ssmDescribeSessionsInput.Filters = filters
	ssmClient := ssm.New(p.session)
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
