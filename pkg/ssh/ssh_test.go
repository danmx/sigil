package ssh

// AWS Mock generation
//go:generate go run github.com/golang/mock/mockgen -self_package=github.com/danmx/sigil/pkg/ssh -package ssh -destination aws_mock_test.go github.com/danmx/sigil/pkg/aws Cloud,CloudInstances,CloudSessions,CloudSSH

import (
	"os"
	"path"
	"testing"

	"github.com/danmx/sigil/pkg/aws"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestStart verifies start ssh method
func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockCloudSSH(ctrl) // skipcq: SCC-compile

	target := "i-xxxxxxxxxxxxxxxx1"
	targetType := aws.TargetTypeInstanceID
	mfa := "123456"
	region := "eu-west-1"
	profile := "west"
	var port uint64 = 22
	pubKey := path.Join(os.TempDir(), "sigil_test.pub")
	osUser := "ec2-user"
	genKey := true
	input := StartInput{
		MFAToken:   &mfa,
		Region:     &region,
		Profile:    &profile,
		Target:     &target,
		TargetType: &targetType,
		PortNumber: &port,
		PublicKey:  &pubKey,
		OSUser:     &osUser,
		GenKeyPair: &genKey,
	}

	gomock.InOrder(
		m.EXPECT().NewWithConfig(gomock.Eq(&aws.Config{
			Region:   *input.Region,
			Profile:  *input.Profile,
			MFAToken: *input.MFAToken,
		})).Return(nil),
		m.EXPECT().StartSSH(
			gomock.Eq(*input.TargetType),
			gomock.Eq(*input.Target),
			gomock.Eq(*input.OSUser),
			gomock.Eq(*input.PortNumber),
			gomock.Any(),
		).Return(nil),
	)

	assert.NoError(t, input.start(m))
}
