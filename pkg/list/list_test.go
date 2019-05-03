package list

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInstancesString(t *testing.T) {
	resultText := `Index  Name       Instance ID          IP Address  Private DNS Name
1      testNode1  i-xxxxxxxxxxxxxxxx1  10.10.10.1  test1.local
2      testNode2  i-xxxxxxxxxxxxxxxx2  10.10.10.2  test2.local
`
	resultWide := `Index  Name       Instance ID          IP Address  Private DNS Name  Hostname       OS Name       OS Version  OS Type
1      testNode1  i-xxxxxxxxxxxxxxxx1  10.10.10.1  test1.local       testHostname1  Amazon Linux  2           Linux
2      testNode2  i-xxxxxxxxxxxxxxxx2  10.10.10.2  test2.local       testHostname2  Ubuntu        18.04       Linux
`
	resultJSON := `[{"hostname":"testHostname1","ip_address":"10.10.10.1","instance_id":"i-xxxxxxxxxxxxxxxx1","private_dns_name":"test1.local","instance_name":"testNode1","os_name":"Amazon Linux","os_type":"Linux","os_version":"2"},{"hostname":"testHostname2","ip_address":"10.10.10.2","instance_id":"i-xxxxxxxxxxxxxxxx2","private_dns_name":"test2.local","instance_name":"testNode2","os_name":"Ubuntu","os_type":"Linux","os_version":"18.04"}]
`
	resultYAML := `- hostname: testHostname1
  ip_address: 10.10.10.1
  instance_id: i-xxxxxxxxxxxxxxxx1
  private_dns_name: test1.local
  instance_name: testNode1
  os_name: Amazon Linux
  os_type: Linux
  os_version: "2"
- hostname: testHostname2
  ip_address: 10.10.10.2
  instance_id: i-xxxxxxxxxxxxxxxx2
  private_dns_name: test2.local
  instance_name: testNode2
  os_name: Ubuntu
  os_type: Linux
  os_version: "18.04"
`

	outputType := "instances"

	output := &StartOutput{
		Instances: []*Instance{
			&Instance{
				Hostname:       stringPointer("testHostname1"),
				IPAddress:      stringPointer("10.10.10.1"),
				InstanceID:     stringPointer("i-xxxxxxxxxxxxxxxx1"),
				PrivateDNSName: stringPointer("test1.local"),
				Name:           stringPointer("testNode1"),
				OSName:         stringPointer("Amazon Linux"),
				OSType:         stringPointer("Linux"),
				OSVersion:      stringPointer("2"),
			},
			&Instance{
				Hostname:       stringPointer("testHostname2"),
				IPAddress:      stringPointer("10.10.10.2"),
				InstanceID:     stringPointer("i-xxxxxxxxxxxxxxxx2"),
				PrivateDNSName: stringPointer("test2.local"),
				Name:           stringPointer("testNode2"),
				OSName:         stringPointer("Ubuntu"),
				OSType:         stringPointer("Linux"),
				OSVersion:      stringPointer("18.04"),
			},
		},
		outputType: &outputType,
	}

	assert := assert.New(t)

	output.format = stringPointer("wrong")
	_, err := output.String()
	assert.NotNil(err)
	output.format = stringPointer("text")
	outString, err := output.String()
	assert.Nil(err)
	assert.Equal(resultText, outString)
	output.format = stringPointer("wide")
	outString, err = output.String()
	assert.Nil(err)
	assert.Equal(resultWide, outString)
	output.format = stringPointer("json")
	outString, err = output.String()
	assert.Nil(err)
	assert.Equal(resultJSON, outString)
	output.format = stringPointer("yaml")
	outString, err = output.String()
	assert.Nil(err)
	assert.Equal(resultYAML, outString)
}

func TestSessionsString(t *testing.T) {
	resultText := `Index  Session ID       Target               Start Date
1      test-1234567890  i-xxxxxxxxxxxxxxxx1  2019-05-03T10:08:44Z
`
	resultWide := `Index  Session ID       Target               Start Date            Owner                                           Status
1      test-1234567890  i-xxxxxxxxxxxxxxxx1  2019-05-03T10:08:44Z  arn:aws:sts::0123456789:assumed-role/test/test  Connected
`
	resultJSON := `[{"session_id":"test-1234567890","target_instance":"i-xxxxxxxxxxxxxxxx1","status":"Connected","start_date":"2019-05-03T10:08:44Z","owner":"arn:aws:sts::0123456789:assumed-role/test/test"}]
`
	resultYAML := `- session_id: test-1234567890
  target_instance: i-xxxxxxxxxxxxxxxx1
  status: Connected
  start_date: "2019-05-03T10:08:44Z"
  owner: arn:aws:sts::0123456789:assumed-role/test/test
`

	outputType := "sessions"

	output := &StartOutput{
		Sessions: []*Session{
			&Session{
				SessionID: stringPointer("test-1234567890"),
				Target:    stringPointer("i-xxxxxxxxxxxxxxxx1"),
				Status:    stringPointer("Connected"),
				StartDate: stringPointer("2019-05-03T10:08:44Z"),
				Owner:     stringPointer("arn:aws:sts::0123456789:assumed-role/test/test"),
			},
		},
		outputType: &outputType,
	}

	assert := assert.New(t)

	output.format = stringPointer("wrong")
	_, err := output.String()
	assert.NotNil(err)
	output.format = stringPointer("text")
	outString, err := output.String()
	assert.Nil(err)
	assert.Equal(resultText, outString)
	output.format = stringPointer("wide")
	outString, err = output.String()
	assert.Nil(err)
	assert.Equal(resultWide, outString)
	output.format = stringPointer("json")
	outString, err = output.String()
	assert.Nil(err)
	assert.Equal(resultJSON, outString)
	output.format = stringPointer("yaml")
	outString, err = output.String()
	assert.Nil(err)
	assert.Equal(resultYAML, outString)
}

// Helper functions

func stringPointer(v string) *string {
	return &v
}
