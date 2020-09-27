package aws

// AWS EC2 Instance Connect Mock generation
//go:generate go run github.com/golang/mock/mockgen -self_package=github.com/danmx/sigil/pkg/aws -package aws -destination ec2instanceconnect_aws_mock_test.go github.com/aws/aws-sdk-go/service/ec2instanceconnect/ec2instanceconnectiface EC2InstanceConnectAPI

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2instanceconnect"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestUploadPublicKey verifies upload of a public key
func TestUploadPublicKey(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockEC2InstanceConnectAPI(ctrl) // skipcq: SCC-compile

	availabilityZone := "a"
	instanceID := "1"
	osUser := "ec2user"
	pubKey := "testKey"

	output := &ec2instanceconnect.SendSSHPublicKeyOutput{
		RequestId: aws.String("testSuccess"),
		Success:   aws.Bool(true),
	}

	gomock.InOrder(
		m.EXPECT().SendSSHPublicKey(&ec2instanceconnect.SendSSHPublicKeyInput{
			AvailabilityZone: &availabilityZone,
			InstanceId:       &instanceID,
			InstanceOSUser:   &osUser,
			SSHPublicKey:     &pubKey,
		}).Return(output, nil),
	)

	assert.NoError(t, uploadPublicKey(m, []byte(pubKey), osUser, instanceID, availabilityZone))

	output = &ec2instanceconnect.SendSSHPublicKeyOutput{
		Success:   aws.Bool(false),
		RequestId: aws.String("testFailure"),
	}

	gomock.InOrder(
		m.EXPECT().SendSSHPublicKey(&ec2instanceconnect.SendSSHPublicKeyInput{
			AvailabilityZone: &availabilityZone,
			InstanceId:       &instanceID,
			InstanceOSUser:   &osUser,
			SSHPublicKey:     &pubKey,
		}).Return(output, nil),
	)

	assert.Error(t, uploadPublicKey(m, []byte(pubKey), osUser, instanceID, availabilityZone))
}
