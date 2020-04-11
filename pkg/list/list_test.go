package list

import (
	"testing"

	"github.com/danmx/sigil/pkg/aws"
	"github.com/danmx/sigil/pkg/aws/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// TestInstancesString verifies correctness of preparing a list of instanes
func TestInstancesString(t *testing.T) {
	resultText := `Index  Name       Instance ID          IP Address  Private DNS Name
1      testNode1  i-xxxxxxxxxxxxxxxx1  10.10.10.1  test1.local
2      testNode2  i-xxxxxxxxxxxxxxxx2  10.10.10.2  test2.local
`
	resultWide := `Index  Name       Instance ID          IP Address  Private DNS Name  Hostname       OS Name       OS Version  OS Type
1      testNode1  i-xxxxxxxxxxxxxxxx1  10.10.10.1  test1.local       testHostname1  Amazon Linux  2           Linux
2      testNode2  i-xxxxxxxxxxxxxxxx2  10.10.10.2  test2.local       testHostname2  Ubuntu        18.04       Linux
`
	resultJSON := `[{"hostname":"testHostname1","ip_address":"10.10.10.1","id":"i-xxxxxxxxxxxxxxxx1","private_dns_name":"test1.local","name":"testNode1","os_name":"Amazon Linux","os_type":"Linux","os_version":"2"},{"hostname":"testHostname2","ip_address":"10.10.10.2","id":"i-xxxxxxxxxxxxxxxx2","private_dns_name":"test2.local","name":"testNode2","os_name":"Ubuntu","os_type":"Linux","os_version":"18.04"}]
`
	resultYAML := `- hostname: testHostname1
  ip_address: 10.10.10.1
  id: i-xxxxxxxxxxxxxxxx1
  private_dns_name: test1.local
  name: testNode1
  os_name: Amazon Linux
  os_type: Linux
  os_version: "2"
- hostname: testHostname2
  ip_address: 10.10.10.2
  id: i-xxxxxxxxxxxxxxxx2
  private_dns_name: test2.local
  name: testNode2
  os_name: Ubuntu
  os_type: Linux
  os_version: "18.04"
`

	instances := []*aws.Instance{
		{
			Hostname:       "testHostname1",
			IPAddress:      "10.10.10.1",
			ID:             "i-xxxxxxxxxxxxxxxx1",
			PrivateDNSName: "test1.local",
			Name:           "testNode1",
			OSName:         "Amazon Linux",
			OSType:         "Linux",
			OSVersion:      "2",
		},
		{
			Hostname:       "testHostname2",
			IPAddress:      "10.10.10.2",
			ID:             "i-xxxxxxxxxxxxxxxx2",
			PrivateDNSName: "test2.local",
			Name:           "testNode2",
			OSName:         "Ubuntu",
			OSType:         "Linux",
			OSVersion:      "18.04",
		},
	}

	a := assert.New(t)

	_, err := instancesToString("wrong", instances)
	a.NotNil(err)
	outString, err := instancesToString(FormatText, instances)
	a.Nil(err)
	a.Equal(resultText, outString)
	outString, err = instancesToString(FormatWide, instances)
	a.Nil(err)
	a.Equal(resultWide, outString)
	outString, err = instancesToString(FormatJSON, instances)
	a.Nil(err)
	a.Equal(resultJSON, outString)
	outString, err = instancesToString(FormatYAML, instances)
	a.Nil(err)
	a.Equal(resultYAML, outString)
}

// TestSessionsString verifies correctness of preparing a list of sessions
func TestSessionsString(t *testing.T) {
	resultText := `Index  Session ID       Target               Start Date
1      test-1234567890  i-xxxxxxxxxxxxxxxx1  2019-05-03T10:08:44Z
`
	resultWide := `Index  Session ID       Target               Start Date            Owner                                           Status
1      test-1234567890  i-xxxxxxxxxxxxxxxx1  2019-05-03T10:08:44Z  arn:aws:sts::0123456789:assumed-role/test/test  Connected
`
	resultJSON := `[{"session_id":"test-1234567890","target":"i-xxxxxxxxxxxxxxxx1","status":"Connected","start_date":"2019-05-03T10:08:44Z","owner":"arn:aws:sts::0123456789:assumed-role/test/test"}]
`
	resultYAML := `- session_id: test-1234567890
  target: i-xxxxxxxxxxxxxxxx1
  status: Connected
  start_date: "2019-05-03T10:08:44Z"
  owner: arn:aws:sts::0123456789:assumed-role/test/test
`

	sessions := []*aws.Session{
		{
			SessionID: "test-1234567890",
			Target:    "i-xxxxxxxxxxxxxxxx1",
			Status:    "Connected",
			StartDate: "2019-05-03T10:08:44Z",
			Owner:     "arn:aws:sts::0123456789:assumed-role/test/test",
		},
	}

	a := assert.New(t)

	_, err := sessionsToString("wrong", sessions)
	a.NotNil(err)
	outString, err := sessionsToString(FormatText, sessions)
	a.Nil(err)
	a.Equal(resultText, outString)
	outString, err = sessionsToString(FormatWide, sessions)
	a.Nil(err)
	a.Equal(resultWide, outString)
	outString, err = sessionsToString(FormatJSON, sessions)
	a.Nil(err)
	a.Equal(resultJSON, outString)
	outString, err = sessionsToString(FormatYAML, sessions)
	a.Nil(err)
	a.Equal(resultYAML, outString)
}

// TestListInstances verifies listing instances
func TestListInstances(t *testing.T) {
	instances := []*aws.Instance{
		{
			Hostname:       "testHostname1",
			IPAddress:      "10.10.10.1",
			ID:             "i-xxxxxxxxxxxxxxxx1",
			PrivateDNSName: "test1.local",
			Name:           "testNode1",
			OSName:         "Amazon Linux",
			OSType:         "Linux",
			OSVersion:      "2",
		},
		{
			Hostname:       "testHostname2",
			IPAddress:      "10.10.10.2",
			ID:             "i-xxxxxxxxxxxxxxxx2",
			PrivateDNSName: "test2.local",
			Name:           "testNode2",
			OSName:         "Ubuntu",
			OSType:         "Linux",
			OSVersion:      "18.04",
		},
	}
	interactive := false
	format := FormatText
	input := StartInput{
		OutputFormat: &format,
		Interactive:  &interactive,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockCloudInstances(ctrl)

	m.EXPECT().ListInstances().Return(instances, nil)

	assert.NoError(t, input.listInstances(m))
	// TODO test integractive part
}

// TestListSessions verifies listing sessions
func TestListSessions(t *testing.T) {
	sessions := []*aws.Session{
		{
			SessionID: "test-1234567890",
			Target:    "i-xxxxxxxxxxxxxxxx1",
			Status:    "Connected",
			StartDate: "2019-05-03T10:08:44Z",
			Owner:     "arn:aws:sts::0123456789:assumed-role/test/test",
		},
	}
	interactive := false
	format := FormatText
	input := StartInput{
		OutputFormat: &format,
		Interactive:  &interactive,
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	m := mocks.NewMockCloudSessions(ctrl)

	m.EXPECT().ListSessions().Return(sessions, nil)

	assert.NoError(t, input.listSessions(m))
	// TODO test integractive part
}
