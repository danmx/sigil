package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringMapString(t *testing.T) {
	tagsString := "Name=test,env=dev,owner=example.com"
	mapTags := map[string]string{
		"Name":  "test",
		"env":   "dev",
		"owner": "example.com",
	}
	errorTagsString := "Name=test,env=dev,owner=example.com,MissingValue"

	mapType := &StringMapStringType{}

	assert := assert.New(t)
	err := mapType.Set(tagsString)
	assert.Nil(err)

	assert.EqualValues(mapTags, mapType.Map)
	err = mapType.Set(errorTagsString)
	assert.NotNil(err)

	err = mapType.Set(mapType.String())
	assert.Nil(err)
	assert.EqualValues(mapTags, mapType.Map)
}
