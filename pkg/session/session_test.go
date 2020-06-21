package session

// AWS Mock generation
//go:generate go run github.com/golang/mock/mockgen -self_package=github.com/danmx/sigil/pkg/session -package session -destination aws_mock_test.go github.com/danmx/sigil/pkg/aws Cloud,CloudInstances,CloudSessions,CloudSSH

import (
	"testing"

	"github.com/danmx/sigil/pkg/aws"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestStart verifies start session method for different configurations
func TestStart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockCloudInstances(ctrl) // skipcq: SCC-compile

	mfa := "123456"
	region := "eu-west-1"
	profile := "west"
	input := StartInput{
		MFAToken: &mfa,
		Region:   &region,
		Profile:  &profile,
	}
	// Instance ID
	target := "i-xxxxxxxxxxxxxxxx1"
	targetType := aws.TargetTypeInstanceID
	input.Target = &target
	input.TargetType = &targetType
	gomock.InOrder(
		m.EXPECT().NewWithConfig(gomock.Eq(&aws.Config{
			Region:   *input.Region,
			Profile:  *input.Profile,
			MFAToken: *input.MFAToken,
		})).Return(nil),
		m.EXPECT().StartSession(gomock.Eq(*input.TargetType), gomock.Eq(*input.Target)).Return(nil),
	)
	assert.NoError(t, input.start(m))
	// DNS
	target = "test.local"
	targetType = aws.TargetTypePrivateDNS
	gomock.InOrder(
		m.EXPECT().NewWithConfig(gomock.Eq(&aws.Config{
			Region:   *input.Region,
			Profile:  *input.Profile,
			MFAToken: *input.MFAToken,
		})).Return(nil),
		m.EXPECT().StartSession(gomock.Eq(*input.TargetType), gomock.Eq(*input.Target)).Return(nil),
	)
	assert.NoError(t, input.start(m))
	// Name
	target = "Backend"
	targetType = aws.TargetTypeName
	gomock.InOrder(
		m.EXPECT().NewWithConfig(gomock.Eq(&aws.Config{
			Region:   *input.Region,
			Profile:  *input.Profile,
			MFAToken: *input.MFAToken,
		})).Return(nil),
		m.EXPECT().StartSession(gomock.Eq(*input.TargetType), gomock.Eq(*input.Target)).Return(nil),
	)
	assert.NoError(t, input.start(m))
	// Deprecated Name
	targetType = aws.DeprecatedTargetTypeName
	gomock.InOrder(
		m.EXPECT().NewWithConfig(gomock.Eq(&aws.Config{
			Region:   *input.Region,
			Profile:  *input.Profile,
			MFAToken: *input.MFAToken,
		})).Return(nil),
		m.EXPECT().StartSession(gomock.Eq(*input.TargetType), gomock.Eq(*input.Target)).Return(nil),
	)
	assert.NoError(t, input.start(m))
}
