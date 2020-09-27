package aws

// AWS SSM Mock generation
//go:generate go run github.com/golang/mock/mockgen -self_package=github.com/danmx/sigil/pkg/aws -package aws -destination ssm_aws_mock_test.go github.com/aws/aws-sdk-go/service/ssm/ssmiface SSMAPI

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestTerminateSession verifies session termination
func TestTerminateSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := NewMockSSMAPI(ctrl) // skipcq: SCC-compile

	sessionID := "testID"

	gomock.InOrder(
		m.EXPECT().TerminateSession(&ssm.TerminateSessionInput{
			SessionId: &sessionID,
		}).Return(new(ssm.TerminateSessionOutput), nil),
	)

	assert.NoError(t, terminateSession(m, sessionID))
}
