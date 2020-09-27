package aws

// AWS SSM Mock generation
//go:generate go run github.com/golang/mock/mockgen -self_package=github.com/danmx/sigil/pkg/aws -package aws -destination ssm_aws_mock_test.go github.com/aws/aws-sdk-go/service/ssm/ssmiface SSMAPI
// AWS EC2 Mock generation
//go:generate go run github.com/golang/mock/mockgen -self_package=github.com/danmx/sigil/pkg/aws -package aws -destination ec2_aws_mock_test.go github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestFetchInstances verifies instances fetching
func TestFetchInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockSSMAPI(ctrl) // skipcq: SCC-compile

	IDs := []string{"1", "2", "3"}

	instances := make(map[string]*Instance)

	input := &ssm.DescribeInstanceInformationInput{
		MaxResults: aws.Int64(maxResults),
		Filters: []*ssm.InstanceInformationStringFilter{
			{
				Key:    aws.String("PingStatus"),
				Values: []*string{aws.String("Online")},
			},
			{
				Key:    aws.String("InstanceIds"),
				Values: aws.StringSlice(IDs),
			},
		},
	}

	gomock.InOrder(
		m.EXPECT().DescribeInstanceInformationPages(input, gomock.Any()).Return(nil),
	)
	list, err := fetchInstances(m, IDs)
	assert.Equal(t, instances, list)
	assert.NoError(t, err)
}

// TestFilterInstances verifies instances filtering
func TestFilterInstances(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockEC2API(ctrl) // skipcq: SCC-compile

	tags := []TagValues{
		{
			Key:    "TestKey1",
			Values: []string{"1"},
		},
	}

	instances := make(map[string]*Instance)

	outputInstances := []*Instance{}

	input := &ec2.DescribeInstancesInput{
		InstanceIds: make([]*string, 0, len(instances)),
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("instance-state-name"),
				Values: []*string{aws.String("running")},
			},
		},
	}
	for _, tag := range tags {
		input.Filters = append(input.Filters, &ec2.Filter{
			Name:   aws.String("tag:" + tag.Key),
			Values: aws.StringSlice(tag.Values),
		})
	}
	for _, instance := range instances {
		input.InstanceIds = append(input.InstanceIds, &instance.ID)
	}

	gomock.InOrder(
		m.EXPECT().DescribeInstancesPages(input, gomock.Any()).Return(nil),
	)
	list, err := filterInstances(m, tags, instances)
	assert.Equal(t, outputInstances, list)
	assert.NoError(t, err)
}

// TestListSessions verifies session listing
func TestListSessions(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockSSMAPI(ctrl) // skipcq: SCC-compile

	filters := SessionFilters{
		After:  "after",
		Before: "before",
		Target: "target",
		Owner:  "owner",
	}

	sessions := []*Session{}

	input := &ssm.DescribeSessionsInput{
		State: aws.String("Active"),
		Filters: []*ssm.SessionFilter{
			{
				Key:   aws.String("InvokedAfter"),
				Value: &filters.After,
			},
			{
				Key:   aws.String("InvokedBefore"),
				Value: &filters.Before,
			},
			{
				Key:   aws.String("Target"),
				Value: &filters.Target,
			},
			{
				Key:   aws.String("Owner"),
				Value: &filters.Owner,
			},
		},
	}

	output := &ssm.DescribeSessionsOutput{
		NextToken: nil,
		Sessions:  []*ssm.Session{},
	}

	gomock.InOrder(
		m.EXPECT().DescribeSessions(input).Return(output, nil),
	)

	list, err := listSessions(m, filters)
	assert.Equal(t, sessions, list)
	assert.NoError(t, err)
}
