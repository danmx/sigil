package aws

// Helpers Mock generation
//go:generate go run github.com/golang/mock/mockgen -self_package=github.com/danmx/sigil/pkg/aws -package aws -destination helpers_mock_test.go github.com/danmx/sigil/pkg/aws/helpers OSExecIface,OSIface
// AWS EC2 Mock generation
//go:generate go run github.com/golang/mock/mockgen -self_package=github.com/danmx/sigil/pkg/aws -package aws -destination ec2_aws_mock_test.go github.com/aws/aws-sdk-go/service/ec2/ec2iface EC2API

import (
	"errors"
	"testing"

	awsCloud "github.com/aws/aws-sdk-go/aws"
	_ "github.com/aws/aws-sdk-go/aws/request" // for mocking EC2API
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestVerifyDependencies verifies the dependency verifier
func TestVerifyDependencies(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockOSExecIface(ctrl) // skipcq: SCC-compile

	gomock.InOrder(
		m.EXPECT().LookPath(pluginName).Return("", errors.New("")),
	)

	assert.Error(t, verifyDependencies(m))

	gomock.InOrder(
		m.EXPECT().LookPath(pluginName).Return("/usr/local/bin/session-manager-plugin", nil),
	)

	assert.NoError(t, verifyDependencies(m))
}

// TestAppendUserAgent verifies env var was set
func TestAppendUserAgent(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockOSIface(ctrl) // skipcq: SCC-compile

	suffix := "test"

	gomock.InOrder(
		m.EXPECT().LookupEnv(execEnvVar).Return("", false),
		m.EXPECT().Setenv(execEnvVar, suffix).Return(nil),
	)

	assert.NoError(t, appendUserAgent(m, suffix))

	gomock.InOrder(
		m.EXPECT().LookupEnv(execEnvVar).Return("value", true),
		m.EXPECT().Setenv(execEnvVar, "value/"+suffix).Return(nil),
	)

	assert.NoError(t, appendUserAgent(m, suffix))
}

// TestGetFilters tests fetching instance filters
func TestGetFilters(t *testing.T) {
	target := "testTarget"
	testMap := map[string][]*ec2.Filter{
		TargetTypeInstanceID: {
			{
				Name:   awsCloud.String("instance-id"),
				Values: []*string{&target},
			},
			{
				Name:   awsCloud.String("instance-state-name"),
				Values: []*string{awsCloud.String("running")},
			},
		},
		TargetTypePrivateDNS: {
			{
				Name:   awsCloud.String("private-dns-name"),
				Values: []*string{&target},
			},
			{
				Name:   awsCloud.String("instance-state-name"),
				Values: []*string{awsCloud.String("running")},
			},
		},
		TargetTypeName: {
			{
				Name:   awsCloud.String("tag:Name"),
				Values: []*string{&target},
			},
			{
				Name:   awsCloud.String("instance-state-name"),
				Values: []*string{awsCloud.String("running")},
			},
		},
	}
	for key, value := range testMap {
		filter, err := getFilters(key, target)
		assert.Equal(t, filter, value)
		assert.NoError(t, err)
	}
	_, err := getFilters("invalid", target)
	assert.Error(t, err)
}

// TestGetFirstInstance tests listing the first instance
func TestGetFirstInstance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockEC2API(ctrl) // skipcq: SCC-compile

	filters := []*ec2.Filter{}

	input := &ec2.DescribeInstancesInput{
		Filters:    filters,
		MaxResults: awsCloud.Int64(maxResults),
	}

	gomock.InOrder(
		m.EXPECT().DescribeInstancesPages(input, gomock.Any()).Return(nil),
		m.EXPECT().DescribeInstancesPages(input, gomock.Any()).Return(errors.New("")),
	)
	_, err := getFirstInstance(m, filters)
	assert.NoError(t, err)
	_, err = getFirstInstance(m, filters)
	assert.Error(t, err)
}
